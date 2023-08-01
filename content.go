package file

import (
	"github.com/valyala/fasttemplate"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
)

type Content struct {
	err  error
	size int64
	data []byte
}

func (c *Content) String() string                         { return lua.B2S(c.data) }
func (c *Content) Type() lua.LValueType                   { return lua.LTObject }
func (c *Content) AssertFloat64() (float64, bool)         { return 0, false }
func (c *Content) AssertString() (string, bool)           { return c.String(), true }
func (c *Content) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c *Content) Peek() lua.LValue                       { return c }

func (c *Content) execL(L *lua.LState) int {
	if !c.ok() {
		L.Push(lua.S2L("not found data"))
		return 1
	}

	lv := L.Get(1)
	obj, ok := lv.(lua.IndexEx)
	if !ok {
		L.Push(lua.S2L("invalid object must have _index , got:" + lv.Type().String()))
		return 1
	}
	t := fasttemplate.New(lua.B2S(c.data), "${", "}")

	s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		item := obj.Index(L, tag)
		return w.Write(lua.S2B(item.String()))
	})

	L.Push(lua.S2L(s))
	return 1
}

func (c *Content) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(c.ok())
	case "size":
		return lua.LInt(c.size)
	case "data":
		return lua.B2L(c.data)
	case "exec":
		return lua.NewFunction(c.execL)
	}

	return lua.LNil
}

func (c *Content) ok() bool {
	return c.err == nil
}
