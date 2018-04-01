// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pdumode

import (
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
)

// SMSCAddress is the address of the SMSC.
// The SMCSAddress is similar to a TPDU Address, but the binary form is
// marshalled differently, hence the subtype.
// The Type-of-number should typically be TonNational or TonInternational,
// but that is not enforced.
// The NumberingPlan should typically be NpISDN, but that is not enforced either.
type SMSCAddress tpdu.Address

// MarshalBinary marshals the SMSC Address into binary.
func (a *SMSCAddress) MarshalBinary() (dst []byte, err error) {
	addr, err := semioctet.Encode([]byte(a.Addr))
	if err != nil {
		return nil, tpdu.EncodeError("addr", err)
	}
	l := len(addr) + 1 // in octets and includes the toa
	dst = make([]byte, 2, l+1)
	dst[0] = byte(l)
	dst[1] = a.TOA
	dst = append(dst, addr...)
	return dst, nil
}

// UnmarshalBinary unmarshals an SMSC Address from a TPDU field.
// It returns the number of bytes read from the source, and any error detected
// while decoding.
func (a *SMSCAddress) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 1 {
		return 0, tpdu.DecodeError("length", 0, tpdu.ErrUnderflow)
	}
	l := int(src[0]) // len is octets including toa
	if l < 1 {
		return 1, tpdu.DecodeError("length", 0, tpdu.ErrUnderflow)
	}
	if len(src) < 2 {
		return 1, tpdu.DecodeError("toa", 1, tpdu.ErrUnderflow)
	}
	toa := src[1]
	ri := 2
	l-- // encoded length includes toa
	if len(src) < ri+l {
		return len(src), tpdu.DecodeError("addr", ri, tpdu.ErrUnderflow)
	}
	baddr, n, err := semioctet.Decode(make([]byte, l*2), src[ri:ri+l])
	ri += n
	if err != nil {
		return ri, tpdu.DecodeError("addr", ri-n, err)
	}
	a.Addr = string(baddr)
	a.TOA = toa
	return ri, nil
}
