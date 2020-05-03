// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package sms

import "github.com/warthog618/sms/encoding/tpdu"

// EncoderOption is an optional mutator for the Encoder.
type EncoderOption interface {
	ApplyEncoderOption(*Encoder)
}

// DecodeOption defines options for Decode.
type DecodeOption interface {
	ApplyDecodeOption(*DecodeConfig)
}

// UnmarshalOption defines options for Unmarhsal.
type UnmarshalOption interface {
	ApplyUnmarshalOption(*UnmarshalConfig)
}

// WithTemplate specifies the TPDU to be used as the template for encoding.
func WithTemplate(t tpdu.TPDU) EncoderOption {
	return tpduTemplate{t}
}

type tpduTemplate struct {
	t tpdu.TPDU
}

func (o tpduTemplate) ApplyEncoderOption(e *Encoder) {
	e.pdu = o.t
}

type templateOption struct {
	tpdu.Option
}

func (o templateOption) ApplyEncoderOption(e *Encoder) {
	o.ApplyTPDUOption(&e.pdu)
}

// WithTemplateOption wraps a TPDU option in a TemplateOption so it can be
// applied to an Encoder template PDU.
func WithTemplateOption(option tpdu.Option) EncoderOption {
	return templateOption{option}
}

var (
	// AsSubmit indicates that generated PDUs will be of type SmsSubmit.
	AsSubmit = templateOption{tpdu.SmsSubmit}

	// AsDeliver indicates that generated PDUs will be of type SmsDeliver.
	AsDeliver = templateOption{tpdu.SmsDeliver}

	// As8Bit indicates that generated PDUs encode user data as 8bit.
	As8Bit = templateOption{tpdu.Dcs8BitData}

	// AsUCS2 indicates that generated PDUs encode user data as UCS2.
	AsUCS2 = templateOption{tpdu.DcsUCS2Data}

	// AsMO indicates that the TPDU originated from the mobile station.
	AsMO = directionOption{tpdu.MO}

	// AsMT indicates that the TPDU as destined for the mobile station.
	AsMT = directionOption{tpdu.MT}

	// WithAllCharsets specifies that all character sets are available for
	// encoding or decoding.
	//
	// This is the default policy for decoding.
	WithAllCharsets = AllCharsetsOption{}

	// WithDefaultCharset specifies that only the default character set is
	// available for encoding or decoding.
	//
	// This is the default policy for encoding.
	WithDefaultCharset = CharsetOption{}
)

// To specifies the DA for a SMS-SUBMIT TPDU.
func To(number string) EncoderOption {
	addr := tpdu.NewAddress(tpdu.FromNumber(number))
	return templateOption{tpdu.WithDA(addr)}
}

// From specifies the OA for a SMS-DELIVER TPDU.
func From(number string) EncoderOption {
	addr := tpdu.NewAddress(tpdu.FromNumber(number))
	return templateOption{tpdu.WithOA(addr)}
}

// AllCharsetsOption specifies that all charactersets are available for encoding.
type AllCharsetsOption struct{}

// ApplyEncoderOption applies the AllCharsetsOption to an Encoder.
func (o AllCharsetsOption) ApplyEncoderOption(e *Encoder) {
	e.eopts = append(e.eopts, tpdu.WithAllCharsets)
}

// WithCharset creates an CharsetOption.
func WithCharset(nli ...int) CharsetOption {
	return CharsetOption{nli}
}

// CharsetOption defines the character sets available for encoding or decoding.
type CharsetOption struct {
	nli []int
}

// ApplyEncoderOption applies the CharsetOption to an Encoder.
func (o CharsetOption) ApplyEncoderOption(e *Encoder) {
	e.eopts = append(e.eopts, tpdu.WithCharset(o.nli...))
}

// ApplyDecodeOption applies the CharsetOption to decoding.
func (o CharsetOption) ApplyDecodeOption(cc *DecodeConfig) {
	cc.dopts = append(cc.dopts, tpdu.WithCharset(o.nli...))
}

// WithLockingCharset creates an LockingCharsetOption.
func WithLockingCharset(nli ...int) LockingCharsetOption {
	return LockingCharsetOption{nli}
}

// LockingCharsetOption defines the locking character sets available for
// encoding or decoding.
type LockingCharsetOption struct {
	nli []int
}

// ApplyEncoderOption applies the LockingCharsetOption to an Encoder.
func (o LockingCharsetOption) ApplyEncoderOption(e *Encoder) {
	e.eopts = append(e.eopts, tpdu.WithLockingCharset(o.nli...))
}

// ApplyDecodeOption applies the LockingCharsetOption to decoding.
func (o LockingCharsetOption) ApplyDecodeOption(cc *DecodeConfig) {
	cc.dopts = append(cc.dopts, tpdu.WithLockingCharset(o.nli...))
}

// WithShiftCharset creates an ShiftCharsetOption.
func WithShiftCharset(nli ...int) ShiftCharsetOption {
	return ShiftCharsetOption{nli}
}

// ShiftCharsetOption defines the shift character sets available for encoding
// or decoding.
type ShiftCharsetOption struct {
	nli []int
}

// ApplyEncoderOption applies the ShiftCharsetOption to an Encoder.
func (o ShiftCharsetOption) ApplyEncoderOption(e *Encoder) {
	e.eopts = append(e.eopts, tpdu.WithShiftCharset(o.nli...))
}

// ApplyDecodeOption applies the ShiftCharsetOption to decoding.
func (o ShiftCharsetOption) ApplyDecodeOption(cc *DecodeConfig) {
	cc.dopts = append(cc.dopts, tpdu.WithShiftCharset(o.nli...))
}

type directionOption struct {
	d tpdu.Direction
}

func (o directionOption) ApplyUnmarshalOption(d *UnmarshalConfig) {
	d.dirn = o.d
}
