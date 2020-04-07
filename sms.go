// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

// Package sms provides encoders and decoders for SMS PDUs.
package sms

import (
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
)

// ConcatConfig contains configuration option for Concatenate.
type ConcatConfig struct {
	dopts []tpdu.UDDecodeOption
}

// Concatenate converts a set of concatenated TPDUs into a UTF-8 message.
//
// Assumes segments are the component TPDUs of a segmented message, in correct order.
// This is the case for segments returned by the Collector.
// It can be tested using IsSegmentedMessage.
func Concatenate(segments []*tpdu.TPDU, options ...ConcatOption) ([]byte, error) {
	cfg := ConcatConfig{}
	for _, option := range options {
		option.ApplyConcatOption(&cfg)
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

// IsSegmentedMessage confirms that the segment TPDUs are suitable for being
// concatenated into a message.
func IsSegmentedMessage(segments []*tpdu.TPDU) bool {
	if len(segments) == 0 {
		return false
	}
	baseSegs, _, baseConcatRef, ok := segments[0].ConcatInfo()
	if !ok {
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

// DecodeConfig contains configuration options for Decode.
type DecodeConfig struct {
	dirn tpdu.Direction
}

// Decode converts a binary SMS TPDU into the corresponding TPDU object.
func Decode(src []byte, options ...DecodeOption) (*tpdu.TPDU, error) {
	cfg := DecodeConfig{}
	for _, option := range options {
		option.ApplyDecodeOption(&cfg)
	}
	t := tpdu.TPDU{Direction: cfg.dirn}
	err := t.UnmarshalBinary(src)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
