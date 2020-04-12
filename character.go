// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import (
	"strings"
)

// Character reprensents character.txt.
type Character struct {
	Name           string `json:"name"`
	ImagePath      string `json:"image_path"`
	SampleWavePath string `json:"sample_wave_path"`
	Author         string `json:"author"`
	WebURLPath     string `json:"web_url_path"`
}

// NewCharacterFromText creates Character from a text.
func NewCharacterFromText(text string) (*Character, error) {
	name, imagePath, sampleWavePath, author, webURLPath := "", "", "", "", ""
	for _, l := range strings.Split(text, "\n") {
		es := strings.Split(l, "=")
		if len(es) != 2 {
			continue
		}
		switch es[0] {
		case "name":
			name = es[1]
		case "image":
			imagePath = es[1]
		case "sample":
			sampleWavePath = es[1]
		case "author":
			author = es[1]
		case "web":
			webURLPath = es[1]
		}
	}
	return &Character{
		Name: name, ImagePath: imagePath, SampleWavePath: sampleWavePath, Author: author, WebURLPath: webURLPath,
	}, nil
}

type characterFactory interface {
	New(string) (*Character, error)
}

type characterFactoryDefault struct {
}

func (cf characterFactoryDefault) New(text string) (*Character, error) {
	return NewCharacterFromText(text)
}

// CharacterReader reads Character from filesystem.
type CharacterReader interface {
	Read(string) (*Character, error)
}

// CharacterReaderDefault is a default CharacterReader.
type characterReaderDefault struct {
	fr fileReader
	cf characterFactory
}

// NewCharacterReader creates a default CharacterReader that reads Character from filesystem.
func NewCharacterReader() CharacterReader {
	return characterReaderDefault{
		fr: fileReaderDefault{},
		cf: characterFactoryDefault{},
	}
}

// Read Character from file.
// Usually CharacterReader reads Character from file system.
func (cr characterReaderDefault) Read(filename string) (*Character, error) {
	t, err := cr.fr.Read(filename)
	if err != nil {
		return nil, err
	}
	return cr.cf.New(t)
}
