package file

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"os"
	"path/filepath"
)

type LInfo struct {
	path string
	fd   os.FileInfo
	ext  string
	err  error
}

func (i LInfo) String() string                         { return auxlib.B2S(i.Byte()) }
func (i LInfo) Type() lua.LValueType                   { return lua.LTObject }
func (i LInfo) AssertFloat64() (float64, bool)         { return 0, false }
func (i LInfo) AssertString() (string, bool)           { return "", false }
func (i LInfo) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (i LInfo) Peek() lua.LValue                       { return i }

func NewLInfo(path string, fd os.FileInfo, err error) LInfo {
	return LInfo{
		path: path,
		fd:   fd,
		ext:  filepath.Ext(path),
		err:  err,
	}
}

func (i LInfo) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("path", i.path)
	enc.KV("ext", i.ext)
	enc.KV("mtime", i.MTime())
	enc.KV("size", i.fd.Size())
	enc.End("}")
	return enc.Bytes()
}

func (i LInfo) ok() bool {
	return i.err == nil

}

func (i LInfo) MTime() int64 {
	if i.ok() {
		return i.fd.ModTime().Unix()
	}
	return 0
}

func (i LInfo) showL(L *lua.LState) int {
	if L.Console == nil {
		return 0
	}
	L.Output(i.String())
	return 0
}

func (i LInfo) ReadAll() []byte {
	f, err := os.Open(i.path)
	if err != nil {
		xEnv.Error("%s open fail %v", i.path, err)
		return nil
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		xEnv.Error("%s read fail %v", i.path, err)
		return nil
	}
	return content
}

func (i LInfo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(i.ok())

	case "scan":
		return lua.NewFunction(func(co *lua.LState) int {
			L.Push(NewLScan(co, i.path))
			return 1
		})

	case "content":
		return lua.B2L(i.ReadAll())

	case "name":
		return lua.S2L(i.fd.Name())

	case "path":
		return lua.S2L(i.path)

	case "ext":
		return lua.S2L(i.ext)

	case "size":
		return lua.LNumber(i.fd.Size())

	case "mtime":
		return lua.LNumber(i.MTime())

	case "ctime":
		return lua.LNumber(i.ctime())

	case "atime":
		return lua.LNumber(i.atime())

	case "dir":
		return lua.LBool(i.fd.IsDir())

	case "not_ext":
		return lua.LBool(i.ext == "")

	case "show":
		return lua.NewFunction(i.showL)

	}

	return lua.LNil
}
