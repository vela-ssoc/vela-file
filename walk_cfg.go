package file

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"go.uber.org/ratelimit"
)

var subscript int = 0

type walkConfig struct {
	dir    bool
	deep   int
	name   string
	path   []string
	ignore *cond.Cond
	filter *cond.Cond
	limit  ratelimit.Limiter
	co     *lua.LState
	pipe   *pipe.Chains
	Finish *pipe.Chains
}

//local w = file.walk("/var/logs")
//w.ignore("*.log")
//w.filter("*java*")

func (cfg *walkConfig) append(path string) {
	cfg.path = append(cfg.path, path)
}

func newWalkConfig(L *lua.LState) *walkConfig {
	subscript++

	n := L.GetTop()
	if n == 0 {
		L.RaiseError("not found path")
		return nil
	}

	w := &walkConfig{
		co:   xEnv.Clone(L),
		dir:  false,
		name: fmt.Sprintf("walk.%d", subscript),
		pipe: pipe.New(pipe.Env(xEnv)),
	}

	for i := 1; i <= n; i++ {
		w.append(L.CheckString(i))
	}

	return w
}
