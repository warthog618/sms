// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package sms

import "github.com/warthog618/sms/encoding/tpdu"

// EncoderOption is an optional mutator for the Encoder.
type EncoderOption interface {
	ApplyEncoderOption(*Encoder)
}

// ConcatOption defines options for Concatenate.
type ConcatOption interface {
	ApplyConcatOption(*ConcatConfig)
}

// DecodeOption defines options for Concatenate.
type DecodeOption interface {
	ApplyDecodeOption(*DecodeConfig)
}

// WithTemplate specifies the TPDU to be used as the template for encoding.
func WithTemplate(t tpdu.TPDU) EncoderOption {
	return tpduTemplate{t}
}

type tpduTemplate struct {
	t tpdu.TPDU
}

// ApplyEncoderOption applies the template option to the Encoder template PDU.
func (o tpduTemplate) ApplyEncoderOption(e *Encoder) {
	e.pdu = o.t
}

// TemplateOption wraps a TPDU option so it can be applied to an Encoder
// template PDU.
type TemplateOption struct {
	tpdu.Option
}

// ApplyEncoderOption applies the template option to the Encoder template PDU.
func (o TemplateOption) ApplyEncoderOption(e *Encoder) {
	o.ApplyTPDUOption(&e.pdu)
}

// WithTemplateOption wraps a TPDU option in a TemplateOption so it can be
// applied to an Encoder template PDU.
func WithTemplateOption(option tpdu.Option) TemplateOption {
	return TemplateOption{option}
}

var (
	// AsSubmit indicates that generated PDUs will be of type SmsSubmit.
	AsSubmit = TemplateOption{tpdu.SmsSubmit}

	// AsDeliver indicates that generated PDUs will be of type SmsDeliver.
	AsDeliver = TemplateOption{tpdu.SmsDeliver}

	// As8Bit indicates that generated PDUs encide user data as 8bit.
	As8Bit = TemplateOption{tpdu.Dcs8BitData}

	// AsUCS2 indicates that generated PDUs encode user data as UCS2.
	AsUCS2 = TemplateOption{tpdu.DcsUCS2Data}

	// AsMO indicates that the TPDU originated from the mobile station.
	AsMO = DirectionOption{tpdu.MO}

	// AsMT indicates that the TPDU as destined for the mobile station.
	AsMT = DirectionOption{tpdu.MT}
)

// To specifies the DA for a SMS-SUBMIT TPDU.
func To(number string) TemplateOption {
	addr := tpdu.NewAddress(tpdu.FromNumber(number))
	return TemplateOption{tpdu.WithDA(addr)}
}

// From specifies the OA for a SMS-DELIVER TPDU.
func From(number string) TemplateOption {
	addr := tpdu.NewAddress(tpdu.FromNumber(number))
	return TemplateOption{tpdu.WithOA(addr)}
}

// AllCharsetsOption specifies that all charactersets are available for encoding.
type AllCharsetsOption struct{}

// ApplyEncoderOption applies the AllCharsetsOption to an Encoder.
func (o AllCharsetsOption) ApplyEncoderOption(e *Encoder) {
	e.eopts = append(e.eopts, tpdu.WithAllCharsets)
}

// WithAllCharsets specifies that all charactersets are available for encoding.
var WithAllCharsets = AllCharsetsOption{}

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

// ApplyConcatOption applies the CharsetOption to a decoding during
// Contatenate.
func (o CharsetOption) ApplyConcatOption(cc *ConcatConfig) {
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

// ApplyConcatOption applies the LockingCharsetOption to a decoding during
// Contatenate.
func (o LockingCharsetOption) ApplyConcatOption(cc *ConcatConfig) {
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

// ApplyConcatOption applies the ShiftCharsetOption to a decoding during
// Contatenate.
func (o ShiftCharsetOption) ApplyConcatOption(cc *ConcatConfig) {
	cc.dopts = append(cc.dopts, tpdu.WithShiftCharset(o.nli...))
}

// DirectionOption defines the direction.
type DirectionOption struct {
	d tpdu.Direction
}

// ApplyDecodeOption applies the ShiftCharsetOption to an Encoder.
func (o DirectionOption) ApplyDecodeOption(d *DecodeConfig) {
	d.dirn = o.d
}
