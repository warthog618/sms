// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package pdumode

import (
	"encoding/hex"
)

// Decoder converts a PDU into the SMSC address and TPDU that it contains.
type Decoder struct{}

// Decode decodes the binary form of the PDU provided by the modem.
//
// The PDU is comprised of the SMSC address and the TPDU, which are returned.
// The TPDU is in binary form, ready to be unmarshalled.
func (Decoder) Decode(src []byte) (*SMSCAddress, []byte, error) {
	smsc := SMSCAddress{}
	n, err := smsc.UnmarshalBinary(src)
	if err != nil {
		return nil, nil, err
	}
	return &smsc, src[n:], nil
}

// DecodeString decodes the hex string provided by the modem.
//
// The PDU is comprised of the SMSC address and the TPDU, which are returned.
// The TPDU is in binary form, ready to be unmarshalled.
func (d Decoder) DecodeString(s string) (*SMSCAddress, []byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, err
	}
	return d.Decode(b)
}

// Decode decodes the binary form of the PDU provided by the modem.
//
// The PDU is comprised of the SMSC address and the TPDU, which are returned.
// The TPDU is in binary form, ready to be unmarshalled.
func Decode(src []byte) (*SMSCAddress, []byte, error) {
	d := Decoder{}
	return d.Decode(src)
}

// DecodeString decodes the hex string provided by the modem.
//
// The PDU is comprised of the SMSC address and the TPDU, which are returned.
// The TPDU is in binary form, ready to be unmarshalled.
func DecodeString(s string) (*SMSCAddress, []byte, error) {
	d := Decoder{}
	return d.DecodeString(s)
}

// Encoder converts an SMSC address and TPDU into a PDU.
type Encoder struct{}

// Encode marshals the PDU into binary form.
// The PDU is comprised of the SMSC address and the TPDU.
// The TPDU has been marshalled into binary and is ready to be transmitted.
func (Encoder) Encode(smsc SMSCAddress, tpdu []byte) ([]byte, error) {
	dst, err := smsc.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dst = append(dst, tpdu...)
	return dst, nil
}

// EncodeToString encodes the PDU into the hex string expected by the modem.
// The PDU is comprised of the SMSC address and the TPDU.
// The TPDU is in binary ready to be transmitted.
func (e Encoder) EncodeToString(smsc SMSCAddress, tpdu []byte) (string, error) {
	p, err := e.Encode(smsc, tpdu)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(p), nil
}

// Encode marshals the PDU into binary form.
// The PDU is comprised of the SMSC address and the TPDU.
// The TPDU has been marshalled into binary and is ready to be transmitted.
func Encode(smsc SMSCAddress, tpdu []byte) ([]byte, error) {
	e := Encoder{}
	return e.Encode(smsc, tpdu)
}

// EncodeToString encodes the PDU into the hex string expected by the modem.
// The PDU is comprised of the SMSC address and the TPDU.
// The TPDU is in binary ready to be transmitted.
func EncodeToString(smsc SMSCAddress, tpdu []byte) (string, error) {
	e := Encoder{}
	return e.EncodeToString(smsc, tpdu)
}
