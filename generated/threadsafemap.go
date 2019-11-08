// Start fileheader template
// Code generated by go generate; DO NOT EDIT.
// This file was generated by FactomGenerate robots

// Start Generated Code

package generated

import (
	"sync"

	"github.com/FactomProject/factomd/common"
)

// End fileheader template

// Start threadsafemap generated go code

type foo struct {
	sync.Mutex
	common.Name
	internalMap map[int]string
}

func (q *foo) Init(parent common.NamedObject, name string, size int) *foo {
	q.Name.Init(parent, name)
	q.internalMap = make(map[int]string, size)
	return q
}

func (q *foo) Put(index int, value string) {
	q.Lock()
	q.internalMap[index] = value
	q.Unlock()
}

func (q *foo) Get(index int) string {
	q.Lock()
	defer q.Unlock()
	return q.internalMap[index]
}

// End threadsafemap generated go code
//
// Start filetail template
// Code generated by go generate; DO NOT EDIT.
// End filetail template
// End Generated Code