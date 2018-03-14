// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gsm7_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/gsm7"
)

type decoderPattern struct {
	name string
	in   []byte
	out  []byte
	err  error
}

type encoderPattern struct {
	name string
	in   []byte
	out  []byte
	err  error
}

func testDecoder(t *testing.T, d gsm7.Decoder, patterns []decoderPattern) {
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := d.Decode(p.in)
			if err != p.err {
				t.Errorf("error decoding %v: %v", p.in, err)
			}
			if !bytes.Equal(out, p.out) {
				t.Errorf("failed to decode: %v, expected %v, got %v", p.in, p.out, out)
			}
		}
		t.Run(p.name, f)
	}
}

func testEncoder(t *testing.T, e gsm7.Encoder, patterns []encoderPattern) {
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := e.Encode(p.in)
			if err != p.err {
				t.Errorf("error decoding %v: %v", p.in, err)
			}
			if !bytes.Equal(out, p.out) {
				t.Errorf("failed to decode: %v expected %v, got %v", p.in, p.out, out)
			}
		}
		t.Run(p.name, f)
	}
}

func TestDecode(t *testing.T) {
	d := gsm7.NewDecoder()
	p := []decoderPattern{
		{"base", []byte("message"), []byte("message"), nil},
		{"ext", []byte("\x1b\x28\x1b\x29"), []byte("{}"), nil},
		{"escaped", []byte("mes\x1b\x40sage"), []byte("mes|sage"), nil},
		{"double escaped", []byte("mes\x1b\x1b\x40sage"), []byte("mes ¡sage"), nil},
		{"dangling escape", []byte("message\x1b"), []byte("message "), nil},
	}
	testDecoder(t, d, p)
}

func TestDecoderWithCharset(t *testing.T) {
	set := map[byte]rune{'m': 'M', 'e': 'E', 's': 'S', 'a': 'A', 'g': 'G'}
	d := gsm7.NewDecoder().WithCharset(set)
	p := []decoderPattern{
		{"base", []byte("message"), []byte("MESSAGE"), nil},
		{"ext", []byte("\x1b\x28\x1b\x29"), []byte("{}"), nil},
		{"escaped", []byte("mes\x1b\x40sage"), []byte("MES|SAGE"), nil},
		{"double escaped", []byte("mes\x1b\x1b\x40sage"), []byte("MES  SAGE"), nil},
		{"dangling escape", []byte("message\x1b"), []byte("MESSAGE "), nil},
		{"unknown", []byte("mesMsage"), []byte("MES SAGE"), nil},
	}
	testDecoder(t, d, p)
}

func TestDecoderWithExtCharset(t *testing.T) {
	ext := map[byte]rune{0x40: 'Q'}
	d := gsm7.NewDecoder().WithExtCharset(ext)
	p := []decoderPattern{
		{"base", []byte("\x40"), []byte("¡"), nil},
		{"ext", []byte("\x1b\x40"), []byte("Q"), nil},
	}
	testDecoder(t, d, p)
}

func TestDecoderStrict(t *testing.T) {
	set := map[byte]rune{'m': 'M', 'e': 'E', 's': 'S', 'a': 'A', 'g': 'G'}
	ext := map[byte]rune{'e': 'E', 'x': 'X', 't': 'T'}
	d := gsm7.NewDecoder().Strict().WithCharset(set).WithExtCharset(ext)
	p := []decoderPattern{
		{"known", []byte("message"), []byte("MESSAGE"), nil},
		{"ext", []byte("\x1be\x1bx\x1bt"), []byte("EXT"), nil},
		{"unknown", []byte("mesMsage"), nil, gsm7.ErrInvalidSeptet('M')},
		{"unknown ext", []byte("mes\x1bmsage"), nil, gsm7.ErrInvalidSeptet('m')},
	}
	testDecoder(t, d, p)
}

func TestEncode(t *testing.T) {
	e := gsm7.NewEncoder()
	p := []encoderPattern{
		{"base", []byte("message"), []byte("message"), nil},
		{"ext", []byte("{}"), []byte("\x1b\x28\x1b\x29"), nil},
		{"escaped", []byte("mes|sage"), []byte("mes\x1b\x40sage"), nil},
		{"invalid", []byte("mesŞsage"), nil, gsm7.ErrInvalidUTF8('Ş')},
	}
	testEncoder(t, e, p)
}

func TestEncoderWithCharset(t *testing.T) {
	set := map[rune]byte{'Ş': 0x40}
	e := gsm7.NewEncoder().WithCharset(set)
	p := []encoderPattern{
		{"base", []byte("Ş"), []byte("\x40"), nil},
		{"ext", []byte("|"), []byte("\x1b\x40"), nil},
	}
	testEncoder(t, e, p)
}

func TestEncoderWithExtCharset(t *testing.T) {
	ext := map[rune]byte{'Ş': 0x40}
	e := gsm7.NewEncoder().WithExtCharset(ext)
	p := []encoderPattern{
		{"base", []byte("¡"), []byte("\x40"), nil},
		{"ext", []byte("Ş"), []byte("\x1b\x40"), nil},
	}
	testEncoder(t, e, p)
}

// TestErrInvalidUTF8 tests that the errors can be stringified.
// It is fragile, as it compares the strings exactly, but its main purpose is
// to confirm the Error function doesn't recurse, as that is bad.
func TestErrInvalidSeptet(t *testing.T) {
	patterns := []byte{0x00, 0xa0, 0x0a, 0x9a, 0xa9, 0xff}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := gsm7.ErrInvalidSeptet(p)
			expected := fmt.Sprintf("gsm7: invalid septet 0x%02x", int(err))
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify %02x, expected '%s', got '%s'", p, expected, s)
			}
		}
		t.Run(fmt.Sprintf("%x", p), f)
	}
}

// TestErrInvalidUTF8 tests that the errors can be stringified.
// It is fragile, as it compares the strings exactly, but its main purpose is
// to confirm the Error function doesn't recurse, as that is bad.
func TestErrInvalidUTF8(t *testing.T) {
	patterns := []byte{0x00, 0xa0, 0x0a, 0x9a, 0xa9, 0xff}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := gsm7.ErrInvalidUTF8(p)
			expected := fmt.Sprintf("gsm7: invalid utf8 '%c' (0x%04x)", rune(err), int(err))
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify %02x, expected '%s', got '%s'", p, expected, s)
			}
		}
		t.Run(fmt.Sprintf("%x", p), f)
	}
}
