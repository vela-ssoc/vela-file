package file

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/lua"
)

func newFileGlob(L *lua.LState) *glob {
	gl := &glob{
		co:   xEnv.Clone(L),
		pipe: pipe.New(pipe.Env(xEnv)),
	}

	n := L.GetTop()
	if n == 0 {
		return gl
	}

	for i := 1; i <= n; i++ {
		gl.append(L.CheckString(i))
	}
	return gl
}

func (gl *glob) String() string                         { return fmt.Sprintf("%p", gl) }
func (gl *glob) Type() lua.LValueType                   { return lua.LTObject }
func (gl *glob) AssertFloat64() (float64, bool)         { return 0, false }
func (gl *glob) AssertString() (string, bool)           { return "", false }
func (gl *glob) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (gl *glob) Peek() lua.LValue                       { return gl }

func (gl *glob) append(pattern string) {
	gl.patterns = append(gl.patterns, pattern)
}

func (gl *glob) pipeL(L *lua.LState) int {
	gl.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (gl *glob) runL(L *lua.LState) int {
	gl.run()
	return 0
}

func (gl *glob) wrapL(L *lua.LState) int {
	err := gl.err.Wrap()
	if err != nil {
		L.Push(lua.S2L(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func (gl *glob) r() lua.LValue {
	rn := len(gl.result)
	if rn == 0 {
		return lua.Slice{}
	}

	rv := lua.Slice{}
	for i := 0; i < rn; i++ {
		rv = append(rv, lua.S2L(gl.result[i]))
	}

	return rv
}

func (gl *glob) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return lua.NewFunction(gl.pipeL)
	case "run":
		return lua.NewFunction(gl.runL)
	case "wrap":
		return lua.NewFunction(gl.wrapL)

	case "result":
		return gl.r()

	}

	return lua.LNil

}
