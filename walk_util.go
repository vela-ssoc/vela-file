package file

import "os"

func match(path string, p []func(string) bool) bool {
	n := len(p)
	if n == 0 {
		return false
	}

	for i := 0; i < n; i++ {
		if p[i](path) {
			return true
		}
	}

	return false
}

func depth(path string) int {
	n := 1
	size := len(path)
	for i := 0; i < size; i++ {
		if path[i] == os.PathSeparator {
			n++
		}
	}
	return n
}
