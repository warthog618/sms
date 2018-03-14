// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ucs2_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/ucs2"
)

type decodePattern struct {
	name string
	in   []byte
	out  []rune
	err  error
}

func TestDecode(t *testing.T) {
	patterns := []decodePattern{
		{"nil", nil, nil, nil},
		{"empty", []byte(""), nil, nil},
		{"odd", []byte{1, 2, 3, 4, 5}, nil, ucs2.ErrInvalidLength},
		{"howdy", []byte{0x4F, 0x60, 0x59, 0x7D, 0xFF, 0x01, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x64, 0x00, 0x79},
			[]rune("你好！Howdy"), nil},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			dst, err := ucs2.Decode(p.in)
			if err != p.err {
				t.Errorf("error decoding %v: %v", p.in, err)
			}
			if string(dst) != string(p.out) {
				t.Errorf("failed to decode %v: expected '%v', got %v", p.in, p.out, dst)
			}
		}
		t.Run(p.name, f)
	}
}

type encodePattern struct {
	name string
	in   []rune
	out  []byte
	err  error
}

func TestEncode(t *testing.T) {
	patterns := []encodePattern{
		{"nil", nil, nil, nil},
		{"empty", []rune(""), nil, nil},
		{"howdy", []rune("你好！Howdy"),
			[]byte{0x4F, 0x60, 0x59, 0x7D, 0xFF, 0x01, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x77, 0x00, 0x64, 0x00, 0x79}, nil},
		{"invalid", []rune("\U0010FFFF"), nil, ucs2.ErrInvalidRune('\U0010FFFF')},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			dst, err := ucs2.Encode([]rune(p.in))
			if err != p.err {
				t.Errorf("error encoding %v: %v", p.in, err)
			}
			if !bytes.Equal(p.out, dst) {
				t.Errorf("failed to encode %v: expected %v, got %v", p.in, p.out, dst)
			}
		}
		t.Run(p.name, f)
	}
}

func TestErrInvalidRune(t *testing.T) {
	patterns := []byte{0x00, 0xa0, 0x0a, 0x9a, 0xa9, 0xff}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := ucs2.ErrInvalidRune(p)
			expected := fmt.Sprintf("ucs2: invalid rune: %U", p)
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify %x, expected '%s', got '%s'", p, expected, s)
			}
		}
		t.Run(fmt.Sprintf("%x", p), f)
	}
}
