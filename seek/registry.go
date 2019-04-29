// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

var registry map[string]interface{} = make(map[string]interface{})

func register(name string, maker interface{}) {
	registry[name] = maker
}

func tryMake(name string, s ISignaler) (interface{}, bool) {
	maker := registry[name]
	if nil != maker {
		return maker.(func(signaler ISignaler) interface{})(s), true
	}

	return nil, false
}
