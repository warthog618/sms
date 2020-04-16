// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package pdumode

import (
	"encoding/hex"
)

// PDU represents the PDU exchanged with the GSM modem.
type PDU struct {
	// SMCS Address
	SMSC SMSCAddress

	// TPDU in binary form
	TPDU []byte
}

// UnmarshalBinary decodes the binary form of the PDU provided by the modem.
//
// Returns the unmarshalled PDU, or an error if unmarshalling fails.
func UnmarshalBinary(src []byte) (p *PDU, err error) {
	pdu := PDU{}
	err = pdu.UnmarshalBinary(src)
	if err != nil {
		return
	}
	p = &pdu
	return
}

// UnmarshalHexString decodes the hex string provided by the modem.
func UnmarshalHexString(s string) (p *PDU, err error) {
	pdu := PDU{}
	err = pdu.UnmarshalHexString(s)
	if err != nil {
		return
	}
	p = &pdu
	return
}

// UnmarshalBinary decodes the binary form of the PDU provided by the modem.
func (p *PDU) UnmarshalBinary(src []byte) error {
	n, err := p.SMSC.UnmarshalBinary(src)
	if err != nil {
		return err
	}
	p.TPDU = src[n:]
	return nil
}

// UnmarshalHexString decodes the hex string provided by the modem.
func (p *PDU) UnmarshalHexString(s string) error {
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	return p.UnmarshalBinary(b)
}

// MarshalBinary marshals the PDU into binary form.
func (p *PDU) MarshalBinary() ([]byte, error) {
	dst, err := p.SMSC.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dst = append(dst, p.TPDU...)
	return dst, nil
}

// MarshalHexString encodes the PDU into the hex string expected by the modem.
func (p *PDU) MarshalHexString() (string, error) {
	b, err := p.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
