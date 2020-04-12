// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type phonemesReaderMock struct {
	mock.Mock
}

func (m *phonemesReaderMock) Read(fn string) (*Phonemes, error) {
	args := m.Called(fn)
	return args.Get(0).(*Phonemes), args.Error(1)
}

type characterReaderMock struct {
	mock.Mock
}

func (m *characterReaderMock) Read(fn string) (*Character, error) {
	args := m.Called(fn)
	return args.Get(0).(*Character), args.Error(1)
}

type affixesReaderMock struct {
	mock.Mock
}

func (m *affixesReaderMock) Read(fn string) (*Affixes, error) {
	args := m.Called(fn)
	return args.Get(0).(*Affixes), args.Error(1)
}

var testVoicebank = &Voicebank{
	Path: "fake_filepath",
	PhonemesMap: map[string]*Phonemes{
		"": testPhonemes,
		"hoge": testPhonemes,
		"foo": testPhonemes,
		"bar": &Phonemes{},
	},
	Character: testCharacter,
	Affixes: testAffixes,
}

var dummyFileInfos = []os.FileInfo {
	dummyFileInfo{n: "", s: 0, m: os.ModeDir, t: time.Date(1985, 2, 4, 0, 0, 0, 0, time.UTC)},
	dummyFileInfo{n: "hoge", s: 0, m: os.ModeDir, t: time.Date(1985, 2, 4, 0, 0, 0, 0, time.UTC)},
	dummyFileInfo{n: "foo", s: 0, m: os.ModeDir, t: time.Date(1985, 2, 4, 0, 0, 0, 0, time.UTC)},
	dummyFileInfo{n: "bar", s: 0, m: os.ModeDir, t: time.Date(1985, 2, 4, 0, 0, 0, 0, time.UTC)},
	dummyFileInfo{n: "not_a_directory.txt",s: 0, m: os.ModeDevice, t: time.Date(1985, 2, 4, 0, 0, 0, 0, time.UTC)},
}

func TestVoicebankReaderReadsFileSuccessfully(t *testing.T) {
	const path = "fake_filepath"
	mockedDirectoryEnumerator := new(directoryEnumeratorMock)
	mockedPhonemesReader := new(phonemesReaderMock)
	mockedCharacterReader := new(characterReaderMock)
	mockedAffixesReader := new(affixesReaderMock)
	sut := &voicebankReaderDefault{
		pr: mockedPhonemesReader,
		cr: mockedCharacterReader,
		ar: mockedAffixesReader,
		de: mockedDirectoryEnumerator,
	}
	expected := testVoicebank
	mockedDirectoryEnumerator.On("Enumerate", path).Return(dummyFileInfos, nil)
	mockedPhonemesReader.On("Read", resolvePath(resolvePath(path, ""), "oto.ini")).Return(testPhonemes, nil)
	mockedPhonemesReader.On("Read", resolvePath(resolvePath(path, "hoge"), "oto.ini")).Return(testPhonemes, nil)
	mockedPhonemesReader.On("Read", resolvePath(resolvePath(path, "foo"), "oto.ini")).Return(testPhonemes, nil)
	mockedPhonemesReader.On("Read", resolvePath(resolvePath(path, "bar"), "oto.ini")).Return(&Phonemes{}, nil)
	mockedCharacterReader.On("Read", resolvePath(path, "character.txt")).Return(testCharacter, nil)
	mockedAffixesReader.On("Read", resolvePath(path, "prefix.map")).Return(testAffixes, nil)
	actual, err := sut.Read(path)
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expected, actual)
	mockedAffixesReader.AssertExpectations(t)
	mockedCharacterReader.AssertExpectations(t)
	mockedDirectoryEnumerator.AssertExpectations(t)
	mockedPhonemesReader.AssertExpectations(t)
}
