package file

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"go.uber.org/ratelimit"
	"time"
)

//func (w *walk) String() string                         { return fmt.Sprintf("%p", w) }
//func (w *walk) Typ() lua.LValueType                   { return lua.LTObject }
//func (w *walk) AssertFloat64() (float64, bool)         { return 0, false }
//func (w *walk) AssertString() (string, bool)           { return "", false }
//func (w *walk) AssertFunction() (*lua.LFunction, bool) { return nil, false }
//func (w *walk) Peek() lua.LValue                       { return w }

//func (w *walk) setExt(L *lua.LState) int {
//	n := L.GetTop()
//	if n == 0 {
//		return 0
//	}
//
//	w.ext = make([]string, n)
//	for i := 1; i <= n; i++ {
//		w.ext[i-1] = L.CheckString(i)
//	}
//	return 0
//}
//
//func (w *walk) setNotExt(L *lua.LState) int {
//	n := L.GetTop()
//	if n == 0 {
//		return 0
//	}
//
//	w.notExt = make([]string, n)
//	for i := 1; i <= n; i++ {
//		w.ext[i-1] = L.CheckString(i)
//	}
//	return 0
//}

func (w *walk) dirL(L *lua.LState) int {
	w.cfg.dir = L.IsTrue(1)
	return 0
}

func (w *walk) limitL(L *lua.LState) int {
	rv := L.IsInt(1)
	pre := L.IsInt(2)
	if rv <= 0 {
		return 0
	}

	if pre <= 0 {
		w.cfg.limit = ratelimit.New(rv)
	} else {
		pv := time.Duration(pre) * time.Second
		w.cfg.limit = ratelimit.New(rv, ratelimit.Per(pv))
	}
	return 0
}

func (w *walk) run(L *lua.LState) int {
	xEnv.Start(L, w).From(L.CodeVM()).Do()
	return 0
}

func (w *walk) scanL(L *lua.LState) int {
	w.pretreatment()
	xEnv.Spawn(0, w.handle)
	w.scan()
	return 0
}

func (w *walk) ignoreL(L *lua.LState) int {
	if w.cfg.ignore == nil {
		w.cfg.ignore = cond.CheckMany(L)
	} else {
		w.cfg.ignore.CheckMany(L)
	}
	return 0
}

func (w *walk) filterL(L *lua.LState) int {
	if w.cfg.filter == nil {
		w.cfg.filter = cond.CheckMany(L)
	} else {
		w.cfg.filter.CheckMany(L)
	}
	return 0
}

func (w *walk) pipeL(L *lua.LState) int {
	w.cfg.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (w *walk) deepL(L *lua.LState) int {
	n := L.IsInt(1)
	if n == 0 {
		return 0
	}

	w.cfg.deep = n
	return 0
}

func (w *walk) onFinishL(L *lua.LState) int {
	w.cfg.Finish = pipe.NewByLua(L)
	return 0
}

func (w *walk) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "dir":
		return L.NewFunction(w.dirL)
	case "deep":
		return L.NewFunction(w.deepL)

	case "ignore":
		return lua.NewFunction(w.ignoreL)
	case "filter":
		return lua.NewFunction(w.filterL)
	case "pipe":
		return lua.NewFunction(w.pipeL)
	//case "ext":
	//	return L.NewFunction(w.setExt)
	//case "not_ext":
	//	return L.NewFunction(w.setNotExt)
	//case "wait":
	//	return L.NewFunction(w.wait)

	case "limit":
		return L.NewFunction(w.limitL)
	case "run":
		return L.NewFunction(w.run)
	case "scan":
		return L.NewFunction(w.scanL)
	case "zip":
		return L.NewFunction(w.zip)
	case "on_finish":
		return L.NewFunction(w.onFinishL)

	}
	return lua.LNil
}
