package file

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"os"
	"path/filepath"
)

type dir struct {
	path   string
	data   []os.DirEntry
	filter []func(string) bool
	pipe   *pipe.Chains
	err    error
}

func newLuaFileDir(L *lua.LState) int {
	path := L.CheckString(1)
	data, err := os.ReadDir(path)
	L.Push(&dir{
		path: path,
		data: data,
		err:  err,
		pipe: pipe.New(pipe.Env(xEnv)),
	})
	return 1
}

func (d *dir) ok() bool {
	if d.err != nil {
		return false
	}

	return true
}

func (d *dir) fuzzy() func(string) bool {
	if len(d.filter) == 0 {
		return func(_ string) bool {
			return true
		}
	}

	return func(path string) bool {
		return match(path, d.filter)
	}
}

func (d *dir) Info(idx int) LInfo {
	fd, err := d.data[idx].Info()
	path := filepath.Join(d.path, fd.Name())
	return NewLInfo(path, fd, err)
}
