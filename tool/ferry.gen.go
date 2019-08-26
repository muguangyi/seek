// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	cFileExt string = ".go"
	cGenFile string = "proxy.gen.go"
)

func main() {
	flag.Parse()

	log.Println("--- start parsing ---")
	targets := make([]*astTarget, 0)
	fset := token.NewFileSet()
	filepath.Walk(".", func(file string, fi os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if fi.IsDir() || ignore(file) {
			return nil
		}

		log.Println(file)
		targets = append(targets, scanFile(fset, file)...)

		return nil
	})
	log.Println("--- parse completed ---")

	if len(targets) > 0 {
		log.Println("--- start generating ---")
		if data, ok := genProxyCode(targets); ok {
			save(cGenFile, data)
		}
		log.Println("--- generate completed ---")
	}
}

func ignore(file string) bool {
	ext := path.Ext(file)
	if cFileExt != ext {
		return true
	}

	if strings.Contains(file, cGenFile) {
		return true
	}

	return false
}

func scanFile(fset *token.FileSet, file string) []*astTarget {
	f, err := parser.ParseFile(fset, file, nil, 0)
	if nil != err {
		log.Fatal(err)
	}

	targets := make([]*astTarget, 0)
	ast.Inspect(f, func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}

		if target, ok := parseFile(file, decl); ok {
			targets = append(targets, target)
		}

		return true
	})

	return targets
}

func parseFile(file string, decl *ast.GenDecl) (target *astTarget, ok bool) {
	if token.TYPE != decl.Tok {
		ok = false
		return
	}

	if len(decl.Specs) > 1 {
		ok = false
		return
	}

	spec := decl.Specs[0].(*ast.TypeSpec)
	itype, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return
	}

	ok = true
	target = new(astTarget)
	target.Source = file
	target.Name = spec.Name.Name
	target.Proxy = strings.ToLower(target.Name + "proxy")

	target.Methods = make(map[string]*astMethod)
	for _, field := range itype.Methods.List {
		m := new(astMethod)
		m.Name = field.Names[0].Name
		m.Params = make([]*astField, 0)
		m.Results = make([]*astField, 0)

		ftype, _ := field.Type.(*ast.FuncType)

		if nil != ftype.Params {
			m.PCount = len(ftype.Params.List)
			if m.PCount > 0 {
				for _, param := range ftype.Params.List {
					m.Params = append(m.Params, &astField{
						Name: param.Names[0].Name,
						Type: param.Type.(*ast.Ident).Name,
					})
				}
			}
		}

		if nil != ftype.Results {
			m.RCount = len(ftype.Results.List)
			if m.RCount > 0 {
				for _, result := range ftype.Results.List {
					m.Results = append(m.Results, &astField{
						Name: func() string {
							if len(result.Names) > 0 {
								return result.Names[0].Name
							}

							return ""
						}(),
						Type: result.Type.(*ast.Ident).Name,
					})
				}
			}
		}

		target.Methods[m.Name] = m
	}

	return
}

func genProxyCode(targets []*astTarget) (data []byte, ok bool) {
	input := map[string]interface{}{
		"date":    time.Now().Format("2006-01-02 15:04:05"),
		"targets": targets,
	}

	t, err := template.New("").Funcs(template.FuncMap{"comma": func(index int, length int) string {
		if (index + 1) < length {
			return ","
		}

		return ""
	}}).Funcs(template.FuncMap{"default": func(t string) string {
		switch t {
		case "bool":
			return "false"
		case "int":
			return "0"
		case "int8":
			return "0"
		case "int16":
			return "0"
		case "int32":
			return "0"
		case "int64":
			return "0"
		case "uint":
			return "0"
		case "uint8":
			return "0"
		case "uint16":
			return "0"
		case "uint32":
			return "0"
		case "uint64":
			return "0"
		case "float32":
			return "0"
		case "float64":
			return "0"
		default:
			return "nil"
		}
	}}).Parse(gencode)
	if nil != err {
		log.Fatal(err)
	}

	buff := bytes.NewBufferString("")
	err = t.Execute(buff, input)
	if nil != err {
		log.Fatal(err)
	}

	data, err = format.Source(buff.Bytes())
	if nil != err {
		log.Fatal(err)
	}

	ok = true
	return
}

func save(target string, data []byte) error {
	return ioutil.WriteFile(target, data, 0644)
}

type astField struct {
	Name string
	Type string
}

type astMethod struct {
	Name    string
	PCount  int
	Params  []*astField
	RCount  int
	Results []*astField
}

type astTarget struct {
	Source  string
	Name    string
	Proxy   string
	Methods map[string]*astMethod
}

const gencode = `
//
// This code was generated by ferry tool.
// 
// Changes to this file may cause incorrect behavior and will be lost if the code is regenerated.
// 
// {{.date}}
//

package main

import (
	"github.com/muguangyi/ferry"
)

{{range $index, $target := .targets}}
// {{$target.Name}} from: {{$target.Source}}
type {{$target.Proxy}} struct {
	slot ferry.ISlot
}

{{range $i, $method := $target.Methods}}
func (p *{{$target.Proxy}}) {{$method.Name}}({{range $j, $param := $method.Params}}{{$param.Name}} {{$param.Type}}{{if le $j $method.PCount}}, {{end}}{{end}}){{if gt $method.RCount 0}}({{range $k, $result := $method.Results}}{{$result.Type}}{{comma $k $method.RCount}}{{end}}){{end}}{
	{{if gt $method.RCount 0}}results, err := p.slot.CallWithResult("{{$target.Name}}", "{{$method.Name}}", {{range $j, $param := $method.Params}}{{$param.Name}}{{comma $j $method.PCount}}{{end}})
	{{else}}err := p.slot.Call("{{$target.Name}}", "{{$method.Name}}", {{range $j, $param := $method.Params}}{{$param.Name}}{{comma $j $method.PCount}}{{end}})
	{{end}}
	if nil != err {
		return {{range $k, $result := $method.Results}}{{default $result.Type}}{{comma $k $method.RCount}}{{end}}
	}{{if gt $method.RCount 0}}
	
	return {{range $k, $result := $method.Results}}results[{{$k}}].({{$result.Type}}){{comma $k $method.RCount}}{{end}}{{end}}
}
{{end}}
// end {{$target.Name}}
{{end}}

// Register to ferry
var (
{{range $index, $target := .targets}}{{$target.Proxy}}succ bool = ferry.Register("{{$target.Name}}", func(slot ferry.ISlot) interface{} { return &{{$target.Proxy}}{slot: slot} })
{{end}}
)
`
