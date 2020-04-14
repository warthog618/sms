// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package sms

import (
	"sync/atomic"

	"github.com/warthog618/sms/encoding/tpdu"
)

// Encode builds a set of TPDUs containing the message.
//
// Long messages are split into multiple concatenated TPDUs, while short
// messages may fit in one.
//
// By default messages are encoded into SMS-SUBMIT TPDUs.  This behaviour may
// be overridden via options.
//
// For 8-bit encoding the message is encoded as is.
//
// For 7-bit encoding the message is assumed to contain UTF-8.
//
// For explicit UCS-2 encoding the message is assumed to contain UTF-16,
// encoded as an array of bytes.  This can be created from an array of UTF-16
// runes using ucs2.Encode.
//
// For implicit UCS-2 encoding (the fallback with 7-bit fails) the message is
// assumed to contain UTF-8.
func Encode(msg []byte, options ...EncoderOption) ([]tpdu.TPDU, error) {
	options = append([]EncoderOption{AsSubmit}, options...)
	e := NewEncoder(options...)
	return e.Encode(msg)
}

// Encoder builds SMS TPDUs from simple inputs such as the destination number
// and the message in a UTF8 form.
type Encoder struct {
	// options for encoding UD
	eopts []tpdu.UDEncodeOption

	// options for segmentation
	sopts []tpdu.SegmentationOption

	// The template TPDU for encoding.
	pdu tpdu.TPDU

	// MsgCount is the number of TPDUs encoded.
	MsgCount tpdu.Counter

	// ConcatRef is tghe number of multi-segment messages encoded.
	ConcatRef tpdu.Counter
}

// NewEncoder creates an Encoder.
func NewEncoder(options ...EncoderOption) *Encoder {
	e := Encoder{}
	for _, option := range options {
		option.ApplyEncoderOption(&e)
	}
	if e.MsgCount == nil {
		e.MsgCount = &Counter{}
	}
	if e.ConcatRef == nil {
		e.ConcatRef = &Counter{}
	}
	return &e
}

// Encode builds a set of TPDUs containing the message.
//
// Long messages are split into multiple concatenated TPDUs, while short
// messages may fit in one.
//
// By default messages are encoded into SMS-DELIVER TPDUs.  This behaviour may
// be overridden via options, either to NewEncoder or Encode.
//
// For 8-bit encoding the message is encoded as is.
//
// For 7-bit encoding the message is assumed to contain UTF-8.
//
// For explicit UCS-2 encoding the message is assumed to contain UTF-16,
// encoded as an array of bytes.  This can be created from an array of UTF-16
// runes using ucs2.Encode.
//
// For implicit UCS-2 encoding (the fallback with 7-bit fails) the message is
// assumed to contain UTF-8.
func (e Encoder) Encode(msg []byte, options ...EncoderOption) ([]tpdu.TPDU, error) {
	for _, option := range options {
		option.ApplyEncoderOption(&e)
	}
	sopts := append(e.sopts, tpdu.WithMR(e.MsgCount), tpdu.WithConcatRef(e.ConcatRef))
	// take the DCS in the template TPDU as a hint...
	alpha, _ := e.pdu.DCS.Alphabet()
	switch alpha {
	case tpdu.Alpha8Bit, tpdu.AlphaUCS2:
		return e.pdu.Segment(msg, sopts...), nil
	default:
		// encode as GSM7, or failing that UCS2...
		d, udh, alpha := tpdu.EncodeUserData(msg, e.eopts...)
		dcs, err := e.pdu.DCS.WithAlphabet(alpha)
		if err != nil {
			return nil, ErrDcsConflict
		}
		if dcs != e.pdu.DCS {
			e.pdu.SetDCS(byte(dcs))
		}
		if udh != nil {
			e.pdu.SetUDH(append(e.pdu.UDH[:0:0], udh...))
		}
		return e.pdu.Segment(d, sopts...), nil
	}
}

// Counter is an implementation of the tpdu.Counter interface.
//
// It also provides a Read method on the current value for diagnostic purposes.
type Counter struct {
	c int64
}

// Count increments and returns the counter.
func (c *Counter) Count() int {
	return int(atomic.AddInt64(&c.c, 1))
}

// Read returns the counter.
func (c *Counter) Read() int {
	return int(atomic.LoadInt64(&c.c))
}
