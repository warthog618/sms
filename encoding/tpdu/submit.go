// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Submit represents a SMS-Submit PDU as defined in 3GPP TS 23.038 Section 9.2.2.2.
type Submit struct {
	BaseTPDU
	mr byte
	da Address
	vp ValidityPeriod
}

// NewSubmit creates a Submit TPDU and initialises non-zero fields.
func NewSubmit() *Submit {
	return &Submit{BaseTPDU: BaseTPDU{firstOctet: byte(MtSubmit), udhiMask: 0x40}}
}

// DA returns the Submit da.
func (s *Submit) DA() Address {
	return s.da
}

// MaxUDL returns the maximum number of octets that can be encoded into the UD.
// Note that for 7bit encoding this can result in up to 160 septets.
func (s *Submit) MaxUDL() int {
	return 140
}

// MR returns the Submit mr.
func (s *Submit) MR() byte {
	return s.mr
}

// SetDA sets the Submit oa field.
func (s *Submit) SetDA(da Address) {
	s.da = da
}

// SetMR sets the Submit mr field.
func (s *Submit) SetMR(mr byte) {
	s.mr = mr
}

// SetVP sets the validity period and the corresponding VPF bits
// in the firstOctet.
func (s *Submit) SetVP(vp ValidityPeriod) {
	s.firstOctet = s.firstOctet&^0x0c | byte(vp.Format<<2)
	s.vp = vp
}

// VP returns the Submit vp.
func (s *Submit) VP() ValidityPeriod {
	return s.vp
}

// MarshalBinary marshals an SMS-Submit TPDU.
func (s *Submit) MarshalBinary() ([]byte, error) {
	b := []byte{s.firstOctet, s.mr}
	da, err := s.da.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	b = append(b, da...)
	b = append(b, s.pid, s.dcs)
	if s.vp.Format != VpfNotPresent {
		vp, verr := s.vp.MarshalBinary()
		if verr != nil {
			return nil, EncodeError("vp", verr)
		}
		b = append(b, vp...)
	}
	ud, err := s.encodeUserData()
	if err != nil {
		return nil, EncodeError("ud", err)
	}
	b = append(b, ud...)
	return b, nil
}

// UnmarshalBinary unmarshals an SMS-Submit TPDU.
// In the case of error the Submit will be partially unmarshalled, up to
// the point that the decoding error was detected.
func (s *Submit) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	s.firstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	s.mr = src[ri]
	ri++
	n, err := s.da.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("da", ri, err)
	}
	ri += n
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	s.pid = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("dcs", ri, ErrUnderflow)
	}
	s.dcs = src[ri]
	ri++
	vpf := ValidityPeriodFormat((s.firstOctet >> 3) & 0x3)
	n, err = s.vp.UnmarshalBinary(src[ri:], vpf)
	if err != nil {
		return DecodeError("vp", ri, err)
	}
	ri += n
	err = s.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	s.udhiMask = 0x40
	return nil
}

func decodeSubmit(src []byte) (TPDU, error) {
	s := NewSubmit()
	if err := s.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return s, nil
}

// RegisterSubmitDecoder registers a decoder for the Submit TPDU.
func RegisterSubmitDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtSubmit, MO, decodeSubmit)
}
