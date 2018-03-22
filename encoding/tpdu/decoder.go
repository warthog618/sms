// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"encoding"

	"github.com/pkg/errors"
)

// Decoder converts binary TPDUs to the corresponding TPDU implementation.
type Decoder struct {
	// d maps MessageType and Direction to the corresponding decoder.
	d map[byte]ConcreteDecoder
}

// DecoderOption is function that modifies an existing Decoder.
type DecoderOption func(*Decoder) error

// NewDecoder creates a new Decoder.
func NewDecoder(opts ...DecoderOption) (*Decoder, error) {
	d := &Decoder{map[byte]ConcreteDecoder{}}
	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// TPDU represents the minimal interface provided by TPDU implementations.
type TPDU interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	// MTI provides a clue as to the underlying TPDU type.
	MTI() MessageType
	// Alphabet defines how the UD field is encoded.
	Alphabet() (Alphabet, error)
	// UD provides the UserData, the format of which depends on the Alphabet.
	UD() UserData
	// UDH provides the UserDataHeader, which is optional and so may be empty.
	UDH() UserDataHeader
}

// ConcreteDecoder is a function that decodes a binary TPDU into a particular TPDU struct.
type ConcreteDecoder func([]byte) (TPDU, error)

// RegisterDecoder registers a decoder for the given MessageType and Direction.
func (d *Decoder) RegisterDecoder(mt MessageType, drn Direction, f ConcreteDecoder) error {
	k := byte(mt) | (byte(drn) << 2)
	if _, ok := d.d[k]; ok {
		return errors.New("decoder already registered")
	}
	d.d[k] = f
	return nil
}

// Decode returns the TPDU decoded from the SMS TPDU in src.
// The direction of the SMS must be provided so the the decoder can correctly
// determine the type of TPDU from the MTI. (the same MTI is used for different
// TPDUs depending on whether the SMS is being sent to the MS, or is from the MS.)
//
// The reverse of this operation is MarshalBinary on the returned TPDU.
func (d *Decoder) Decode(src []byte, drn Direction) (TPDU, error) {
	if len(src) < 1 {
		return nil, DecodeError("firstOctet", 0, ErrUnderflow)
	}
	firstOctet := src[0]
	k := (firstOctet & 0x3) | byte(drn<<2)
	f, ok := d.d[k]
	if !ok {
		return nil, DecodeError("firstOctet", 0, ErrUnsupportedMTI(firstOctet&0x03))
	}
	pdu, err := f(src)
	if err != nil {
		return nil, err
	}
	return pdu, nil
}
