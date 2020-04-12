// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"strings"
)

// Affix represents a single line on prefix.map of UTAU.
type Affix struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

// Affixes represents a prefix.map of UTAU
type Affixes map[string]*Affix

// NewAffixesFromText creates Affixes from the text on prefix.map.
func NewAffixesFromText(text string) (*Affixes, error) {
	res := Affixes{}
	for _, l := range strings.Split(text, "\n") {
		es := strings.Split(l, "\t")
		if len(es) != 3 {
			continue
		}
		res[es[0]] = &Affix{Prefix: es[1], Suffix: es[2]}
	}
	return &res, nil
}

type affixesFactory interface {
	New(string) (*Affixes, error)
}

type affixesFactoryDefault struct {
}

func (af affixesFactoryDefault) New(text string) (*Affixes, error) {
	return NewAffixesFromText(text)
}

// AffixesReader reads Affixes from filesystem.
type AffixesReader interface {
	Read(string) (*Affixes, error)
}

// AffixesReaderDefault reads Affixes from filesystem.
type affixesReaderDefault struct {
	fr fileReader
	af affixesFactory
}

// NewAffixesReaderDefault creates a default AffixesReader that reads Affixes from filesystem.
func NewAffixesReaderDefault() AffixesReader {
	return affixesReaderDefault{
		fr: fileReaderDefault{},
		af: affixesFactoryDefault{},
	}
}

// Read Affixes from the file specified by filename.
func (ar affixesReaderDefault) Read(filename string) (*Affixes, error) {
	t, err := ar.fr.Read(filename)
	if err != nil {
		return nil, err
	}
	return ar.af.New(t)
}
