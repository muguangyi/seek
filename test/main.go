// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/seek/seek"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go seek.Serve("127.0.0.1:9999")

	go seek.Startup("127.0.0.1:9999", "util",
		seek.NewSignaler("IMath", newMath(), true))

	go seek.Startup("127.0.0.1:9999", "logic",
		seek.NewSignaler("IGame", newGame(&wg), true),
		seek.NewSignaler("ILobby", newLobby(&wg), true))

	wg.Wait()
	seek.Close()
	fmt.Println("Completed!")
}

//go:generate seek.gen
