package file

import (
	"bufio"
	"context"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	vswitch "github.com/vela-ssoc/vela-switch"
	"os"
)

/*
	local s = file.scan("www.baidu.com").case().pipe()
	s.case("* re 123").pipe()
	s.case("* re 123").pipe()
	s.case("* re 123").pipe()
	s.pipe()
	s.do()


*/

type LScan struct {
	ignore *cond.Ignore
	filter *cond.Combine
	ctx    context.Context
	file   string
	fd     *os.File
	err    error
	vsh    *vswitch.Switch
	pipe   *pipe.Chains
	co     *lua.LState
}

func (ls *LScan) Byte() []byte {
	return nil
}

func (ls *LScan) ok() bool {
	return ls.err == nil
}

func (ls *LScan) run() {
	if !ls.ok() {
		xEnv.Errorf("%s scan fail %v", ls.file, ls.err)
		return
	}

	defer func() {
		ls.fd.Close()
		ls.fd = nil
	}()

	scan := bufio.NewScanner(ls.fd)
	line := 0
	for scan.Scan() {
		select {
		case <-ls.co.Context().Done():
			xEnv.Errorf("%s scan over", ls.file)
			return
		default:
			text := scan.Text()
			line++
			if ls.ignore.Match(text) || !ls.filter.Match(text) {
				goto next
			}
			ls.vsh.Do(text)
			ls.pipe.Call2(text, line, ls.co)

		next:
			if err := scan.Err(); err != nil {
				xEnv.Errorf("%s scan over %v", ls.file, err)
				return
			}
		}
	}
}

func NewLScan(L *lua.LState, file string) *LScan {
	ls := &LScan{
		co:     xEnv.Clone(L),
		file:   file,
		pipe:   pipe.New(pipe.Env(xEnv)),
		vsh:    vswitch.NewL(L),
		ignore: cond.NewIgnore(),
		filter: cond.NewCombine(),
	}

	fd, err := os.Open(ls.file)
	ls.fd = fd
	ls.err = err
	return ls
}
