// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package spanish_test

import (
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/gsm7/charset/spanish"
)

type testPattern struct {
	b byte
	r rune
}

var (
	extSize = 20

	// just a random sample for smoke testing
	extTestPatterns = []testPattern{
		{0x0A, '\f'},
		{0x2f, '\\'},
		{0x40, '|'},
		{0x4f, 'Ó'},
		{0x65, '€'},
		{0x6f, 'ó'},
		{0x75, 'ú'},
	}
)

func TestNewExtDecoder(t *testing.T) {
	d := spanish.NewExtDecoder()
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

func TestNewExtEncoder(t *testing.T) {
	e := spanish.NewExtEncoder()
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
