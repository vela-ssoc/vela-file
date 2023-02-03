package file

import (
	"os"
	"testing"
)

func TestDirMtime(t *testing.T) {
	v, err := os.Stat("../rule.d")
	if err != nil {
		t.Logf("not found %v", err)
		return
	}

	t.Logf("mtime:%d size:%d", v.ModTime().Unix(), v.Size())
}
