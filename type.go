package file

import (
	"github.com/h2non/filetype"
	"os"
)

type Type struct {
	Err       error  `json:"err"`
	Extension string `json:"extension"`
	Typ       string `json:"type"`
	Subtype   string `json:"subtype"`
	Value     string `json:"value"`
}

func Filetype(path string) *Type {
	ft := &Type{}
	file, err := os.Open(path)
	if err != nil {
		ft.Err = err
		return ft
	}
	defer file.Close()

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		ft.Err = err
		return ft
	}

	kind, err := filetype.Match(head)
	if err != nil {
		ft.Err = err
		return ft
	}

	ft.Typ = kind.MIME.Type
	ft.Extension = kind.Extension
	ft.Subtype = kind.MIME.Subtype
	ft.Value = kind.MIME.Value
	return ft
}
