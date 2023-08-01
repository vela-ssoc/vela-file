package file

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/lua"
)

func (t *Type) String() string                         { return lua.B2S(t.Bytes()) }
func (t *Type) Type() lua.LValueType                   { return lua.LTObject }
func (t *Type) AssertFloat64() (float64, bool)         { return 0, false }
func (t *Type) AssertString() (string, bool)           { return "", false }
func (t *Type) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (t *Type) Peek() lua.LValue                       { return t }

func (t *Type) Bytes() []byte {
	chunk, _ := json.Marshal(t)
	return chunk
}

func (t *Type) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "err":
		if t.Err == nil {
			return lua.LNil
		}
		return lua.S2L(t.Err.Error())
	case "ext":
		return lua.S2L(t.Extension)
	case "type":
		return lua.S2L(t.Typ)
	case "subtype":
		return lua.S2L(t.Subtype)
	case "value":
		return lua.S2L(t.Value)
	}

	return lua.LNil
}
