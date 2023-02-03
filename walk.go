package file

import (
	"context"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
)

const (
	Run int = iota + 1
	Init
	Err
)

var (
	walkTypeof = reflect.TypeOf((*walk)(nil)).String()
)

type walk struct {
	lua.SuperVelaData
	cfg    *walkConfig
	output chan LInfo
	ctx    context.Context
	stop   context.CancelFunc
	offset int
	dirs   int32
	files  int32
}

func newWalk(cfg *walkConfig) *walk {
	w := &walk{cfg: cfg}
	return w
}

func (w *walk) Name() string {
	return w.cfg.name
}

func (w *walk) Type() string {
	return walkTypeof
}

func (w *walk) Close() error {
	if w.stop != nil {
		w.stop()
	}

	if w.output != nil {
		close(w.output)
	}

	w.V(lua.VTClose)
	return nil
}

func (w *walk) pretreatment() {
	ctx, stop := context.WithCancel(context.Background())
	w.ctx = ctx
	w.stop = stop
	w.offset = 0
	w.output = make(chan LInfo, 64)
}
func (w *walk) Start() error {
	w.pretreatment()
	xEnv.Spawn(0, w.handle)
	xEnv.Spawn(0, w.scan)
	return nil
}
