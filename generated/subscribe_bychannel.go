// Start fileheader template
// Code generated by go generate; DO NOT EDIT.
// This file was generated by FactomGenerate robots

// Start Generated Code

package generated

import (
	. "github.com/FactomProject/factomd/common/pubsubtypes"
	. "github.com/FactomProject/factomd/pubsub"
)

// End fileheader template

// Start subscribeBychannel generated go code

// Channel subscriber has the basic necessary function implementations. All this does is add a wrapper with typing.
type Subscribe_ByChannel_Hash_type struct {
	*Channel
}

// type the Read function
func (s *Subscribe_ByChannel_Hash_type) Read() Hash {
	return s.Channel.Read().(Hash) // cast the return to the specific type
}

// type the ReadWithInfo function
func (s *Subscribe_ByChannel_Hash_type) ReadWithInfo() (Hash, bool) {
	v, ok := <-s.Updates
	return v.(Hash), ok
}

// Create a typed instance form a generic instance
func Subscribe_ByChannel_Hash(p *Channel) *Subscribe_ByChannel_Hash_type {
	return &Subscribe_ByChannel_Hash_type{p}
}

// End subscribe_bychannel generated code
//
// Start subscribeBychannel generated go code

// Channel subscriber has the basic necessary function implementations. All this does is add a wrapper with typing.
type Subscribe_ByChannel_IMsg_type struct {
	*Channel
}

// type the Read function
func (s *Subscribe_ByChannel_IMsg_type) Read() IMsg {
	return s.Channel.Read().(IMsg) // cast the return to the specific type
}

// type the ReadWithInfo function
func (s *Subscribe_ByChannel_IMsg_type) ReadWithInfo() (IMsg, bool) {
	v, ok := <-s.Updates
	return v.(IMsg), ok
}

// Create a typed instance form a generic instance
func Subscribe_ByChannel_IMsg(p *Channel) *Subscribe_ByChannel_IMsg_type {
	return &Subscribe_ByChannel_IMsg_type{p}
}

// End subscribe_bychannel generated code
//
// Start filetail template
// Code generated by go generate; DO NOT EDIT.
// End filetail template
// End Generated Code
