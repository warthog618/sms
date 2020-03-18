// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package message

import (
	"sync"

	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/sar"
)

// Encoder builds Submit TPDUs from simple inputs such as the destination
// number and the message in a UTF8 form.
type Encoder struct {
	ude      UDEncoder
	s        Segmenter
	sf       SubmitFactory
	mutex    sync.Mutex // covers msgCount
	msgCount int
}

type EncoderConfig struct {
	ude   UDEncoder
	eopts []tpdu.UDEncoderOption
	s     Segmenter
	sopts []sar.SegmenterOption
	sf    SubmitFactory
}

// UDEncoder converts a UTF-8 message into the corresponding TPDU user data.
type UDEncoder interface {
	Encode(msg string) (tpdu.UserData, tpdu.UserDataHeader, tpdu.Alphabet, error)
}

// Segmenter segments a large outgoing message into the set of Submit TPDUs
// required to contain it.
type Segmenter interface {
	Segment(msg []byte, t *tpdu.Submit) []tpdu.Submit
}

// SubmitFactory is a function that creates Submit TPDUs.
type SubmitFactory func(options ...tpdu.SubmitOption) *tpdu.Submit

// EncoderOption is a construction option for the Encoder.
type EncoderOption interface {
	applyEncoderOption(*EncoderConfig)
}

func WithCharset(nli ...int) CharsetOption {
	return CharsetOption{nli}
}

type CharsetOption struct {
	nli []int
}

func (o CharsetOption) applyEncoderOption(ec *EncoderConfig) {
	ec.ude = tpdu.NewUDEncoder(tpdu.WithCharset(o.nli...))
}

func (o CharsetOption) applyReassemblerOption(rc *ReassemblerConfig) {
	rc.dopts = append(rc.dopts, tpdu.WithCharset(o.nli...))
}

func WithLockingCharset(nli ...int) LockingCharsetOption {
	return LockingCharsetOption{nli}
}

type LockingCharsetOption struct {
	nli []int
}

func (o LockingCharsetOption) applyEncoderOption(ec *EncoderConfig) {
	ec.ude = tpdu.NewUDEncoder(tpdu.WithLockingCharset(o.nli...))
}

func (o LockingCharsetOption) applyReasemblerOption(rc *ReassemblerConfig) {
	rc.dopts = append(rc.dopts, tpdu.WithLockingCharset(o.nli...))
}

func WithShiftCharset(nli ...int) ShiftCharsetOption {
	return ShiftCharsetOption{nli}
}

type ShiftCharsetOption struct {
	nli []int
}

func (o ShiftCharsetOption) applyEncoderOption(ec *EncoderConfig) {
	ec.ude = tpdu.NewUDEncoder(tpdu.WithShiftCharset(o.nli...))
}

func (o ShiftCharsetOption) applyReasemblerOption(rc *ReassemblerConfig) {
	rc.dopts = append(rc.dopts, tpdu.WithShiftCharset(o.nli...))
}

// NewEncoder creates an Encoder.
func NewEncoder(options ...EncoderOption) *Encoder {
	ec := EncoderConfig{}
	for _, option := range options {
		option.applyEncoderOption(&ec)
	}
	if ec.ude == nil {
		ec.ude = tpdu.NewUDEncoder(ec.eopts...)
	}
	if ec.s == nil {
		ec.s = sar.NewSegmenter(ec.sopts...)
	}
	if ec.sf == nil {
		ec.sf = tpdu.NewSubmit
	}
	e := Encoder{ude: ec.ude, s: ec.s, sf: ec.sf}
	return &e
}

type UDEncoderOption struct {
	ude UDEncoder
}

func (o UDEncoderOption) applyEncoderOption(ec *EncoderConfig) {
	ec.ude = o.ude
}

// WithUDEncoder specifies the user data encoder to be used when encoding messages.
func WithUDEncoder(ude UDEncoder) UDEncoderOption {
	return UDEncoderOption{ude}
}

type SegmenterOption struct {
	s Segmenter
}

func (o SegmenterOption) applyEncoderOption(ec *EncoderConfig) {
	ec.s = o.s
}

// WithSegmenter specifies the segmenter to be used when encoding messages.
func WithSegmenter(s Segmenter) SegmenterOption {
	return SegmenterOption{s}
}

type SubmitFactoryOption struct {
	sf SubmitFactory
}

func (o SubmitFactoryOption) applyEncoderOption(ec *EncoderConfig) {
	ec.sf = o.sf
}

// WithSubmitFactory specifies the factory for the template Submit TPDU for
// encoding messages.
func WithSubmitFactory(sf SubmitFactory) SubmitFactoryOption {
	return SubmitFactoryOption{sf}
}

// FromSubmitPDU specifies a static template Submit TPDU for encoding
// messages.
func FromSubmitPDU(t *tpdu.Submit) SubmitFactoryOption {
	sf := func(options ...tpdu.SubmitOption) *tpdu.Submit {
		options = append([]tpdu.SubmitOption{tpdu.FromSubmit(t)}, options...)
		return tpdu.NewSubmit(options...)
	}
	return SubmitFactoryOption{sf}
}

// Encode builds a set of Submit TPDUs from the destination number and UTF8 message.
// Long messages are split into multiple concatenated TPDUs, while short
// messages may fit in one.
func (e *Encoder) Encode(number, msg string) ([]tpdu.Submit, error) {
	d, udh, alpha, err := e.ude.Encode(msg)
	if err != nil {
		return nil, err
	}
	s := e.sf(
		tpdu.To(number),
		tpdu.WithUserDataHeader(udh),
		tpdu.WithAlphabet(alpha),
	)
	return e.segment(d, s), nil
}

// Encode8Bit builds a set of Submit TPDUs from the destination number and raw
// binary message.
// Long messages are split into multiple concatenated TPDUs, while short
// messages may fit in one.
func (e *Encoder) Encode8Bit(number string, d []byte) ([]tpdu.Submit, error) {
	s := e.sf(
		tpdu.To(number),
		tpdu.WithAlphabet(tpdu.Alpha8Bit),
	)
	return e.segment(d, s), nil
}

func (e *Encoder) segment(d []byte, s *tpdu.Submit) []tpdu.Submit {
	segments := e.s.Segment(d, s)
	e.mutex.Lock()
	for _, sg := range segments {
		e.msgCount++
		sg.MR = byte(e.msgCount)
	}
	e.mutex.Unlock()
	return segments
}
