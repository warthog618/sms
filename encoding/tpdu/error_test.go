// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/warthog618/sms/encoding/tpdu"
)

type decodeTestPattern struct {
	Field  string
	Offset int
	Err    error
}

type DecodeError interface {
	Field() string
	Offset() int
}

type EncodeError interface {
	Field() string
}

// TestDecodeError tests that the errors can be stringified.
// It is fragile, as it compares the strings exactly, but its main purpose is
// to confirm the Error function doesn't recurse, as that is bad.
func TestDecodeError(t *testing.T) {
	patterns := []decodeTestPattern{
		{"nil", 0, nil},
		{"err", 2, errors.New("an error")},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := tpdu.NewDecodeError(p.Field, p.Offset, p.Err)
			expected := fmt.Sprintf("tpdu: error decoding %s at octet %d: %v", p.Field, p.Offset, p.Err)
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify, expected '%s', got '%s'", expected, s)
			}
		}
		t.Run(p.Field, f)
	}
	// nested
	f := func(t *testing.T) {
		err := tpdu.NewDecodeError("nested", 40, tpdu.NewDecodeError("inner", 2, nil))
		expected := fmt.Sprintf("tpdu: error decoding nested.inner at octet 42: %v", nil)
		s := err.Error()
		if s != expected {
			t.Errorf("failed to stringify, expected '%s', got '%s'", expected, s)
		}
	}
	t.Run("nested", f)
}

// TestEncodeError tests that the errors can be stringified.
// It is fragile, as it compares the strings exactly, but its main purpose is
// to confirm the Error function doesn't recurse, as that is bad.
func TestEncodeError(t *testing.T) {
	patterns := []decodeTestPattern{
		{"nil", 0, nil},
		{"err", 2, errors.New("an error")},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := tpdu.EncodeError(p.Field, p.Err)
			expected := fmt.Sprintf("tpdu: error encoding %s: %v", p.Field, p.Err)
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify, expected '%s', got '%s'", expected, s)
			}
		}
		t.Run(p.Field, f)
	}
	// nested
	f := func(t *testing.T) {
		err := tpdu.EncodeError("nested", tpdu.EncodeError("inner", nil))
		expected := fmt.Sprintf("tpdu: error encoding nested.inner: %v", nil)
		s := err.Error()
		if s != expected {
			t.Errorf("failed to stringify, expected '%s', got '%s'", expected, s)
		}
	}
	t.Run("nested", f)
}

// TestErrUnsupportedMTI tests that the errors can be stringified.
// It is fragile, as it compares the strings exactly, but its main purpose is
// to confirm the Error function doesn't recurse, as that is bad.
func TestErrUnsupportedSmsType(t *testing.T) {
	patterns := []byte{0x00, 0xa0, 0x0a, 0x9a, 0xa9, 0xff}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := tpdu.ErrUnsupportedSmsType(p)
			expected := fmt.Sprintf("unsupported SMS type: 0x%x", uint(err))
			s := err.Error()
			if s != expected {
				t.Errorf("failed to stringify %02x, expected '%s', got '%s'", p, expected, s)
			}
		}
		t.Run(fmt.Sprintf("%x", p), f)
	}
}
