// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

// Package sms provides encoders and decoders for SMS PDUs.
package sms

import (
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
)

// DecodeConfig contains configuration option for Decode.
type DecodeConfig struct {
	dopts []tpdu.UDDecodeOption
}

// Decode returns the UTF-8 message contained in a set of TPDUs.
//
// For concatenated messages the segments assumed to be the component TPDUs, in
// correct order. This is the case for segments returned by the Collector. It
// can be tested using IsCompleteMessage.
func Decode(segments []*tpdu.TPDU, options ...DecodeOption) ([]byte, error) {
	cfg := DecodeConfig{}
	for _, option := range options {
		option.ApplyDecodeOption(&cfg)
	}
	if len(cfg.dopts) == 0 {
		cfg.dopts = []tpdu.UDDecodeOption{tpdu.WithAllCharsets}
	}
	bl := 0
	ts := make([][]byte, len(segments))
	var danglingSurrogate ucs2.ErrDanglingSurrogate
	for i, s := range segments {
		a, _ := s.Alphabet()
		ud := s.UD
		if danglingSurrogate != nil {
			ud = append([]byte(danglingSurrogate), ud...)
			danglingSurrogate = nil
		}
		d, err := tpdu.DecodeUserData(ud, s.UDH, a, cfg.dopts...)
		if err != nil {
			switch e := err.(type) {
			case ucs2.ErrDanglingSurrogate:
				danglingSurrogate = e
			default:
				return nil, err
			}
		}
		ts[i] = d
		bl += len(d)
	}
	if danglingSurrogate != nil {
		return nil, danglingSurrogate
	}
	m := make([]byte, 0, bl)
	for _, t := range ts {
		m = append(m, t...)
	}
	return m, nil
}

// IsCompleteMessage confirms that the TPDUs contain all the sgements required
// to reassemble a complete message and are in the correct order.
func IsCompleteMessage(segments []*tpdu.TPDU) bool {
	if len(segments) == 0 {
		return false
	}
	baseSegs, _, baseConcatRef, ok := segments[0].ConcatInfo()
	if !ok {
		if len(segments) == 1 {
			return true
		}
		return false
	}
	if baseSegs != len(segments) {
		return false
	}
	for i, s := range segments {
		segs, seqno, concatRef, ok := s.ConcatInfo()
		if !ok {
			return false
		}
		if segs != baseSegs {
			return false
		}
		if concatRef != baseConcatRef {
			return false
		}
		if seqno != i+1 {
			return false
		}
	}
	return true
}

// UnmarshalConfig contains configuration options for Unmarshal.
type UnmarshalConfig struct {
	dirn tpdu.Direction
}

// Unmarshal converts a binary SMS TPDU into the corresponding TPDU object.
func Unmarshal(src []byte, options ...UnmarshalOption) (*tpdu.TPDU, error) {
	cfg := UnmarshalConfig{}
	for _, option := range options {
		option.ApplyUnmarshalOption(&cfg)
	}
	t := tpdu.TPDU{Direction: cfg.dirn}
	err := t.UnmarshalBinary(src)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
