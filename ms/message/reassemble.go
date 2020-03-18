// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package message

import (
	"encoding/binary"

	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
	"github.com/warthog618/sms/ms/sar"
)

// Message represents a message received from an origination number.
// The message is provided in UTF8.
// The message was contained in the associated TPDUs.
type Message struct {
	Msg    string
	Number string
	TPDUs  []*tpdu.Deliver
}

// Reassembler is responsible for collecting TPDUs and building Messages from
// them using the DataDecoder (typically a tpdu.UDDecoder).
type Reassembler struct {
	c Collector
	d UDDecoder
}

// UDDecoder provides a Decode method to convert the user data from a TPDU
// into the corresponding UTF-8 message.
type UDDecoder interface {
	Decode(ud tpdu.UserData, udh tpdu.UserDataHeader, alpha tpdu.Alphabet) ([]byte, error)
}

// Collector collects the segments of a concatenated SMS and returns the
// completed set when available.
type Collector interface {
	Collect(pdu *tpdu.Deliver) ([]*tpdu.Deliver, error)
	Close()
}

// NewReassembler creates a Reassembler.
func NewReassembler(options ...ReassemblerOption) *Reassembler {
	rc := ReassemblerConfig{}
	for _, option := range options {
		option.applyReassemblerOption(&rc)
	}
	if rc.c == nil {
		rc.c = sar.NewCollector()
	}
	if rc.d == nil {
		rc.d = tpdu.NewUDDecoder(rc.dopts...)
	}
	r := Reassembler{rc.c, rc.d}
	return &r
}

type ReassemblerConfig struct {
	d     UDDecoder
	dopts []tpdu.UDDecoderOption
	c     Collector
	//copts []sar.CollectorOption
}

type ReassemblerOption interface {
	applyReassemblerOption(*ReassemblerConfig)
}

func WithCollector(c Collector) CollectorOption {
	return CollectorOption{c}
}

type CollectorOption struct {
	c Collector
}

func (o CollectorOption) applyReassemblerOption(r *ReassemblerConfig) {
	r.c = o.c
}

func WithDataDecoder(d UDDecoder) DataDecoderOption {
	return DataDecoderOption{d}
}

type DataDecoderOption struct {
	d UDDecoder
}

func (o DataDecoderOption) applyReassemblerOption(r *ReassemblerConfig) {
	r.d = o.d
}

// Close terminates the reassembler and all the reassembly pipes currently active.
func (r *Reassembler) Close() {
	r.c.Close()
}

// Reassemble takes a binary Deliver TPDU and adds it to the reassembly
// collection. If the Deliver is the last TPDU in a set then the completed
// Message is returned.
func (r *Reassembler) Reassemble(b []byte) (*Message, error) {
	d := tpdu.NewDeliver()
	err := d.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	segments, err := r.c.Collect(d)
	if err != nil {
		return nil, err
	}
	if segments != nil {
		return r.concatenate(segments)
	}
	return nil, nil
}

// Concatenate converts a set of concatenated TPDUs into a Message.
// The User Data in each TPDU is converted to UTF-8 using the DataDecoder.
func (r *Reassembler) concatenate(segments []*tpdu.Deliver) (*Message, error) {
	bl := 0
	ts := make([][]byte, len(segments))
	var danglingSurrogate tpdu.UserData
	for i, s := range segments {
		a, _ := s.Alphabet()
		ud := s.UD
		if danglingSurrogate != nil {
			ud = append(danglingSurrogate, ud...)
			danglingSurrogate = nil
		}
		d, err := r.d.Decode(ud, s.UDH, a)
		if err != nil {
			switch e := err.(type) {
			case ucs2.ErrDanglingSurrogate:
				danglingSurrogate = []byte{0, 0}
				binary.BigEndian.PutUint16(danglingSurrogate, e.Surrogate())
			default:
				return nil, err
			}
		}
		ts[i] = d
		bl += len(d)
	}
	m := make([]byte, 0, bl)
	for _, t := range ts {
		m = append(m, t...)
	}
	return &Message{string(m), segments[0].OA.Number(), segments}, nil
}
