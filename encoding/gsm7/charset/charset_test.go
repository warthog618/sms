// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package charset_test

import (
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/gsm7/charset"
)

type testPattern struct {
	b byte
	r rune
}

var charsetName = []string{
	"Default",
	"Turkish",
	"Spanish",
	"Portuguese",
	"Bengali",
	"Gujarati",
	"Hindi",
	"Kannada",
	"Malayalam",
	"Oriya",
	"Punjabi",
	"Tamil",
	"Telugu",
	"Urdu",
}

type language struct {
	nli             int
	mapSize         int
	extSize         int
	extDuplicates   int
	lockingPatterns []testPattern
	shiftPatterns   []testPattern
}

func TestDecoder(t *testing.T) {
	for _, l := range languageTests {
		d := charset.NewDecoder(l.nli)
		if len(d) != l.mapSize {
			t.Errorf("expected %s map length of %d, got %d", charsetName[l.nli], l.mapSize, len(d))
		}
		for _, p := range l.lockingPatterns {
			f := func(t *testing.T) {
				r, ok := d[p.b]
				if !ok {
					t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
				} else if r != p.r {
					t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
				}
			}
			t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
		}
	}
}

func TestExtDecoder(t *testing.T) {
	for _, l := range languageTests {
		d := charset.NewExtDecoder(l.nli)
		if len(d) != l.extSize {
			t.Errorf("expected %s map length of %d, got %d", charsetName[l.nli], l.extSize, len(d))
		}
		for _, p := range l.shiftPatterns {
			f := func(t *testing.T) {
				r, ok := d[p.b]
				if !ok {
					t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
				} else if r != p.r {
					t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
				}
			}
			t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
		}
	}
}

func TestEncoder(t *testing.T) {
	for _, l := range languageTests {
		e := charset.NewEncoder(l.nli)
		if len(e) != l.mapSize {
			t.Errorf("expected %s map length of %d, got %d", charsetName[l.nli], l.mapSize, len(e))
		}
		for _, p := range l.lockingPatterns {
			f := func(t *testing.T) {
				b, ok := e[p.r]
				if !ok {
					t.Errorf("failed to encode '%c': expected 0x%02x, got no match", p.r, p.b)
				} else if b != p.b {
					t.Errorf("failed to encode '%c': expected 0x%02x, got 0x%02x", p.r, p.b, b)
				}
			}
			t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
		}
	}
}

func TestExtEncoder(t *testing.T) {
	for _, l := range languageTests {
		e := charset.NewExtEncoder(l.nli)
		if len(e) != l.extSize-l.extDuplicates {
			t.Errorf("expected %s map length of %d, got %d", charsetName[l.nli], l.extSize-l.extDuplicates, len(e))
		}
		for _, p := range l.shiftPatterns {
			f := func(t *testing.T) {
				b, ok := e[p.r]
				if !ok {
					t.Errorf("failed to encode '%c': expected 0x%02x, got no match", p.r, p.b)
				} else if b != p.b {
					t.Errorf("failed to encode '%c': expected 0x%02x, got 0x%02x", p.r, p.b, b)
				}
			}
			t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
		}
	}
}

func TestDefaultDecoder(t *testing.T) {
	d := charset.DefaultDecoder()
	l := languageTests[0]
	if len(d) != l.mapSize {
		t.Errorf("expected map length of %d, got %d", l.mapSize, len(d))
	}
	for _, p := range l.lockingPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
	}
}

func TestDefaultExtDecoder(t *testing.T) {
	d := charset.DefaultExtDecoder()
	l := languageTests[0]
	if len(d) != l.extSize {
		t.Errorf("expected map length of %d, got %d", l.extSize, len(d))
	}
	for _, p := range l.shiftPatterns {
		f := func(t *testing.T) {
			r, ok := d[p.b]
			if !ok {
				t.Errorf("failed to decode 0x%02x: expected '%c', got no match", p.b, p.r)
			} else if r != p.r {
				t.Errorf("failed to decode 0x%02x: expected '%c', got '%c'", p.b, p.r, r)
			}
		}
		t.Run(fmt.Sprintf("%s x%02x", charsetName[l.nli], p.b), f)
	}
}

