// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"github.com/warthog618/sms/encoding/gsm7"
	"github.com/warthog618/sms/encoding/semioctet"
)

// Address represents a phone number.
type Address struct {
	TOA  byte
	Addr string
}

// NewAddress creates an Address and initialises the TOA.
func NewAddress() *Address {
	// TOA must have bit 7 set.
	a := Address{TOA: 0x80}
	return &a
}

// MarshalBinary marshals an Address into binary.
// It returns the marshalled address and any error detected
// while marshalling.
func (a *Address) MarshalBinary() (dst []byte, err error) {
	ton := TypeOfNumber((a.TOA >> 4) & 0x07)
	var addr []byte
	var l int // is digits and ignores the toa
	switch ton {
	case TonAlphanumeric:
		e := gsm7.NewEncoder().WithExtCharset(nil) // without escapes
		addr, err = e.Encode([]byte(a.Addr))
		if err != nil {
			return nil, EncodeError("addr", err)
		}
		l = len(addr)
		addr = gsm7.Pack7Bit(addr, 0)
	default:
		addr, err = semioctet.Encode([]byte(a.Addr))
		if err != nil {
			return nil, EncodeError("addr", err)
		}
		l = len(a.Addr)
	}
	dst = make([]byte, 2, l+2)
	dst[0] = byte(l)
	dst[1] = a.TOA
	dst = append(dst, addr...)
	return dst, nil
}

// UnmarshalBinary unmarshals an Address from a binary TPDU.
// It returns the number of bytes read from the source, and any error detected
// while unmarshalling.
func (a *Address) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 2 {
		return 0, DecodeError("addr", 0, ErrUnderflow)
	}
	l := int(src[0]) // len is digits and ignores toa
	toa := src[1]
	ton := TypeOfNumber((toa >> 4) & 0x07)
	ri := 2
	switch ton {
	case TonAlphanumeric:
		// l is digits, i.e. GSM7 septets, and so requires conversion to octets....
		ol := (l*7 + 7) / 8
		if len(src) < ri+ol {
			return len(src), DecodeError("addr", ri, ErrUnderflow)
		}
		u := gsm7.Unpack7Bit(src[ri:ri+ol], 0)
		d := gsm7.NewDecoder().WithExtCharset(nil).Strict() // without escapes
		baddr, err := d.Decode(u)
		if err != nil {
			return ri, DecodeError("addr", ri, err)
		}
		ri += ol
		a.Addr = string(baddr)
	default:
		sl := (l + 1) / 2
		if len(src) < ri+sl {
			return len(src), DecodeError("addr", ri, ErrUnderflow)
		}
		baddr, n, err := semioctet.Decode(make([]byte, l), src[ri:ri+sl])
		ri += n
		if err != nil {
			return ri, DecodeError("addr", ri-n, err)
		}
		if n != sl || len(baddr) < l {
			return ri, DecodeError("addr", ri-n, ErrUnderflow)
		}
		a.Addr = string(baddr)
	}
	a.TOA = toa
	return ri, nil
}

// Number returns the stringified number corresponding to the Address.
func (a Address) Number() string {
	if a.TypeOfNumber() == TonInternational {
		return "+" + a.Addr
	}
	return a.Addr
}

// NumberingPlan extracts the NPI field from the TOA.
func (a *Address) NumberingPlan() NumberingPlan {
	return NumberingPlan(a.TOA & 0x0f)
}

// SetNumberingPlan sets the NPI field in the TOA.
func (a *Address) SetNumberingPlan(np NumberingPlan) {
	a.TOA = (a.TOA &^ 0x0f) | byte(np&0x0f)
}

// SetTypeOfNumber sets the TON field in the TOA.
func (a *Address) SetTypeOfNumber(ton TypeOfNumber) {
	a.TOA = (a.TOA &^ 0x30) | (byte(ton&0x3) << 4)
}

// TypeOfNumber extracts the TON field from the TOA.
func (a *Address) TypeOfNumber() TypeOfNumber {
	return TypeOfNumber((a.TOA >> 4) & 0x07)
}

// TypeOfNumber corresponds to bits 6,5,4 of the Address TOA field.
// i.e. 1xxxyyyy,
// as defined in 3GPP TS 23.040 Section 9.1.2.5.
type TypeOfNumber int

const (
	// TonUnknown indicates the type of the number is unknown.
	TonUnknown TypeOfNumber = iota
	// TonInternational indicates the number is international.
	TonInternational
	// TonNational indicates the number is national.
	TonNational
	// TonNetworkSpecific indicates the number is specific to the carrier network.
	TonNetworkSpecific
	// TonSubscriberNumber indicates the number is a subscriber number.
	TonSubscriberNumber
	// TonAlphanumeric indicates the number is in alphanumeric format.
	TonAlphanumeric
	// TonAbbreviated indicates the number is in abbreviated format.
	TonAbbreviated
	// TonExtension is reserved for future extension.
	TonExtension
)

// NumberingPlan corresponds to bits 4,3,2,1 of the Address TOA field.
// i.e. 1yyyxxxx
// as defined in 3GPP TS 23.040 Section 9.1.2.5
type NumberingPlan int

const (
	// NpUnknown indicates the numbering plan is unknown.
	NpUnknown NumberingPlan = iota
	// NpISDN indicates the number is in ISDN/E.164 format.
	NpISDN
	_
	// NpData indicates a data numbering plan (X.121).
	NpData
	// NpTelex indicates a telex numbering plan.
	NpTelex
	// NpScSpecificA indicates a service center specific numbering plan.
	NpScSpecificA
	// NpScSpecificB indicates a service center specific numbering plan.
	NpScSpecificB
	_
	// NpNational indicates a national numbering plan.
	NpNational
	// NpPrivate indicates a private numbering plan.
	NpPrivate
	// NpErmes indicates the ERMES (ETSI DE/PS 3 01-3) numbering plan.
	NpErmes
	// NpExtension is reserved for future extensions.
	NpExtension = 0x0f
	// all other values reserved.
)
