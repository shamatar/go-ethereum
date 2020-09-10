// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"fmt"
	"sync"

	"github.com/holiman/uint256"
)

const limit = 2048

var stackPool = sync.Pool{
	New: func() interface{} {
		st := new(Stack)
		return st
	},
}

// Stack is an object for basic stack operations. Items popped to the stack are
// expected to be changed and modified. stack does not take care of adding newly
// initialised objects.
type Stack struct {
	data   [limit]uint256.Int
	length int
}

func newstack() *Stack {
	return stackPool.Get().(*Stack)
}

func returnStack(s *Stack) {
	s.length = 0
	stackPool.Put(s)
}

// Data returns the underlying uint256.Int slice
func (st *Stack) Data() []uint256.Int {
	tmp := make([]uint256.Int, st.length)
	copy(tmp, st.data[:st.length])
	return tmp
}

func (st *Stack) push(d uint256.Int) {
	// NOTE push limit (1024) is checked in baseCheck
	st.data[st.length] = d
	st.length++
}
func (st *Stack) pushN(ds ...uint256.Int) {
	toAdd := len(ds)
	copy(st.data[st.length:st.length+toAdd], ds)
	st.length += toAdd
}

func (st *Stack) pop() (ret uint256.Int) {
	ret = st.data[st.length-1]
	st.length--
	return
}

func (st *Stack) len() int {
	return st.length
}

func (st *Stack) swap(n int) {
	st.data[st.length-n], st.data[st.length-1] = st.data[st.length-1], st.data[st.length-n]
}

func (st *Stack) dup(n int) {
	st.push(st.data[st.length-n])
}

func (st *Stack) peek() *uint256.Int {
	return &st.data[st.length-1]
}

// Back returns the n'th item in stack
func (st *Stack) Back(n int) *uint256.Int {
	return &st.data[st.length-n-1]
}

// Print dumps the content of the stack
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		for i, val := range st.data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}

var rStackPool = sync.Pool{
	New: func() interface{} {
		return &ReturnStack{data: make([]uint32, 0, 10)}
	},
}

// ReturnStack is an object for basic return stack operations.
type ReturnStack struct {
	data []uint32
}

func newReturnStack() *ReturnStack {
	return rStackPool.Get().(*ReturnStack)
}

func returnRStack(rs *ReturnStack) {
	rs.data = rs.data[:0]
	rStackPool.Put(rs)
}

func (st *ReturnStack) push(d uint32) {
	st.data = append(st.data, d)
}

// A uint32 is sufficient as for code below 4.2G
func (st *ReturnStack) pop() (ret uint32) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}
