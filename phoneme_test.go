// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccessfulCasesOfNewPhonemFromLine(t *testing.T) {
	type TestCase struct {
		input    string
		expected *Phoneme
	}
	for i, tc := range []TestCase{
		{"ファイル名=エイリアス,0,1,2,4,8", &Phoneme{Filename: "ファイル名", Alias: "エイリアス", LeftBlank: 0, Consonant: 1, RightBlank: 2, PreUtterance: 4, Overlap: 8}},
		{"_ああいあうえあ.wav=- あ,500,125,-500,500,250", &Phoneme{Filename: "_ああいあうえあ.wav", Alias: "- あ", LeftBlank: 500, Consonant: 125, RightBlank: -500, PreUtterance: 500, Overlap: 250}},
	} {
		t.Logf("Test case %v.; `%v` should be interpreted as `%v`.", i+1, tc.input, tc.expected)
		actual, err := NewPhonemeFromLine(tc.input)
		assert.Equal(t, nil, err)
		assert.EqualValues(t, tc.expected, actual)
	}
}

func TestFailedCasesOfNewPhonemeFromLine(t *testing.T) {
	for i, tc := range []string{
		"This is invalid line for oto.ini",
		"LackOf=Parameters",
		"LackOf=Parameters,0",
		"LackOf=Parameters,0,1,2,4",
		"TooMany=Parameters,0,1,2,4,8,16",
		"NeedsAlias,0,1,2,4,8",
		"ShouldBe=Float,0,1,2,4,NoImString:D",
	} {
		t.Logf("Test case %v.; `%v` cannnot be interpreted as Phoneme.", i+1, tc)
		_, err := NewPhonemeFromLine(tc)
		assert.Error(t, err)
	}
}

func TestSuccessfulCasesOfNewPhonemsFromText(t *testing.T) {
	const input = "ファイル名=エイリアス,0,1,2,4,8\n_ああいあうえあ.wav=- あ,500,125,-500,500,250\n"
	expected := &Phonemes{
		&Phoneme{Filename: "ファイル名", Alias: "エイリアス", LeftBlank: 0, Consonant: 1, RightBlank: 2, PreUtterance: 4, Overlap: 8},
		&Phoneme{Filename: "_ああいあうえあ.wav", Alias: "- あ", LeftBlank: 500, Consonant: 125, RightBlank: -500, PreUtterance: 500, Overlap: 250},
	}
	actual, err := NewPhonemesFromText(input)
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expected, actual)
}

type phonemesFactoryMock struct {
	mock.Mock
}

func (m *phonemesFactoryMock) New(text string) (*Phonemes, error) {
	args := m.Called(text)
	return args.Get(0).(*Phonemes), args.Error(1)
}

var testPhonemes = &Phonemes{
	&Phoneme{Filename: "ファイル名", Alias: "エイリアス", LeftBlank: 0, Consonant: 1, RightBlank: 2, PreUtterance: 4, Overlap: 8},
	&Phoneme{Filename: "_ああいあうえあ.wav", Alias: "- あ", LeftBlank: 500, Consonant: 125, RightBlank: -500, PreUtterance: 500, Overlap: 250},
}

func TestPhonemesReaderReadsFileSuccessfully(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text."
	mockedFileReader := new(fileReaderMock)
	mockedPhonemesFactory := new(phonemesFactoryMock)
	sut := &PhonemesReaderDefault{
		fr: mockedFileReader,
		pf: mockedPhonemesFactory,
	}
	expected := testPhonemes
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedPhonemesFactory.On("New", fakeFileText).Return(expected, nil)
	actual, err := sut.Read(testCase)
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expected, actual)
	mockedFileReader.AssertExpectations(t)
	mockedPhonemesFactory.AssertExpectations(t)
}

func TestPhonemesReaderReadsFileInFailWhenReadingFileFails(t *testing.T) {
	const testCase = "testCase"
	mockedFileReader := new(fileReaderMock)
	mockedPhonemesFactory := new(phonemesFactoryMock)
	sut := &PhonemesReaderDefault{
		fr: mockedFileReader,
		pf: mockedPhonemesFactory,
	}
	expected := errors.New("FAILED")
	mockedFileReader.On("Read", testCase).Return("", expected)
	_, err := sut.Read(testCase)
	assert.Equal(t, expected, err)
	mockedFileReader.AssertExpectations(t)
	mockedPhonemesFactory.AssertNumberOfCalls(t, "New", 0)
}

func TestPhonemesReaderReadsFileInFailWhenParsingTextFails(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text."
	mockedFileReader := new(fileReaderMock)
	mockedPhonemesFactory := new(phonemesFactoryMock)
	sut := &PhonemesReaderDefault{
		fr: mockedFileReader,
		pf: mockedPhonemesFactory,
	}
	expected := errors.New("FAILED")
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedPhonemesFactory.On("New", fakeFileText).Return((*Phonemes)(nil), expected)
	_, err := sut.Read(testCase)
	assert.Equal(t, expected, err)
	mockedFileReader.AssertExpectations(t)
	mockedPhonemesFactory.AssertExpectations(t)
}
