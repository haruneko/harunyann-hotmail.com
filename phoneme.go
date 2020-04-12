// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"errors"
	"strconv"
	"strings"
)

// Phoneme represents one piece of wave file that contains voice.
type Phoneme struct {
	Filename     string  `json:"filename"`
	Alias        string  `json:"alias"`
	LeftBlank    float64 `json:"left_blank"`
	Consonant    float64 `json:"consonant"`
	RightBlank   float64 `json:"right_blank"`
	PreUtterance float64 `json:"pre_utterance"`
	Overlap      float64 `json:"overlap"`
}

// Phonemes represents a single oto.ini in UTAU library.
type Phonemes []*Phoneme

// NewPhonemesFromText creates Phonemes from a single oto.ini.
func NewPhonemesFromText(t string) (*Phonemes, error) {
	res := Phonemes{}
	for _, l := range strings.Split(t, "\n") {
		p, err := NewPhonemeFromLine(l)
		if err != nil {
			// Log here if you want.
			continue
		}
		res = append(res, p)
	}
	return &res, nil
}

// NewPhonemeFromLine creates Phoneme from a single line in a single oto.ini.
func NewPhonemeFromLine(l string) (*Phoneme, error) {
	es := strings.Split(l, ",")
	if len(es) != 6 {
		return nil, errors.New("The given line does not contain 6 elements; `" + l + "`")
	}
	if !strings.Contains(es[0], "=") {
		return nil, errors.New("The given line does not contain alias; `" + es[0] + "`")
	}
	fa := strings.Split(es[0], "=")
	if len(fa) != 2 {
		return nil, errors.New("The given line does not contain valid alias; `" + es[0] + "`")
	}
	r1, r2, r3, r4, r5, err := parse(es)
	if err != nil {
		return nil, err
	}
	return &Phoneme{
		Filename:     fa[0],
		Alias:        fa[1],
		LeftBlank:    r1,
		Consonant:    r2,
		RightBlank:   r3,
		PreUtterance: r4,
		Overlap:      r5,
	}, nil
}

func parse(es []string) (float64, float64, float64, float64, float64, error) {
	r1, err := strconv.ParseFloat(es[1], 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	r2, err := strconv.ParseFloat(es[2], 64)
	if err != nil {
		return r1, 0, 0, 0, 0, err
	}
	r3, err := strconv.ParseFloat(es[3], 64)
	if err != nil {
		return r1, r2, 0, 0, 0, err
	}
	r4, err := strconv.ParseFloat(es[4], 64)
	if err != nil {
		return r1, r2, r3, 0, 0, err
	}
	r5, err := strconv.ParseFloat(es[5], 64)
	if err != nil {
		return r1, r2, r3, r4, 0, err
	}
	return r1, r2, r3, r4, r5, nil
}

type phonemesFactory interface {
	New(string) (*Phonemes, error)
}

type phonemesFactoryDefault struct {
}

func (pf phonemesFactoryDefault) New(filename string) (*Phonemes, error) {
	return NewPhonemesFromText(filename)
}

// PhonemesReader reads Phonemes from the file on file system.
type PhonemesReader interface {
	Read(string) (*Phonemes, error)
}

// PhonemesReaderDefault is a default PhonemeReader.
type PhonemesReaderDefault struct {
	fr fileReader
	pf phonemesFactory
}

// NewPhonemesReader create PhonemesReader that uses file system as its source.
func NewPhonemesReader() PhonemesReader {
	return PhonemesReaderDefault{
		fr: fileReaderDefault{},
		pf: phonemesFactoryDefault{},
	}
}

func (pr PhonemesReaderDefault) Read(filename string) (*Phonemes, error) {
	t, err := pr.fr.Read(filename)
	if err != nil {
		return nil, err
	}
	return pr.pf.New(t)
}
