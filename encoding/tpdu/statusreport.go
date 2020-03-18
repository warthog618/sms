// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// StatusReport represents a SMS-Status-Report PDU as defined in 3GPP TS 23.038
// Section 9.2.2.3.
type StatusReport struct {
	TPDU
	MR   byte
	RA   Address
	SCTS Timestamp
	DT   Timestamp
	ST   byte
	PI   byte
}

// NewStatusReport creates a StatusReport TPDU and initialises non-zero fields.
func NewStatusReport() *StatusReport {
	return &StatusReport{TPDU: TPDU{FirstOctet: byte(MtCommand)}}
}

// SetDCS sets the StatusReport dcs field and the corresponding bit of the PI.
func (s *StatusReport) SetDCS(dcs byte) {
	s.PI = s.PI | 0x02
	s.TPDU.DCS = DCS(dcs)
}

// SetPID sets the StatusReport pid field and the corresponding bit of the PI.
func (s *StatusReport) SetPID(pid byte) {
	s.PI = s.PI | 0x01
	s.TPDU.PID = pid
}

// SetUD sets the StatusReport ud field and the corresponding bit of the PI.
func (s *StatusReport) SetUD(ud UserData) {
	s.PI = s.PI | 0x04
	s.TPDU.UD = ud
}

// SetUDH sets the User Data Header of the StatusReport and the corresponding
// bit of the PI.
func (s *StatusReport) SetUDH(udh UserDataHeader) {
	s.PI = s.PI | 0x04
	s.TPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Status-Report TPDU.
func (s *StatusReport) MarshalBinary() ([]byte, error) {
	b := []byte{s.FirstOctet, s.MR}
	ra, err := s.RA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("ra", err)
	}
	b = append(b, ra...)
	scts, err := s.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	dt, err := s.DT.MarshalBinary()
	if err != nil {
		return nil, EncodeError("dt", err)
	}
	b = append(b, dt...)
	b = append(b, s.ST)
	if s.PI == 0x00 {
		return b, nil
	}
	b = append(b, s.PI)
	if s.PI&0x01 == 0x01 {
		b = append(b, s.PID)
	}
	if s.PI&0x02 == 0x02 {
		b = append(b, byte(s.DCS))
	}
	if s.PI&0x4 == 0x4 {
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
	s.FirstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	s.MR = src[ri]
	ri++
	n, err := s.RA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("ra", ri, err)
	}
	ri += n
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err = s.SCTS.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	if len(src) < ri+7 {
		return DecodeError("dt", ri, ErrUnderflow)
	}
	err = s.DT.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("dt", ri, err)
	}
	ri += 7
	if len(src) <= ri {
		return DecodeError("st", ri, ErrUnderflow)
	}
	s.ST = src[ri]
	ri++
	if len(src) > ri {
		return s.unmarshalOptionals(ri, src)
	}
	return nil
}

// unmarshal the optional fields at the end if the StatusReport TPDU.
func (s *StatusReport) unmarshalOptionals(ri int, src []byte) error {
	s.PI = src[ri]
	ri++
	if s.PI&0x01 == 0x01 {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		s.PID = src[ri]
		ri++
	}
	if s.PI&0x02 == 0x02 {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		s.DCS = DCS(src[ri])
		ri++
	}
	if s.PI&0x04 == 0x04 {
		err := s.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func decodeStatusReport(src []byte) (interface{}, error) {
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
