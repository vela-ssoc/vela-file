package file

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/grep"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/lua"
	"path/filepath"
)

/*
	local d = file.dir("/var/log")
	d.filter('*.log')
	d.pipe(function(v)

	end)
	d.run()

*/

func (d *dir) String() string                         { return fmt.Sprintf("%p", d) }
func (d *dir) Type() lua.LValueType                   { return lua.LTObject }
func (d *dir) AssertFloat64() (float64, bool)         { return 0, false }
func (d *dir) AssertString() (string, bool)           { return "", false }
func (d *dir) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (d *dir) Peek() lua.LValue                       { return d }

func (d *dir) visit(L *lua.LState) int {
	if !d.ok() {
		return 0
	}

	n := len(d.data)
	if n == 0 {
		return 0
	}

	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	for i := 0; i < n; i++ {
		file := d.data[i]
		info, err := file.Info()
		fi := NewLInfo(filepath.Join(d.path, file.Name()), info, err)

		if len(d.filter) > 0 && !match(fi.path, d.filter) {
			continue
		}
		d.pipe.Do(fi, co, func(err error) {
			xEnv.Errorf("%s pipe call fail %v", fi.path, err)
		})
	}

	return 0

}

func (d *dir) pipeL(L *lua.LState) int {
	d.pipe.CheckMany(L, pipe.Seek(0))
	return d.visit(L)
}

func (d *dir) filterL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		return 0
	}

	for i := 1; i <= n; i++ {
		d.filter = append(d.filter, grep.New(L.IsString(i)))
	}

	return 0
}

func (d *dir) r() lua.LValue {
	n := len(d.data)
	if n == 0 {
		return lua.Slice{}
	}

	var rv lua.Slice
	filter := d.fuzzy()

	for i := 0; i < n; i++ {
		fi := d.Info(i)
		if filter(fi.path) {
			rv = append(rv, fi)
		}
	}
	return rv
}

func (d *dir) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(d.ok())
	case "err":
		if d.ok() {
			return lua.LNil
		}
		return lua.S2L(d.err.Error())

	case "count":
		return lua.LInt(len(d.data))
	case "filter":
		return L.NewFunction(d.filterL)
	case "pipe":
		return L.NewFunction(d.pipeL)
	case "result":
		return d.r()
	}

	return lua.LNil
}
