package file

import (
	"archive/zip"
	"github.com/vela-ssoc/vela-kit/lua"
	"path/filepath"
)

func (w *walk) zip(L *lua.LState) int {
	path := L.IsString(1)
	if path == "" {
		return 0
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		L.Push(lua.S2L(err.Error()))
		return 1
	}
	defer r.Close()

	if len(r.File) == 0 {
		return 0
	}

	for _, f := range r.File {
		fd := f.FileInfo()
		if fd.IsDir() {
			continue
		}

		file := NewLInfo(path+"#"+f.Name, fd, nil)
		file.ext = filepath.Ext(path)
		file.fd = fd

		fi, ei := zip.FileInfoHeader(f.FileInfo())
		xEnv.Errorf("name %s zip LInfo %v , error %v", f.Name, fi, ei)
	}

	return 0
}
