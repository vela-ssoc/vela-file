package file

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
	"reflect"
	"strings"
	"time"
)

var fileTypeOf = reflect.TypeOf((*xFile)(nil)).String()

type xFile struct {
	lua.SuperVelaData
	cfg *config
	fd  *os.File
}

func newFile(cfg *config) *xFile {
	xf := &xFile{cfg: cfg}
	xf.V(fileTypeOf, lua.VTInit)
	return &xFile{cfg: cfg}
}

func (xf *xFile) Name() string {
	return xf.cfg.name
}

func (xf *xFile) Type() string {
	return fileTypeOf
}

// format = access_log.YY-MM-dd.HH-mm-ss
func (xf *xFile) filename(now time.Time) string {
	path := strings.ReplaceAll(xf.cfg.path, "YYYY", "2006")
	path = strings.ReplaceAll(path, "MM", "01")
	path = strings.ReplaceAll(path, "dd", "02")
	path = strings.ReplaceAll(path, "HH", "15")
	path = strings.ReplaceAll(path, "mm", "04")
	return now.Format(path)
}

func (xf *xFile) Start() error {
	filename := xf.filename(time.Now())
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return err
	}
	xf.fd = fd
	return nil
}

func (xf *xFile) Close() error {
	xEnv.Errorf("%s file proc close", xf.Name())
	xf.V(lua.VTClose, time.Now())
	return xf.fd.Close()
}

func (xf *xFile) Write(p []byte) (int, error) {
	if xf.cfg.delim != "" {
		p = append(p, xf.cfg.delim...)
	}

	return xf.fd.Write(p)
}

func (xf *xFile) Push(v interface{}) error {
	_, err := auxlib.Push(xf, v)
	return err
}
