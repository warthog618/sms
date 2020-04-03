// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"errors"
	"fmt"
	"io"
)

type decodeError struct {
	Field  string
	Offset int
	Err    error
}

// DecodeError creates a decodeError which identifies the field being decoded,
// and the offset into the byte array where the field starts.
// If the provided error is a nested decodeError then the offset is updated to
// provide the offset from the beginning of the enclosing field, and the field
// names are combined in outer.inner format.
func DecodeError(f string, o int, e error) error {
	if s, ok := e.(decodeError); ok {
		s.Field = fmt.Sprintf("%s.%s", f, s.Field)
		s.Offset = s.Offset + o
		return s
	}
	if e == io.EOF {
		e = ErrUnderflow
	}
	return decodeError{f, o, e}
}

type encodeError struct {
	Field string
	Err   error
}

// EncodeError creates an encodeError which identifies the field being encoded.
// If the provided error is a nested encodeError then the error is returned as
// is rather than wrapping it.
func EncodeError(f string, e error) error {
	if s, ok := e.(encodeError); ok {
		s.Field = fmt.Sprintf("%s.%s", f, s.Field)
		return s
	}
	return encodeError{f, e}
}

func (e encodeError) Error() string {
	return fmt.Sprintf("tpdu: error encoding %s: %v", e.Field, e.Err)
}

func (e decodeError) Error() string {
	return fmt.Sprintf("tpdu: error decoding %s at octet %d: %v", e.Field, e.Offset, e.Err)
}

// ErrUnsupportedSmsType indicates the type of TPDU being decoded is not
// unsupported by the decoder.
type ErrUnsupportedSmsType byte

func (e ErrUnsupportedSmsType) Error() string {
	return fmt.Sprintf("unsupported SMS type: 0x%x", uint(e))
}

var (
	// ErrInvalid indicates the value of a field provided to an encoder is not valid.
	ErrInvalid = errors.New("invalid")
	// ErrOddUCS2Length indicates the length of a binary array containing UCS2
	// characters has an uneven length, and so has split a UCS2 character.
	ErrOddUCS2Length = errors.New("odd UCS2 length")
	// ErrOverlength indicates the binary provided contains more bytes than
	// expected by the TPDU decoder.
	ErrOverlength = errors.New("overlength")
	// ErrMissing indicates a field requiored to marshal an object is missing.
	ErrMissing = errors.New("missing")
	// ErrNonZero indicates a field which is expected to be zeroed, but contains
	// non-zero data.
	ErrNonZero = errors.New("non-zero fill")
	// ErrUnderflow indicates the binary provided does not contain
	// sufficient bytes to correctly decode the TPDU.
	ErrUnderflow = errors.New("underflow")
)
