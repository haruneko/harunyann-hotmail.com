// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

type fileReaderMock struct {
	mock.Mock
}

func (m *fileReaderMock) Read(fn string) (string, error) {
	args := m.Called(fn)
	return args.Get(0).(string), args.Error(1)
}

type directoryEnumeratorMock struct {
	mock.Mock
}

func (m *directoryEnumeratorMock) Enumerate(path string) ([]os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

type dummyFileInfo struct {
	n string
	s int64
	m os.FileMode
	t time.Time
}

func (d dummyFileInfo) Name() string       { return d.n }
func (d dummyFileInfo) Size() int64        { return d.s }
func (d dummyFileInfo) Mode() os.FileMode  { return d.m }
func (d dummyFileInfo) ModTime() time.Time { return d.t }
func (d dummyFileInfo) IsDir() bool        { return d.m.IsDir() }
func (d dummyFileInfo) Sys() interface{}   { return nil }
