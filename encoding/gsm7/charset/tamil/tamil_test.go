// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tamil_test

import (
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/gsm7/charset/tamil"
)

type testPattern struct {
	b byte
	r rune
}

var (
	mapSize = 103
	extSize = 79

	// just a random sample for smoke testing,
	// particularly characters that differ from basic.
	testPatterns = []testPattern{
		{0x01, '\u0b82'},
		{0x0a, '\n'},
		{0x0d, '\r'},
		{0x0f, '\u0b8f'},
		{0x10, '\u0b90'},
		{0x1b, 0x1b},
		{0x1f, '\u0b9e'},
		{0x20, ' '},
		{0x30, '0'},
		{0x42, '\u0bae'},
		{0x50, '\u0bbe'},
		{0x70, 'p'},
		{0x7b, '\u0bd7'},
		{0x7c, '\u0bf0'},
		{0x7f, '\u0bf9'},
	}

	// just a random sample for smoke testing,
	// particularly characters that differ from basic.
	extTestPatterns = []testPattern{
		{0x00, '@'},
		{0x08, '&'},
		{0x0a, '\f'},
		{0x0b, '*'},
		{0x13, '¡'},
		{0x1c, '\u0be6'},
		{0x1e, '\u0be8'},
		{0x1f, '\u0be9'},
		{0x26, '\u0bf3'},
		{0x27, '\u0bf4'},
		{0x2a, '\u0bf5'},
		{0x2e, '\u0bfa'},
		{0x2f, '\\'},
		{0x40, '|'},
		{0x47, 'G'},
		{0x65, '€'},
	}
)

func TestNewDecoder(t *testing.T) {
	d := tamil.NewDecoder()
	if len(d) != mapSize {
		t.Errorf("expected map length of %d, got %d", mapSize, len(d))
	}
	for _, p := range testPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to encode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}

func TestNewExtDecoder(t *testing.T) {
	d := tamil.NewExtDecoder()
	if len(d) != extSize {
		t.Errorf("expected map length of %d, got %d", extSize, len(d))
	}
	for _, p := range extTestPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to encode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}

func TestNewEncoder(t *testing.T) {
	e := tamil.NewEncoder()
	if len(e) != mapSize {
		t.Errorf("expected map length of %d, got %d", mapSize, len(e))
	}
	for _, p := range testPatterns {
		f := func(t *testing.T) {
			b, ok := e[p.r]
			if !ok {
				t.Errorf("failed to encode '%c': expected 0x%02x, got no match", p.r, p.b)
			} else if b != p.b {
				t.Errorf("failed to encode '%c': expected 0x%02x, got 0x%02x", p.r, p.b, b)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}

func TestNewExtEncoder(t *testing.T) {
	e := tamil.NewExtEncoder()
	if len(e) != extSize-2 { // 2 duplicate values - '*' and '¡', which are mapped to lowest key
		t.Errorf("expected map length of %d, got %d", extSize-2, len(e))
	}
	for _, p := range extTestPatterns {
		f := func(t *testing.T) {
			b, ok := e[p.r]
			if !ok {
				t.Errorf("failed to encode '%c': expected 0x%02x, got no match", p.r, p.b)
			} else if b != p.b {
				t.Errorf("failed to encode '%c': expected 0x%02x, got 0x%02x", p.r, p.b, b)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}
