package file

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

const (
	TIMESTAMP = 1609430400
	DAY       = 86400
	HOUR      = 3600
)

type config struct {
	name  string
	path  string
	delim string
}

func newConfig(L *lua.LState) *config {
	val := L.Get(1)
	//tab := L.CheckTable(1)
	cfg := &config{
		name:  "file",
		delim: "\n",
	}

	switch val.Type() {
	case lua.LTString:
		cfg.path = val.String()

	case lua.LTTable:
		tab := val.(*lua.LTable)
		tab.Range(func(key string, val lua.LValue) {
			switch key {
			case "name":
				cfg.name = val.String()
			case "path":
				cfg.path = val.String()
			default:
				L.RaiseError("invalid %s field", key)
				return
			}
		})

	default:
		L.RaiseError("invalid config type %s", val.Type().String())
		return cfg
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := auxlib.Name(cfg.name); e != nil {
		return e
	}

	if e := auxlib.Warp(cfg.delim); e != nil {
		return e
	}
	return nil
}
