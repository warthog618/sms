// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"github.com/warthog618/sms/encoding/gsm7"
)

// BaseTPDU is the base type for SMS TPDUs.
type BaseTPDU struct {
	firstOctet byte
	pid        byte
	dcs        byte
	udh        UserDataHeader
	// ud contains the short message from the User Data.
	// It does not include the User Data Header, which is provided in udh.
	// The interpretation of ud depends on the Alphabet.
	// For Alpha7Bit, ud is an array of GSM7 septets, each septet stored in the lower 7 bits of a byte.
	//  These have NOT been converted to the corresponding UTF8.
	//  Use the gsm7 package to convert to UTF8.
	// For AlphaUCS2, ud is an array of UCS2 characters packed into a byte array in Big Endian.
	//  These have NOT been converted to the corresponding UTF8.
	//  Use the usc2 package to convert to UTF8.
	// For Alpha8Bit, ud contains the raw octets.
	ud       UserData
	udhiMask byte
}

// Alphabet returns the alphabet field from the DCS of the SMS TPDU.
func (p *BaseTPDU) Alphabet() (Alphabet, error) {
	return DCS(p.dcs).Alphabet()
}

// DCS returns the TPDU dcs field.
func (p *BaseTPDU) DCS() DCS {
	return DCS(p.dcs)
}

// FirstOctet returns the TPDU firstOctet field.
func (p *BaseTPDU) FirstOctet() byte {
	return p.firstOctet
}

// PID returns the TPDU pid field.
func (p *BaseTPDU) PID() byte {
	return p.pid
}

// MTI returns the MessageType from the first octet of the SMS TPDU.
func (p *BaseTPDU) MTI() MessageType {
	return MessageType(p.firstOctet & 0x3)
}

// SetDCS sets the TPDU dcs field.
func (p *BaseTPDU) SetDCS(dcs DCS) {
	p.dcs = byte(dcs)
}

// SetFirstOctet sets the TPDU firstOctet field.
func (p *BaseTPDU) SetFirstOctet(fo byte) {
	p.firstOctet = fo
}

// SetPID sets the TPDU pid field.
func (p *BaseTPDU) SetPID(pid byte) {
	p.pid = pid
}

// SetUD sets the TPDU ud field.
func (p *BaseTPDU) SetUD(ud UserData) {
	p.ud = ud
}

// SetUDH sets the User Data Header of the TPDU.
func (p *BaseTPDU) SetUDH(udh UserDataHeader) {
	if udh == nil {
		p.udh = nil
		p.firstOctet = (p.firstOctet &^ p.udhiMask)
	} else {
		p.udh = udh
		p.firstOctet = (p.firstOctet | p.udhiMask)
	}
}

// SetUDHIMask sets the udhiMask field of the TPDU.
func (p *BaseTPDU) SetUDHIMask(udhiMask byte) {
	p.udhiMask = udhiMask
}

// UD returns the User Data.
func (p *BaseTPDU) UD() UserData {
	return p.ud
}

// UDH returns the User Data Header.
func (p *BaseTPDU) UDH() UserDataHeader {
	return p.udh
}

// UDHI returns the User Data Header Indicator bit from the SMS TPDU first octet.
// This is generally the same as testing the length of the udh - unless the dcs
// has been intentionally overwritten to create an inconsistency.
func (p *BaseTPDU) UDHI() bool {
	return p.firstOctet&p.udhiMask != 0
}

// decodeUserData unmarshals the User Data field from the binary src.
func (p *BaseTPDU) decodeUserData(src []byte) error {
	if len(src) < 1 {
		return DecodeError("udl", 0, ErrUnderflow)
	}
	udl := int(src[0])
	if udl == 0 {
		return nil
	}
	var udh UserDataHeader
	sml7 := 0
	ri := 1
	alphabet, err := p.Alphabet()
	if err != nil {
		return DecodeError("alphabet", ri, err)
	}
	if alphabet == Alpha7Bit {
		sml7 = udl
		// length is septets - convert to octets
		udl = (sml7*7 + 7) / 8
	}
	if len(src) < ri+udl {
		return DecodeError("sm", ri, ErrUnderflow)
	}
	if len(src) > ri+udl {
		return DecodeError("ud", ri, ErrOverlength)
	}
	var udhl int // Note that in this context udhl includes itself.
	udhi := p.UDHI()
	if udhi {
		udh = make(UserDataHeader, 0)
		l, err := udh.UnmarshalBinary(src[ri:])
		if err != nil {
			return DecodeError("udh", ri, err)
		}
		udhl = l
		ri += udhl
	}
	if ri == len(src) {
		p.udh = udh
		return nil
	}
	switch alphabet {
	case Alpha7Bit:
		sm, err := decode7Bit(sml7, udhl, src[ri:])
		if err != nil {
			return DecodeError("sm", ri, err)
		}
		p.ud = sm
	case AlphaUCS2:
		if len(src[ri:])&0x01 == 0x01 {
			return DecodeError("sm", ri, ErrOverlength)
		}
		fallthrough
	case Alpha8Bit:
		p.ud = append([]byte(nil), src[ri:]...)
	}
	p.udh = udh
	return nil
}

