// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Submit represents a SMS-Submit PDU as defined in 3GPP TS 23.038 Section 9.2.2.2.
type Submit struct {
	TPDU
	MR byte
	DA Address
	VP ValidityPeriod
}

// NewSubmit creates a Submit TPDU and initialises non-zero fields.
func NewSubmit() *Submit {
	return &Submit{TPDU: TPDU{FirstOctet: byte(MtSubmit)}}
}

// MaxUDL returns the maximum number of octets that can be encoded into the UD.
// Note that for 7bit encoding this can result in up to 160 septets.
func (s *Submit) MaxUDL() int {
	return 140
}

// SetVP sets the validity period and the corresponding VPF bits
// in the firstOctet.
func (s *Submit) SetVP(vp ValidityPeriod) {
	s.FirstOctet = s.FirstOctet&^0x0c | byte(vp.Format<<2)
	s.VP = vp
}

// MarshalBinary marshals an SMS-Submit TPDU.
func (s *Submit) MarshalBinary() ([]byte, error) {
	b := []byte{s.FirstOctet, s.MR}
	da, err := s.DA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	b = append(b, da...)
	b = append(b, s.PID, s.DCS)
	if s.VP.Format != VpfNotPresent {
		vp, verr := s.VP.MarshalBinary()
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
// In the case of error the Submit will be partially unmarshalled, up to the
// point that the decoding error was detected.
func (s *Submit) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	s.FirstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	s.MR = src[ri]
	ri++
	n, err := s.DA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("da", ri, err)
	}
	ri += n
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	s.PID = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("dcs", ri, ErrUnderflow)
	}
	s.DCS = src[ri]
	ri++
	vpf := ValidityPeriodFormat((s.FirstOctet >> 3) & 0x3)
	n, err = s.VP.UnmarshalBinary(src[ri:], vpf)
	if err != nil {
		return DecodeError("vp", ri, err)
	}
	ri += n
	err = s.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func decodeSubmit(src []byte) (interface{}, error) {
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
