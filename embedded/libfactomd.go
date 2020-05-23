// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "C"
import (
	"fmt"
	"github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/modules/worker"
	"os"
	"reflect"
)

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
	engine.Run(params)
}

//export Shutdown
func Shutdown() {
	worker.SendSigInt()
}

func main() {
	Serve()
}
