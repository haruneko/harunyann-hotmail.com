// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"io/ioutil"
	"os"
)

type fileReader interface {
	Read(string) (string, error)
}

type fileReaderDefault struct {
}

func (fr fileReaderDefault) Read(filename string) (string, error) {
	return readFile(filename)
}

func readFile(filename string) (string, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type directoryEnumerator interface {
	Enumerate(string) ([]os.FileInfo, error)
}

type directoryEnumeratorDefault struct {
}

func (de directoryEnumeratorDefault) Enumerate(path string) ([]os.FileInfo, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	res := []os.FileInfo{}
	for _, f := range fs {
		if f.IsDir() {
			res = append(res, f)
		}
	}
	return res, nil
}
