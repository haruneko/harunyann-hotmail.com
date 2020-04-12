// Copyright 2020 Hal@shurabaP.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utau

import "os"

// Voicebank represents a single UTAU library.
type Voicebank struct {
	Path string `json:"path"`
	PhonemesMap map[string]*Phonemes `json:"phoneme_sets"`
	Character *Character `json:"character"`
	Affixes   *Affixes   `json:"affixes"`
}

// VoicebankReader reads voicebanks from file on the file system.
type VoicebankReader interface {
	Read(string) (*Voicebank, error)
}

// VoicebankReaderDefault is a default VoicebankReader.
type voicebankReaderDefault struct {
	pr PhonemesReader
	cr CharacterReader
	ar AffixesReader
	de directoryEnumerator
}

// VoicebankReaderDefault is a default VoicebankReader.
func (vr voicebankReaderDefault) Read(path string) (*Voicebank, error) {
	ds, err := vr.de.Enumerate(path)
	if err != nil {
		return nil, err
	}
	pm := map[string]*Phonemes { }
	for _, d := range ds {
		if !d.IsDir() {
			continue
		}
		ps, err := vr.pr.Read(resolvePath(resolvePath(path, d.Name()), "oto.ini"))
		if err == nil {
			pm[d.Name()] = ps
		}
	}
	c, err := vr.cr.Read(resolvePath(path, "character.txt"))
	if err != nil {
		c = &Character{}
	}
	a, err := vr.ar.Read(resolvePath(path, "prefix.map"))
	if err != nil {
		a = &Affixes{}
	}
	return &Voicebank {
		Path: path,
		PhonemesMap: pm,
		Character: c,
		Affixes: a,
	}, nil
}

func resolvePath(path string, filename string) string {
	return path + string(os.PathSeparator) + filename
}
