package file

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
	"time"
)

func (xf *xFile) pushL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LNil)
		return 0
	}

	for i := 1; i <= n; i++ {
		xf.Push(auxlib.Format(L, 0))
	}
	return 0

}

func (xf *xFile) backupL(L *lua.LState) int {
	filename := xf.filename(time.Now())

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		xEnv.Errorf("file backup fail %v", err)
		return 0
	}

	old := xf.fd
	defer old.Close()
	xf.fd = fd
	return 0
}

func (xf *xFile) dayL(L *lua.LState) int {
	format := L.IsString(1)
	if format == "" {
		format = "2006-01-02"
	}

	return 0
}

func (xf *xFile) startL(L *lua.LState) int {
	xEnv.Start(L, xf).From(L.CodeVM()).Do()
	return 0
}

func (xf *xFile) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "push":
		return lua.NewFunction(xf.pushL)
	case "backup":
		return lua.NewFunction(xf.backupL)
	case "day":
		return lua.NewFunction(xf.dayL)
	default:
		return lua.LNil
	}

	return lua.LNil
}
