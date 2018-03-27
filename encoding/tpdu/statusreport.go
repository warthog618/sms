// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// StatusReport represents a SMS-Status-Report PDU as defined in 3GPP TS 23.038 Section 9.2.2.3.
type StatusReport struct {
	BaseTPDU
	mr   byte
	ra   Address
	scts Timestamp
	dt   Timestamp
	st   byte
	pi   byte
}

// NewStatusReport creates a StatusReport TPDU and initialises non-zero fields.
func NewStatusReport() *StatusReport {
	return &StatusReport{BaseTPDU: BaseTPDU{firstOctet: byte(MtCommand), udhiMask: 0x04}}
}

// RA returns the StatusReport ra.
func (s *StatusReport) RA() Address {
	return s.ra
}

// MR returns the StatusReport mr.
func (s *StatusReport) MR() byte {
	return s.mr
}

// PI returns the StatusReport pi.
func (s *StatusReport) PI() byte {
	return s.pi
}

// SCTS returns the StatusReport scts.
func (s *StatusReport) SCTS() Timestamp {
	return s.scts
}

// DT returns the StatusReport dt.
func (s *StatusReport) DT() Timestamp {
	return s.dt
}

// ST returns the StatusReport st.
func (s *StatusReport) ST() byte {
	return s.st
}

// SetDCS sets the StatusReport dcs field and the corresponding bit of the pi.
func (s *StatusReport) SetDCS(dcs DCS) {
	s.pi = s.pi | 0x02
	s.BaseTPDU.SetDCS(dcs)
}

// SetDT sets the StatusReport dt field.
func (s *StatusReport) SetDT(dt Timestamp) {
	s.dt = dt
}

// SetPI sets the StatusReport pi field.
func (s *StatusReport) SetPI(pi byte) {
	s.pi = pi
}

// SetMR sets the StatusReport mr field.
func (s *StatusReport) SetMR(mr byte) {
	s.mr = mr
}

// SetRA sets the StatusReport ra field.
func (s *StatusReport) SetRA(ra Address) {
	s.ra = ra
}

// SetPID sets the StatusReport pid field and the corresponding bit of the pi.
func (s *StatusReport) SetPID(pid byte) {
	s.pi = s.pi | 0x01
	s.BaseTPDU.SetPID(pid)
}

// SetSCTS sets the StatusReport scts field.
func (s *StatusReport) SetSCTS(scts Timestamp) {
	s.scts = scts
}

// SetST sets the StatusReport st field.
func (s *StatusReport) SetST(st byte) {
	s.st = st
}

// SetUD sets the StatusReport ud field and the corresponding bit of the pi.
func (s *StatusReport) SetUD(ud UserData) {
	s.pi = s.pi | 0x04
	s.BaseTPDU.SetUD(ud)
}

// SetUDH sets the User Data Header of the StatusReport and the corresponding bit of the pi.
func (s *StatusReport) SetUDH(udh UserDataHeader) {
	s.pi = s.pi | 0x04
	s.BaseTPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Status-Report TPDU.
func (s *StatusReport) MarshalBinary() ([]byte, error) {
	b := []byte{s.firstOctet, s.mr}
	ra, err := s.ra.MarshalBinary()
	if err != nil {
		return nil, EncodeError("ra", err)
	}
	b = append(b, ra...)
	scts, err := s.scts.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	dt, err := s.dt.MarshalBinary()
	if err != nil {
		return nil, EncodeError("dt", err)
	}
	b = append(b, dt...)
	b = append(b, s.st)
	if s.pi == 0x00 {
		return b, nil
	}
	b = append(b, s.pi)
	if s.pi&0x01 == 0x01 {
		b = append(b, s.pid)
	}
	if s.pi&0x02 == 0x02 {
		b = append(b, s.dcs)
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

// UnmarshalBinary unmarshals an SMS-Status-Report TPDU.
func (s *StatusReport) UnmarshalBinary(src []byte) error {
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
	n, err := s.ra.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("ra", ri, err)
	}
	ri += n
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err = s.scts.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	if len(src) < ri+7 {
		return DecodeError("dt", ri, ErrUnderflow)
	}
	err = s.dt.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("dt", ri, err)
	}
	ri += 7
	if len(src) <= ri {
		return DecodeError("st", ri, ErrUnderflow)
	}
	s.st = src[ri]
	ri++
	if len(src) > ri {
		return s.unmarshalOptionals(ri, src)
	}
	return nil
}

// unmarshal the optional fields at the end if the StatusReport TPDU.
func (s *StatusReport) unmarshalOptionals(ri int, src []byte) error {
	s.pi = src[ri]
	ri++
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
		s.dcs = src[ri]
		ri++
	}
	if s.pi&0x04 == 0x04 {
		s.udhiMask = 0x04
		err := s.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func decodeStatusReport(src []byte) (TPDU, error) {
	s := NewStatusReport()
	if err := s.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return s, nil
}

// RegisterStatusReportDecoder registers a decoder for the StatusReport TPDU.
func RegisterStatusReportDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtCommand, MT, decodeStatusReport)
}
