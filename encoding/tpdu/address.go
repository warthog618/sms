// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

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

// AddressOption returns a new Address with an option applied.
type AddressOption func(Address) Address

// NewAddress creates an Address and initialises the TOA.
func NewAddress(options ...AddressOption) Address {
	a := Address{}
	for _, option := range options {
		a = option(a)
	}
	// TOA must have bit 7 set.
	a.TOA |= 0x80
	return a
}

// FromNumber creates an AddressOption thats sets the address to the
// international number.
//
// The number may be optionally prefixed with '+'.
func FromNumber(number string) AddressOption {
	return func(a Address) Address {
		a.SetNumber(number)
		return a
	}
}

// MarshalBinary marshals an Address into binary.
//
// It returns the marshalled address and any error detected
// while marshalling.
func (a *Address) MarshalBinary() (dst []byte, err error) {
	ton := a.TypeOfNumber()
	var addr []byte
	var l int // is digits and ignores the toa
	switch ton {
	case TonAlphanumeric:
		e := gsm7.NewEncoder().WithExtCharset(nil) // without escapes
		addr, err = e.Encode([]byte(a.Addr))
		if err != nil {
			return nil, EncodeError("addr", err)
		}
		l = (len(addr)*7 + 3) / 4
		addr = gsm7.Pack7Bit(addr, 0)
	default:
		addr, err = semioctet.Encode([]byte(a.Addr))
		if err != nil {
			return nil, EncodeError("addr", err)
		}
		l = len(a.Addr)
	}
	dst = make([]byte, 0, l+2)
	dst = append(dst, byte(l), a.TOA)
	dst = append(dst, addr...)
	return dst, nil
}

// UnmarshalBinary unmarshals an Address from a binary TPDU.
//
// It returns the number of bytes read from the source, and any error detected
// while unmarshalling.
func (a *Address) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 2 {
		return 0, NewDecodeError("addr", 0, ErrUnderflow)
	}
	l := int(src[0])  // len is semi-octets and ignores toa
	ol := (l + 1) / 2 // octet length
	toa := src[1]
	ton := TypeOfNumber((toa >> 4) & 0x07)
	ri := 2
	if len(src) < ri+ol {
		return len(src), NewDecodeError("addr", ri, ErrUnderflow)
	}
	switch ton {
	case TonAlphanumeric:
		u := gsm7.Unpack7Bit(src[ri:ri+ol], 0)
		if (len(u)*7+3)/4 > l {
			// drop septet of fill
			u = u[:len(u)-1]
		}
		d := gsm7.NewDecoder().WithExtCharset(nil).Strict() // without escapes
		baddr, err := d.Decode(u)
		if err != nil {
			return ri, NewDecodeError("addr", ri, err)
		}
		ri += ol
		a.Addr = string(baddr)
	default:
		baddr, n, err := semioctet.Decode(make([]byte, l), src[ri:ri+ol])
		ri += n
		if err != nil {
			return ri, NewDecodeError("addr", ri-n, err)
		}
		if n != ol || len(baddr) < l {
			return ri, NewDecodeError("addr", ri-n, ErrUnderflow)
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

// SetNumber sets the address to the international number.
//
// The number may be optionally prefixed with '+'.
func (a *Address) SetNumber(number string) {
	if len(number) > 0 && number[0] == '+' {
		number = number[1:]
	}
	a.SetTypeOfNumber(TonInternational)
	a.SetNumberingPlan(NpISDN)
	a.Addr = number
}

// NumberingPlan extracts the NPI field from the TOA.
func (a Address) NumberingPlan() NumberingPlan {
	return NumberingPlan(a.TOA & 0x0f)
}

// SetNumberingPlan sets the NPI field in the TOA.
func (a *Address) SetNumberingPlan(np NumberingPlan) {
	a.TOA = (a.TOA &^ 0x0f) | byte(np&0x0f)
}

// SetTypeOfNumber sets the TON field in the TOA.
func (a *Address) SetTypeOfNumber(ton TypeOfNumber) {
	a.TOA = (a.TOA &^ 0x70) | (byte(ton&0x7) << 4)
}

// TypeOfNumber extracts the TON field from the TOA.
func (a Address) TypeOfNumber() TypeOfNumber {
	return TypeOfNumber((a.TOA >> 4) & 0x07)
}

// TypeOfNumber corresponds to bits 6,5,4 of the Address TOA field.
//
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
