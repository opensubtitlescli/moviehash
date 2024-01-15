package testdata

import (
	"fmt"
	"path"
	"runtime"
)

var filesMap = [][]string{
	{"breakdance.avi", "8e245d9679d31e12"},
	{"dummy.bin",      "61f7751fc2a72bfb"},
}

func FilesMap() ([][]string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("runtime.Caller(0) failed")
	}
	d := path.Dir(file)
	m := [][]string{}
	for _, v := range filesMap {
		f := path.Join(d, v[0])
		m = append(m, []string{f, v[1]})
	}
	return m, nil
}
