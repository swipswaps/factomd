// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "C"
import (
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/fnode"
	"github.com/FactomProject/factomd/modules/ipfs"
	"os"
	"reflect"
	"time"

	"github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/modules/registry"
	"github.com/FactomProject/factomd/modules/worker"
	"github.com/FactomProject/factomd/wsapi"
)

var state interfaces.IState

// export V2Api
func V2Api(method string, params string) string {
	r := new(primitives.JSON2Request)
	r.Method = method
	// FIXME parse params and populate
	//r.Params =
	res, err := wsapi.HandleV2JSONRequest(state, r)
	_ = res
	_ = err
	return `{ "foo": "bar" }`
}

//export Serve
func Serve() {
	fmt.Println("Command Line Arguments:")

	for _, v := range os.Args[1:] {
		fmt.Printf("\t%s\n", v)
	}

	// REVIEW: does this make sense to pass into from python
	params := engine.ParseCmdLine(os.Args[1:])

	fmt.Println("Parameter:")
	s := reflect.ValueOf(params).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%30s %s = %v\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	fmt.Println()
	p := registry.New()
	p.Register(func(w *worker.Thread) {
		engine.Factomd(w, params, params.Sim_Stdin)
		state = fnode.Get(0).State
		ipfs.Start(w, state)

		w.OnRun(func() {
			// poppulate local state
			fmt.Print("Running")
		})
		w.OnComplete(func() {
			fmt.Println("Waiting to Shut Down")
			time.Sleep(time.Second * 5)
		})
	})
	p.Run()
}

//export Shutdown
func Shutdown() {
	worker.SendSigInt()
}

func main() {
	Serve()
}
