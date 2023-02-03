//go:build linux || freebsd || netbsd || openbsd
// +build linux freebsd netbsd openbsd

package file

import (
	"os"
	"syscall"
)

func (i LInfo) ctime() int64 {
	if i.ok() {
		stat := i.fd.Sys().(*syscall.Stat_t)
		return stat.Ctim.Nsec
	}
	return 0
}

func (i LInfo) atime() int64 {
	if i.ok() {
		stat := i.fd.Sys().(*syscall.Stat_t)
		return stat.Atim.Nsec
	}
	return 0
}

func openFile(filename string) (*os.File, error) {
	return os.Open(filename)
}
