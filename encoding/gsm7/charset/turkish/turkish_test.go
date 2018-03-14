// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package turkish_test

import (
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/gsm7/charset/turkish"
)

type testPattern struct {
	b byte
	r rune
}

var (
	mapSize = 128
	extSize = 18

	// just a random sample for smoke testing,
	// particularly characters that differ from basic.
	testPatterns = []testPattern{
		{0x00, '@'},
		{0x0a, '\n'},
		{0x0d, '\r'},
		{0x0f, 'å'},
		{0x10, 'Δ'},
		{0x1b, '\x1b'},
		{0x1c, 'Ş'},
		{0x1f, 'É'},
		{0x20, ' '},
		{0x30, '0'},
		{0x41, 'A'},
		{0x50, 'P'},
		{0x70, 'p'},
		{0x7b, 'ä'},
		{0x7c, 'ö'},
		{0x7f, 'à'},
	}

	// just a random sample for smoke testing,
	// particularly characters that differ from basic.
	extTestPatterns = []testPattern{
		{0x0A, '\f'},
		{0x2f, '\\'},
		{0x40, '|'},
		{0x47, 'Ğ'},
		{0x65, '€'},
		{0x67, 'ğ'},
		{0x69, 'ı'},
		{0x73, 'ş'},
	}
)

func TestNewDecoder(t *testing.T) {
	d := turkish.NewDecoder()
	if testing.Verbose() {
		charset.Display(d)
	}
	if len(d) != mapSize {
		t.Errorf("expected map length of %d, got %d", mapSize, len(d))
	}
	for _, p := range testPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}

func TestNewExtDecoder(t *testing.T) {
	d := turkish.NewExtDecoder()
	if testing.Verbose() {
		charset.Display(d)
	}
	if len(d) != extSize {
		t.Errorf("expected map length of %d, got %d", extSize, len(d))
	}
	for _, p := range extTestPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("x%02x", p.b), f)
	}
}

func TestNewEncoder(t *testing.T) {
	e := turkish.NewEncoder()
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
	e := turkish.NewExtEncoder()
	if len(e) != extSize {
		t.Errorf("expected map length of %d, got %d", extSize, len(e))
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