// decode7Bit decodes the GSM7 encoded binary src into a byte array.
// sml is the number of septets expected, and udhl is the number of octets
// in the UDH, including the UDHL field.
func decode7Bit(sml, udhl int, src []byte) ([]byte, error) {
	var fillBits int
	if udhl > 0 {
		if dangling := udhl % 7; dangling != 0 {
			fillBits = 7 - dangling
		}
		sml = sml - (udhl*8+fillBits)/7
	}
	sm := gsm7.Unpack7Bit(src, fillBits)
	// this is a double check on the math and should never trip...
	if len(sm) < sml {
		return nil, ErrUnderflow
	}
	if len(sm) > sml {
		if len(sm) > sml+1 || sm[sml] != 0 {
			return nil, ErrOverlength
		}
		// drop trailing 0 septet
		sm = sm[:sml]
	}
	return sm, nil
}

// encodeUserData marshals the User Data into binary.
// The User Data Header is also encoded if present.
// If Alphabet is GSM7 then the User Data is assumed to be unpacked GSM7 septets
// and is packed prior to encoding.
// For other alphabet values the User Data is encoded as is.
// No checks of encoded size are performed here as that depends on concrete TPDU type,
// and that can check the length of the returned b.
func (p *BaseTPDU) encodeUserData() (b []byte, err error) {
	udh, err := p.udh.MarshalBinary()
	if err != nil {
		return nil, EncodeError("udh", err)
	}
	ud := p.ud
	alphabet, err := p.Alphabet()
	if err != nil {
		return nil, EncodeError("alphabet", err)
	}
	udl := len(p.ud) // assume octets
	switch alphabet {
	case Alpha7Bit:
		fillBits := 0
		if dangling := len(udh) % 7; dangling != 0 {
			fillBits = 7 - dangling
		}
		ud = gsm7.Pack7Bit(p.ud, fillBits)
		// udl is in septets so convert
		if udl > 0 {
			udl = udl + (len(udh)*8+fillBits)/7
		} else {
			udl = (len(udh) * 8) / 7
		}
	case AlphaUCS2:
		if udl&0x01 == 0x01 {
			return nil, EncodeError("sm", ErrOddUCS2Length)
		}
		fallthrough
	case Alpha8Bit:
		// udl is in octets
		udl = udl + len(udh)
	}
	b = make([]byte, 0, 1+len(udh)+len(ud))
	b = append(b, byte(udl))
	b = append(b, udh...)
	b = append(b, ud...)
	return b, nil
}

// MessageType identifies the type of TPDU encoded in a binary stream,
// as defined in 3GPP TS 23.040 Section 9.2.3.1.
// Note that the direction of the TPDU must also be known to determine
// how to interpret the TPDU.
type MessageType int

const (
	// MtDeliver identifies the message as a SMS-Deliver or SMS-Deliver-Report TPDU.
	MtDeliver MessageType = iota
	// MtSubmit identifies the message as a SMS-Submit or SMS-Submit-Report TPDU.
	MtSubmit
	// MtCommand identifies the message as a SMS-Command or SMS-Status-Report TPDU.
	MtCommand
	// MtReserved identifies the message as an unknown type of SMS TPDU.
	MtReserved
)

// Direction indicates the direction that the SMS TPDU is carried.
type Direction int

const (
	// MT indicates that the SMS TPDU is intended to be received by the MS.
	MT Direction = iota
	// MO indicates that the SMS TPDU is intended to be sent by the MS.
	MO
)
