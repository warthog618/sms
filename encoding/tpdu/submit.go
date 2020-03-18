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

// SubmitOption is a construction option for Submit PDUs.
type SubmitOption interface {
	applySubmitOption(*Submit)
}

// NewSubmit creates a Submit TPDU and initialises non-zero fields.
func NewSubmit(options ...SubmitOption) *Submit {
	s := &Submit{TPDU: TPDU{FirstOctet: byte(MtSubmit)}}
	for _, option := range options {
		option.applySubmitOption(s)
	}
	return s
}

// FromSubmitOption provides a template PDU to be copied into the new PDU.
type FromSubmitOption struct {
	t Submit
}

func (o FromSubmitOption) applySubmitOption(s *Submit) {
	s.Clone(&o.t)
}

// FromSubmit provides a template PDU to be copied into the new PDU.
//
// The template PDU is copied.
func FromSubmit(t *Submit) FromSubmitOption {
	s := Submit{}
	s = *t
	s.UDH = append(t.UDH[:0:0], t.UDH...)
	return FromSubmitOption{s}
}

func (udh UserDataHeader) applySubmitOption(s *Submit) {
	s.SetUDH(append(s.UDH, udh...))
}

// WithUserDataHeader provides a custom base user data header for the Submit PDU.
//
// This user header is copied and may be extended during user data encoding.
func WithUserDataHeader(udh UserDataHeader) UserDataHeader {
	return udh
}

func (alpha Alphabet) applySubmitOption(s *Submit) {
	dcs, err := DCS(s.DCS).WithAlphabet(alpha)
	if err != nil {
		// ignore the template dcs
		dcs, _ = DCS(0).WithAlphabet(alpha)
	}
	s.DCS = dcs
}

// WithAlphabet specifies the alphabet used to encode user data in the Submit PDU.
func WithAlphabet(alpha Alphabet) Alphabet {
	return alpha
}

// To specifies the destination number the PDU is addressed to.
//
// The number is assumed to be in international format.
func To(number string) AddressOption {
	if len(number) > 0 && number[0] == '+' {
		number = number[1:]
	}
	return AddressOption{
		Address{
			TOA:  0x80 | byte(TonInternational<<4) | byte(NpISDN),
			Addr: number,
		}}
}

// AddressOption provides a destination address for the Submit PDU.
type AddressOption struct {
	addr Address
}

func (o AddressOption) applySubmitOption(s *Submit) {
	s.DA = o.addr
}

// Clone copies the Submit PDU attribules into a target.
func (s *Submit) Clone(t *Submit) {
	*s = *t
	s.UDH = append(t.UDH[:0:0], t.UDH...)
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
	b = append(b, s.PID, byte(s.DCS))
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
	s.DCS = DCS(src[ri])
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
