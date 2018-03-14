// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// SubmitReport represents a SMS-Submit-Report PDU as defined in 3GPP TS 23.038 Section 9.2.2.2a.
type SubmitReport struct {
	BaseTPDU
	fcs  byte
	pi   byte
	scts Timestamp
}

// NewSubmitReport creates a SubmitReport TPDU and initialises non-zero fields.
func NewSubmitReport() *SubmitReport {
	return &SubmitReport{BaseTPDU: BaseTPDU{firstOctet: byte(MtSubmit), udhiMask: 0x04}}
}

// FCS returns the SubmitReport fcs.
func (s *SubmitReport) FCS() byte {
	return s.fcs
}

// PI returns the SubmitReport pi.
func (s *SubmitReport) PI() byte {
	return s.pi
}

// SCTS returns the SubmitReport scts.
func (s *SubmitReport) SCTS() Timestamp {
	return s.scts
}

// SetDCS sets the SubmitReport dcs field and the corresponding bit of the pi.
func (s *SubmitReport) SetDCS(dcs DCS) {
	s.pi = s.pi | 0x02
	s.BaseTPDU.SetDCS(dcs)
}

// SetFCS sets the SubmitReport fcs field.
func (s *SubmitReport) SetFCS(fcs byte) {
	s.fcs = fcs
}

// SetPI sets the SubmitReport pi field.
func (s *SubmitReport) SetPI(pi byte) {
	s.pi = pi
}

// SetPID sets the SubmitReport pid field and the corresponding bit of the pi.
func (s *SubmitReport) SetPID(pid byte) {
	s.pi = s.pi | 0x01
	s.BaseTPDU.SetPID(pid)
}

// SetSCTS sets the SubmitReport scts field.
func (s *SubmitReport) SetSCTS(scts Timestamp) {
	s.scts = scts
}

// SetUD sets the SubmitReport ud field and the corresponding bit of the pi.
func (s *SubmitReport) SetUD(ud UserData) {
	s.pi = s.pi | 0x04
	s.BaseTPDU.SetUD(ud)
}

// SetUDH sets the User Data Header of the SubmitReport and the corresponding bit of the pi.
func (s *SubmitReport) SetUDH(udh UserDataHeader) {
	s.pi = s.pi | 0x04
	s.BaseTPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Submit-Report TPDU.
func (s *SubmitReport) MarshalBinary() ([]byte, error) {
	b := []byte{s.firstOctet, s.fcs, s.pi}
	scts, err := s.scts.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	if s.pi&0x01 == 0x01 {
		b = append(b, s.pid)
	}
	if s.pi&0x02 == 0x02 {
		b = append(b, byte(s.dcs))
	}
	if s.pi&0x4 == 0x4 {
		ud, err := s.encodeUserData()
		if err != nil {
			return nil, EncodeError("ud", err)
		}
		b = append(b, ud...)
	}
	return b, nil
}

// UnmarshalBinary unmarshals an SMS-Submit-Report TPDU.
func (s *SubmitReport) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	s.firstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	s.fcs = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	s.pi = src[ri]
	ri++
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err := s.scts.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	if s.pi&0x01 == 0x01 {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		s.pid = src[ri]
		ri++
	}
	if s.pi&0x02 == 0x02 {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		s.SetDCS(DCS(src[ri]))
		ri++
	}
	s.udhiMask = 0x04
	if s.pi&0x04 == 0x04 {
		err := s.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func decodeSubmitReport(src []byte) (TPDU, error) {
	s := NewSubmitReport()
	if err := s.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return s, nil
}

// RegisterSubmitReportDecoder registers a decoder for the SubmitReport TPDU.
func RegisterSubmitReportDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtSubmit, MT, decodeSubmitReport)
}
