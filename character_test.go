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

func TestSuccessfulCasesOfNewCharacterFromText(t *testing.T) {
	type TestCase struct {
		input    string
		expected *Character
	}
	for i, tc := range []TestCase{
		{"", &Character{Name: "", ImagePath: "", SampleWavePath: "", Author: "", WebURLPath: ""}},
		{"name=テスト音源\nimage=テスト画像\nsample=サンプルパス\nauthor=作成者\nweb=WebサイトのURL", &Character{Name: "テスト音源", ImagePath: "テスト画像", SampleWavePath: "サンプルパス", Author: "作成者", WebURLPath: "WebサイトのURL"}},
		{"name=テスト音源\nimage=テスト画像\nThis line should be ignored.", &Character{Name: "テスト音源", ImagePath: "テスト画像", SampleWavePath: "", Author: "", WebURLPath: ""}},
	} {
		t.Logf("Test case %v.; `%v` should be interpreted as `%v`.", i+1, tc.input, tc.expected)
		actual, err := NewCharacterFromText(tc.input)
		assert.Equal(t, nil, err)
		assert.EqualValues(t, tc.expected, actual)
	}
}

func TestFailedCasesOfNewCharacterFromText(t *testing.T) {
	/* Currently no failed cases exist. */
}

type characterFactoryMock struct {
	mock.Mock
}

func (m *characterFactoryMock) New(text string) (*Character, error) {
	args := m.Called(text)
	return args.Get(0).(*Character), args.Error(1)
}

var testCharacter = &Character{Name: "test", ImagePath: "test.png", SampleWavePath: "test.wav", Author: "test author", WebURLPath: "https://example.com"}

func TestCharacterReaderReadsFileSuccessfully(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text."
	mockedFileReader := new(fileReaderMock)
	mockedCharacterFactory := new(characterFactoryMock)
	sut := &characterReaderDefault{
		fr: mockedFileReader,
		cf: mockedCharacterFactory,
	}
	expected := testCharacter
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedCharacterFactory.On("New", fakeFileText).Return(expected, nil)
	actual, err := sut.Read(testCase)
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expected, actual)
	mockedFileReader.AssertExpectations(t)
	mockedCharacterFactory.AssertExpectations(t)
}

func TestCharacterReaderReadsFileInFailWhenReadingFileFails(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text."
	mockedFileReader := new(fileReaderMock)
	mockedCharacterFactory := new(characterFactoryMock)
	sut := &characterReaderDefault{
		fr: mockedFileReader,
		cf: mockedCharacterFactory,
	}
	expected := errors.New("FAILED")
	mockedFileReader.On("Read", testCase).Return("", expected)
	_, err := sut.Read(testCase)
	assert.Equal(t, expected, err)
	mockedFileReader.AssertExpectations(t)
	mockedCharacterFactory.AssertNumberOfCalls(t, "New", 0)
}

func TestCharacterReaderReadsFileInFailWhenParsingTextFails(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text."
	mockedFileReader := new(fileReaderMock)
	mockedCharacterFactory := new(characterFactoryMock)
	sut := &characterReaderDefault{
		fr: mockedFileReader,
		cf: mockedCharacterFactory,
	}
	expected := errors.New("FAILED")
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedCharacterFactory.On("New", fakeFileText).Return((*Character)(nil), expected)
	_, err := sut.Read(testCase)
	assert.Equal(t, expected, err)
	mockedFileReader.AssertExpectations(t)
	mockedCharacterFactory.AssertExpectations(t)
}
