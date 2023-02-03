package file

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

/*
	file.scan("aaa.txt").pipe()
*/

func (ls *LScan) String() string                         { return auxlib.B2S(ls.Byte()) }
func (ls *LScan) Type() lua.LValueType                   { return lua.LTObject }
func (ls *LScan) AssertFloat64() (float64, bool)         { return 0, false }
func (ls *LScan) AssertString() (string, bool)           { return "", false }
func (ls *LScan) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (ls *LScan) Peek() lua.LValue                       { return ls }

func (ls *LScan) ignoreL(L *lua.LState) int {
	ls.ignore.CheckMany(L)
	L.Push(ls)
	return 1
}
func (ls *LScan) filterL(L *lua.LState) int {
	ls.filter.CheckMany(L)
	L.Push(ls)
	return 1
}

func (ls *LScan) pipeL(L *lua.LState) int {
	ls.pipe.CheckMany(L)
	ls.run()
	L.Push(ls)
	return 1
}

func (ls *LScan) doL(L *lua.LState) int {
	ls.run()
	return 0
}

func (ls *LScan) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ignore":
		return lua.NewFunction(ls.ignoreL)
	case "filter":
		return lua.NewFunction(ls.filterL)
	case "pipe":
		return lua.NewFunction(ls.pipeL)
	case "run":
		return lua.NewFunction(ls.doL)
	}

	return lua.LNil
}
