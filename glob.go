package file

import (
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"os"
	"path/filepath"
)

type glob struct {
	co       *lua.LState
	patterns []string
	pipe     *pipe.Chains
	result   []string
	err      *exception.Cause
}

func (gl *glob) visit(path string) {
	stat, err := os.Stat(path)
	fi := NewLInfo(path, stat, err)

	gl.pipe.Do(fi, gl.co, func(err error) {
		xEnv.Errorf("%s filepath glob pipe call fail %v", path, err)
	})
}

func (gl *glob) evaluate() {
	if gl.pipe.Len() == 0 {
		return
	}

	rn := len(gl.result)
	if rn == 0 {
		return
	}

	for i := 0; i < rn; i++ {
		gl.visit(gl.result[i])
	}
}

func (gl *glob) run() {
	n := len(gl.patterns)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		pattern := gl.patterns[i]
		v, err := filepath.Glob(pattern)
		if err != nil {
			gl.err.Try(pattern, err)
			continue
		}
		gl.result = append(gl.result, v...)
	}

	gl.evaluate()
}
