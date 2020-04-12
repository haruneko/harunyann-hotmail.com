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

func TestSuccessfulCasesOfNewAffixesFromText(t *testing.T) {
	type TestCase struct {
		input    string
		expected *Affixes
	}
	for i, tc := range []TestCase{
		{"C4\tweek_\t↓\nC5\tnormal_\t\nC6\tstrong_\t↑\n", &Affixes{"C4": &Affix{"week_", "↓"}, "C5": &Affix{"normal_", ""}, "C6": &Affix{"strong_", "↑"}}},
		{"\n", &Affixes{}},
	} {
		t.Logf("Test case %v.; `%v` should be interpreted as `%v`.", i+1, tc.input, tc.expected)
		actual, err := NewAffixesFromText(tc.input)
		assert.Equal(t, nil, err)
		assert.EqualValues(t, tc.expected, actual)
	}
}

func TestFailedCasesOfNewAffixesFromText(t *testing.T) {
	/* Currently no failed cases exist. */
}

type affixesFactoryMock struct {
	mock.Mock
}

func (m *affixesFactoryMock) New(fn string) (*Affixes, error) {
	args := m.Called(fn)
	return args.Get(0).(*Affixes), args.Error(1)
}

var testAffixes = &Affixes{"C4": &Affix{Prefix: "weak_", Suffix: "↓"}, "C6": &Affix{Prefix: "strong_", Suffix: "↑"}}

func TestAffixesReaderReadsFileSuccessfully(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text code."
	mockedFileReader := new(fileReaderMock)
	mockedAffixesFactory := new(affixesFactoryMock)
	sut := &affixesReaderDefault{
		fr: mockedFileReader,
		af: mockedAffixesFactory,
	}
	expected := testAffixes
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedAffixesFactory.On("New", fakeFileText).Return(expected, nil)
	actual, err := sut.Read(testCase)
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expected, actual)
	mockedFileReader.AssertExpectations(t)
	mockedAffixesFactory.AssertExpectations(t)
}

func TestAffixesReaderReadsFileInFailWhenReadingFileFails(t *testing.T) {
	const testCase = "testCase"
	err := errors.New("test error")
	mockedFileReader := new(fileReaderMock)
	mockedAffixesFactory := new(affixesFactoryMock)
	sut := &affixesReaderDefault{
		fr: mockedFileReader,
		af: mockedAffixesFactory,
	}
	mockedFileReader.On("Read", testCase).Return("", err)
	_, e := sut.Read(testCase)
	assert.Equal(t, err, e)
	mockedFileReader.AssertExpectations(t)
	mockedAffixesFactory.AssertNumberOfCalls(t, "New", 0)
}

func TestAffixesReaderReadsFileInFailWhenParsingTextFails(t *testing.T) {
	const testCase = "testCase"
	const fakeFileText = "This is a fake text code."
	err := errors.New("test error")
	mockedFileReader := new(fileReaderMock)
	mockedAffixesFactory := new(affixesFactoryMock)
	sut := &affixesReaderDefault{
		fr: mockedFileReader,
		af: mockedAffixesFactory,
	}
	mockedFileReader.On("Read", testCase).Return(fakeFileText, nil)
	mockedAffixesFactory.On("New", fakeFileText).Return((*Affixes)(nil), err)
	_, e := sut.Read(testCase)
	assert.Equal(t, err, e)
	mockedFileReader.AssertExpectations(t)
	mockedAffixesFactory.AssertExpectations(t)
}
