// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package message

import (
	"sync"

	"github.com/warthog618/sms/encoding/tpdu"
)

// Encoder builds Submit TPDUs from simple inputs such as the destination
// number and the message in a UTF8 form.
type Encoder struct {
	e        DataEncoder
	s        Segmenter
	mutex    sync.Mutex // covers msgCount and t
	msgCount int
	t        *tpdu.Submit
}

// DataEncoder converts a UTF-8 message into the corresponding TPDU user data.
type DataEncoder interface {
	Encode(msg string) (tpdu.UserData, tpdu.UserDataHeader, tpdu.Alphabet, error)
}

// Segmenter segments a large outgoing message into the set of Submit TPDUs
// required to contain it.
type Segmenter interface {
	Segment(msg []byte, t *tpdu.Submit) []tpdu.Submit
}

// NewEncoder creates an Encoder.
func NewEncoder(e DataEncoder, s Segmenter) *Encoder {
	return &Encoder{e, s, sync.Mutex{}, 0, nil}
}

// SetT sets the template Submit TPDU used by Encode.
// The Submit TPDU is used to populate the fields for encoded Submit TPDUs,
// with the exception of the MR, DA and UD which are explicitly set by Encode.
// Encode also sets the DCS alphabet, and may add elements to the UDH.
// The provided DCS may contain a message class, but will be completely ignored
// if the value is incompatible with setting the alphabet.
func (e *Encoder) SetT(t *tpdu.Submit) {
	e.mutex.Lock()
	e.t = t
	e.mutex.Unlock()
}

// Encode builds a set of Submit TPDUs from the destination number and UTF8 message.
// Long messages are split into multiple concatenated TPDUs, while short messages
// may fit in one.
func (e *Encoder) Encode(number, msg string) ([]tpdu.Submit, error) {
	d, udh, alpha, err := e.e.Encode(msg)
	if err != nil {
		return nil, err
	}
	s := tpdu.NewSubmit()
	e.mutex.Lock()
	if e.t != nil {
		*s = *e.t
		s.SetUDH(append(e.t.UDH(), udh...))
	} else {
		s.SetUDH(udh)
	}
	if len(number) > 0 && number[0] == '+' {
		number = number[1:]
	}
	s.SetDA(tpdu.Address{TOA: 0x80 | byte(tpdu.TonInternational<<4) | byte(tpdu.NpISDN), Addr: number})
	dcs, err := s.DCS().WithAlphabet(alpha)
	if err != nil {
		// ignore the template dcs
		dcs, _ = tpdu.DCS(0).WithAlphabet(alpha)
	}
	s.SetDCS(dcs)
	segments := e.segment(d, s)
	e.mutex.Unlock()
	return segments, nil
}

// Encode8Bit builds a set of Submit TPDUs from the destination number and raw binary message.
// Long messages are split into multiple concatenated TPDUs, while short messages
// may fit in one.
func (e *Encoder) Encode8Bit(number string, d []byte) ([]tpdu.Submit, error) {
	s := tpdu.NewSubmit()
	e.mutex.Lock()
	if e.t != nil {
		*s = *e.t
	}
	if len(number) > 0 && number[0] == '+' {
		number = number[1:]
	}
	s.SetDA(tpdu.Address{TOA: 0x80 | byte(tpdu.TonInternational<<4) | byte(tpdu.NpISDN), Addr: number})
	dcs, err := s.DCS().WithAlphabet(tpdu.Alpha8Bit)
	if err != nil {
		// ignore the template dcs
		dcs, _ = tpdu.DCS(0).WithAlphabet(tpdu.Alpha8Bit)
	}
	s.SetDCS(dcs)
	segments := e.segment(d, s)
	e.mutex.Unlock()
	return segments, nil
}

func (e *Encoder) segment(d []byte, s *tpdu.Submit) []tpdu.Submit {
	segments := e.s.Segment(d, s)
	for _, sg := range segments {
		e.msgCount++
		sg.SetMR(byte(e.msgCount))
	}
	return segments
}
