// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"io"

	"github.com/muguangyi/ferry/codec"
)

const (
	cError                uint = 0  // Error
	cHeartbeat            uint = 1  // Heartbeat
	cRegisterRequest      uint = 2  // Register request
	cHubRegisterResponse  uint = 3  // Register hub response
	cDockRegisterResponse uint = 4  // Register dock response
	cImportRequest        uint = 5  // Import request
	cImportResponse       uint = 6  // Import response
	cQueryRequest         uint = 7  // Query request
	cQueryResponse        uint = 8  // Query response
	cRpcRequest           uint = 9  // RPC request
	cRpcResponse          uint = 10 // RPC response
)

func protoMaker(id uint) IProto {
	switch id {
	case cError:
		return new(protoError)
	case cHeartbeat:
		return new(protoHeartbeat)
	case cRegisterRequest:
		return new(protoRegisterRequest)
	case cHubRegisterResponse:
		return new(protoHubRegisterResponse)
	case cDockRegisterResponse:
		return new(protoDockRegisterResponse)
	case cImportRequest:
		return new(protoImportRequest)
	case cImportResponse:
		return new(protoImportResponse)
	case cQueryRequest:
		return new(protoQueryRequest)
	case cQueryResponse:
		return new(protoQueryResponse)
	case cRpcRequest:
		return new(protoRpcRequest)
	case cRpcResponse:
		return new(protoRpcResponse)
	}

	return nil
}

// Error
type protoError struct {
	Error string
}

func (p *protoError) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Error).Encode(writer)
}

func (p *protoError) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Error, err = any.String()
	return err
}

// Heartbeat
type protoHeartbeat struct {
}

func (p *protoHeartbeat) Marshal(writer io.Writer) error {
	return nil
}

func (p *protoHeartbeat) Unmarshal(reader io.Reader) error {
	return nil
}

// Register request
type protoRegisterRequest struct {
	Slots []string
}

func (p *protoRegisterRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slots).Encode(writer)
}

func (p *protoRegisterRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Slots = make([]string, len(arr))
	for i, iv := range arr {
		p.Slots[i] = iv.(string)
	}

	return nil
}

// Register hub response
type protoHubRegisterResponse struct {
	Port int
}

func (p *protoHubRegisterResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Port).Encode(writer)
}

func (p *protoHubRegisterResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Port, err = any.Int()
	return err
}

// Register dock response
type protoDockRegisterResponse struct {
	Slot string
}

func (p *protoDockRegisterResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slot).Encode(writer)
}

func (p *protoDockRegisterResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Slot, err = any.String()
	return err
}

// Import request
type protoImportRequest struct {
	Slots []string
}

func (p *protoImportRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slots).Encode(writer)
}

func (p *protoImportRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Slots = make([]string, len(arr))
	for i, iv := range arr {
		p.Slots[i] = iv.(string)
	}

	return nil
}

// Import response
type protoImportResponse struct {
	Docks []string
}

func (p *protoImportResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Docks).Encode(writer)
}

func (p *protoImportResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Docks = make([]string, len(arr))
	for i, iv := range arr {
		p.Docks[i] = iv.(string)
	}

	return nil
}

// Query request
type protoQueryRequest struct {
	Slot string
}

func (p *protoQueryRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slot).Encode(writer)
}

func (p *protoQueryRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Slot, err = any.String()
	return err
}

// Query response
type protoQueryResponse struct {
	DockAddr string
}

func (p *protoQueryResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.DockAddr).Encode(writer)
}

func (p *protoQueryResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.DockAddr, err = any.String()
	return err
}

// RPC request
type protoRpcRequest struct {
	Index      int64
	Slot       string
	Method     string
	Args       []interface{}
	WithResult bool
}

func (p *protoRpcRequest) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Slot).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Method).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Args).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.WithResult).Encode(writer)
	if nil != err {
		return err
	}

	return nil
}

func (p *protoRpcRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)

	err := any.Decode(reader)
	if nil != err {
		return err
	}
	p.Index, err = any.Int64()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Slot, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Method, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Args, err = any.Arr()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.WithResult, err = any.Bool()
	if nil != err {
		return err
	}

	return nil
}

// RPC response
type protoRpcResponse struct {
	Index  int64
	Slot   string
	Method string
	Result []interface{}
	Err    string
}

func (p *protoRpcResponse) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Slot).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Method).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Result).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Err).Encode(writer)
	if nil != err {
		return err
	}

	return nil
}

func (p *protoRpcResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)

	err := any.Decode(reader)
	if nil != err {
		return err
	}
	p.Index, err = any.Int64()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Slot, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Method, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Result, err = any.Arr()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Err, err = any.String()
	if nil != err {
		return err
	}

	return nil
}