func TestDefaultEncoder(t *testing.T) {
	e := charset.DefaultEncoder()
	l := languageTests[0]
	if len(e) != l.mapSize {
		t.Errorf("expected map length of %d, got %d", l.mapSize, len(e))
	}
	for _, p := range l.lockingPatterns {
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

func TestDefaultExtEncoder(t *testing.T) {
	e := charset.DefaultExtEncoder()
	l := languageTests[0]
	if len(e) != l.extSize-l.extDuplicates {
		t.Errorf("expected %s map length of %d, got %d", charsetName[l.nli], l.extSize-l.extDuplicates, len(e))
	}
	for _, p := range l.shiftPatterns {
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

var (
	languageTests = []language{
		{charset.Default, 128, 11, 0,
			[]testPattern{
				{0x00, '@'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x10, 'Δ'},
				{0x1b, '\x1b'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, 'A'},
				{0x50, 'P'},
				{0x70, 'p'},
				{0x7f, 'à'},
			},
			[]testPattern{
				{0x0A, '\f'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x65, '€'},
			},
		},
		{charset.Turkish, 128, 18, 0,
			[]testPattern{
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
			},
			[]testPattern{
				{0x0A, '\f'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'Ğ'},
				{0x65, '€'},
				{0x67, 'ğ'},
				{0x69, 'ı'},
				{0x73, 'ş'},
			},
		},
		{charset.Spanish, 128, 20, 0,
			[]testPattern{
				{0x00, '@'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x10, 'Δ'},
				{0x1b, '\x1b'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, 'A'},
				{0x50, 'P'},
				{0x70, 'p'},
				{0x7f, 'à'},
			},
			[]testPattern{
				{0x0A, '\f'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x4f, 'Ó'},
				{0x65, '€'},
				{0x6f, 'ó'},
				{0x75, 'ú'},
			},
		},
		{charset.Portuguese, 128, 37, 0,
			[]testPattern{
				{0x00, '@'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x10, 'Δ'},
				{0x1b, '\x1b'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, 'A'},
				{0x50, 'P'},
				{0x70, 'p'},
				{0x7b, 'ã'},
				{0x7e, 'ü'},
				{0x7f, 'à'},
			},
			[]testPattern{
				{0x0A, '\f'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x4f, 'Ó'},
				{0x65, '€'},
				{0x6f, 'ó'},
				{0x75, 'ú'},
				{0x7b, 'ã'},
				{0x7c, 'õ'},
				{0x7f, 'â'},
			},
		},
		{charset.Bengali, 115, 84, 2, // 2 duplicate values - '*' and '¡', which are mapped to lowest key
			[]testPattern{
				{0x00, '\u0981'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0f, '\u098f'},
				{0x10, '\u0990'},
				{0x1b, 0x1b},
				{0x1f, '\u099e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u09ad'},
				{0x50, '\u09be'},
				{0x70, 'p'},
				{0x7b, '\u09d7'},
				{0x7c, '\u09dc'},
				{0x7f, '\u09f1'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1d, '\u09e9'},
				{0x1e, '\u09ea'},
				{0x1f, '\u09eb'},
				{0x2f, '\\'},
				{0x30, '\u09f6'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Gujaranti, 121, 72, 2, // 2 duplicate values - '*' and '¡', which are mapped to lowest key
			[]testPattern{
				{0x00, '\u0a81'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0f, '\u0a8f'},
				{0x10, '\u0a90'},
				{0x1b, 0x1b},
				{0x1f, '\u0a9e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0aad'},
				{0x50, '\u0abe'},
				{0x70, 'p'},
				{0x7b, '\u0ae0'},
				{0x7c, '\u0ae1'},
				{0x7f, '\u0af1'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1d, '\u0ae7'},
				{0x1e, '\u0ae8'},
				{0x1f, '\u0ae9'},
				{0x25, '\u0aef'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Hindi, 128, 90, 2, // 2 duplicate values - '*' and '¡', which are mapped to lowest key
			[]testPattern{
				{0x00, '\u0981'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0e, 'ऎ'},
				{0x0f, 'ए'},
				{0x10, 'ऐ'},
				{0x1b, 0x1b},
				{0x1c, 'छ'},
				{0x1f, 'ञ'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, 'भ'},
				{0x50, '\u093e'},
				{0x70, 'p'},
				{0x7b, 'ॲ'},
				{0x7c, 'ॻ'},
				{0x7f, 'ॿ'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1d, '१'},
				{0x1e, '२'},
				{0x1f, '३'},
				{0x2f, '\\'},
				{0x30, 'ज़'},
				{0x3a, 'ॱ'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Kannada, 121, 75, 2,
			[]testPattern{
				{0x01, '\u0c82'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0f, '\u0c8f'},
				{0x10, '\u0c90'},
				{0x1b, 0x1b},
				{0x1f, '\u0c9e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0cad'},
				{0x50, '\u0cbe'},
				{0x70, 'p'},
				{0x7b, '\u0cd6'},
				{0x7c, '\u0ce0'},
				{0x7f, '\u0ce3'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1c, '\u0ce6'},
				{0x1e, '\u0ce8'},
				{0x1f, '\u0ce9'},
				{0x26, '\u0cde'},
				{0x27, '\u0cf1'},
				{0x2a, '\u0cf2'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Malayalam, 121, 84, 2,
			[]testPattern{
				{0x01, '\u0d02'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0f, '\u0d0f'},
				{0x10, '\u0d10'},
				{0x1b, 0x1b},
				{0x1f, '\u0d1e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0d2d'},
				{0x50, '\u0d3e'},
				{0x70, 'p'},
				{0x7b, '\u0d60'},
				{0x7d, '\u0d62'},
				{0x7f, '\u0d79'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1c, '\u0d66'},
				{0x1e, '\u0d68'},
				{0x1f, '\u0d69'},
				{0x2f, '\\'},
				{0x30, '\u0d7b'},
				{0x34, '\u0d7f'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Oriya, 117, 77, 2,
			[]testPattern{
				{0x00, '\u0b01'},
				{0x0a, '\n'},
				{0x0b, '\u0b0c'},
				{0x0d, '\r'},
				{0x0f, '\u0b0f'},
				{0x10, '\u0b10'},
				{0x1b, 0x1b},
				{0x1c, '\u0b1b'},
				{0x1f, '\u0b1e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0b2d'},
				{0x50, '\u0b3e'},
				{0x70, 'p'},
				{0x7b, '\u0b57'},
				{0x7c, '\u0b60'},
				{0x7f, '\u0b63'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1c, '\u0b66'},
				{0x1e, '\u0b68'},
				{0x1f, '\u0b69'},
				{0x21, '\u0b6b'},
				{0x2a, '\u0b5f'},
				{0x2c, '\u0b71'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Punjabi, 111, 78, 2,
			[]testPattern{
				{0x00, '\u0a01'},
				{0x0a, '\n'},
				{0x08, '\u0a0a'},
				{0x0d, '\r'},
				{0x0f, '\u0a0f'},
				{0x10, '\u0a10'},
				{0x1b, 0x1b},
				{0x1c, '\u0a1b'},
				{0x1f, '\u0a1e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0a2d'},
				{0x50, '\u0a3e'},
				{0x70, 'p'},
				{0x7b, '\u0a70'},
				{0x7c, '\u0a71'},
				{0x7f, '\u0a74'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1c, '\u0a66'},
				{0x1e, '\u0a68'},
				{0x1f, '\u0a69'},
				{0x21, '\u0a6b'},
				{0x2a, '\u0a5b'},
				{0x2d, '\u0a75'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Tamil, 103, 79, 2,
			[]testPattern{
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
			},
			[]testPattern{
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
			},
		},
		{charset.Telugu, 121, 80, 2,
			[]testPattern{
				{0x00, '\u0c01'},
				{0x0a, '\n'},
				{0x0b, '\u0c0c'},
				{0x0d, '\r'},
				{0x0f, '\u0c0f'},
				{0x10, '\u0c10'},
				{0x1b, 0x1b},
				{0x1f, '\u0c1e'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, '\u0c2d'},
				{0x50, '\u0c3e'},
				{0x70, 'p'},
				{0x7b, '\u0c56'},
				{0x7d, '\u0c61'},
				{0x7f, '\u0c63'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1c, '\u0c66'},
				{0x1e, '\u0c68'},
				{0x1f, '\u0c69'},
				{0x2f, '\\'},
				{0x30, '\u0c7d'},
				{0x32, '\u0c7f'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
		{charset.Urdu, 128, 92, 2,
			[]testPattern{
				{0x00, 'ا'},
				{0x0a, '\n'},
				{0x0d, '\r'},
				{0x0e, 'ٺ'},
				{0x0f, 'ټ'},
				{0x10, 'ث'},
				{0x1b, 0x1b},
				{0x1c, 'ڌ'},
				{0x1f, 'ڊ'},
				{0x20, ' '},
				{0x30, '0'},
				{0x41, 'ض'},
				{0x50, 'ں'},
				{0x70, 'p'},
				{0x7b, '\u0655'},
				{0x7c, '\u0651'},
				{0x7f, '\u0670'},
			},
			[]testPattern{
				{0x00, '@'},
				{0x08, '&'},
				{0x0a, '\f'},
				{0x0b, '*'},
				{0x13, '¡'},
				{0x1d, '۱'},
				{0x1e, '۲'},
				{0x1f, '۳'},
				{0x2f, '\\'},
				{0x40, '|'},
				{0x47, 'G'},
				{0x65, '€'},
			},
		},
	}
)
