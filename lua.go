package file

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"io"
	"os"
)

var (
	xEnv vela.Environment
	max  int64 = 1024 * 1024 * 30
)

/*
	local w = file.open{name , path , delim}
	local st = file.stat(path)
	local wk = file.walk("name")

	local wx = wk.open("/var/log")
	wx.ext(".zip" , ".txt")
	wx.limit()
	wx.run()

*/

func newLuaFileOpen(L *lua.LState) int {
	cfg := newConfig(L)
	var ov *xFile

	proc := L.NewVelaData(cfg.name, fileTypeOf)
	if proc.IsNil() {
		ov = newFile(cfg)
		proc.Set(ov)
	} else {
		ov = proc.Data.(*xFile)
		ov.cfg = cfg
	}

	xEnv.Start(L, ov).From(L.CodeVM()).Do()
	L.Push(proc)
	return 1
}

func newLuaFileStat(L *lua.LState) int {
	path := L.IsString(1)
	if path == "" {
		return 0
	}

	fd, err := os.Stat(path)
	L.Push(NewLInfo(path, fd, err))
	return 1
}

func newLuaFileWalk(L *lua.LState) int {
	cfg := newWalkConfig(L)
	proc := L.NewVelaData(cfg.name, walkTypeof)
	if proc.IsNil() {
		proc.Set(newWalk(cfg))
	} else {
		old := proc.Data.(*walk)
		old.Close()
		proc.Set(newWalk(cfg))
	}

	L.Push(proc)
	return 1
}

func newLuaFileGlob(L *lua.LState) int {
	L.Push(newFileGlob(L))
	return 1
}

func newLuaFileScan(L *lua.LState) int {
	file := L.CheckString(1)
	s := NewLScan(L, file)
	L.Push(s)
	return 1
}

func newLuaFileReadAll(L *lua.LState) int {
	content := &Content{}
	var fd *os.File

	todo := func() int {
		L.Use(content)
		return 1
	}

	path := L.CheckString(1)
	fd, err := os.Open(path)
	if err != nil {
		content.err = err
		return todo()
	}

	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		content.err = err
		return todo()
	}

	if stat.Size() > max {
		content.err = fmt.Errorf("%s too big , size:%d > %d", path, stat.Size(), max)
		return todo()
	}

	data, err := io.ReadAll(fd)
	if err != nil {
		return todo()
	}

	content.data = data
	return todo()
}

func newLuaFileTypeL(L *lua.LState) int {
	path := L.CheckString(1)
	if len(path) == 0 {
		L.RaiseError("invalid file path")
		return 0
	}

	ft := Filetype(path)
	L.Push(ft)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	file := lua.NewUserKV()
	file.Set("open", lua.NewFunction(newLuaFileOpen))
	file.Set("dir", lua.NewFunction(newLuaFileDir))
	file.Set("stat", lua.NewFunction(newLuaFileStat))
	file.Set("walk", lua.NewFunction(newLuaFileWalk))
	file.Set("find", lua.NewFunction(newLuaFileWalk))
	file.Set("glob", lua.NewFunction(newLuaFileGlob))
	file.Set("scan", lua.NewFunction(newLuaFileScan))
	file.Set("read_all", lua.NewFunction(newLuaFileReadAll))
	file.Set("cat", lua.NewFunction(newLuaFileReadAll))
	file.Set("type", lua.NewFunction(newLuaFileTypeL))
	env.Global("file", file)
}
