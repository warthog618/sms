// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// SubmitReport represents a SMS-Submit-Report PDU as defined in 3GPP TS 23.038
// Section 9.2.2.2a.
type SubmitReport struct {
	TPDU
	FCS  byte
	PI   byte
	SCTS Timestamp
}

// NewSubmitReport creates a SubmitReport TPDU and initialises non-zero fields.
func NewSubmitReport() *SubmitReport {
	return &SubmitReport{TPDU: TPDU{FirstOctet: byte(MtSubmit)}}
}

// SetDCS sets the SubmitReport dcs field and the corresponding bit of the pi.
func (s *SubmitReport) SetDCS(dcs byte) {
	s.PI = s.PI | 0x02
	s.TPDU.DCS = dcs
}

// SetPID sets the SubmitReport pid field and the corresponding bit of the pi.
func (s *SubmitReport) SetPID(pid byte) {
	s.PI = s.PI | 0x01
	s.TPDU.PID = pid
}

// SetUD sets the SubmitReport ud field and the corresponding bit of the pi.
func (s *SubmitReport) SetUD(ud UserData) {
	s.PI = s.PI | 0x04
	s.TPDU.UD = ud
}

// SetUDH sets the User Data Header of the SubmitReport and the corresponding
// bit of the pi.
func (s *SubmitReport) SetUDH(udh UserDataHeader) {
	s.PI = s.PI | 0x04
	s.TPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Submit-Report TPDU.
func (s *SubmitReport) MarshalBinary() ([]byte, error) {
	b := []byte{s.FirstOctet, s.FCS, s.PI}
	scts, err := s.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	if s.PI&0x01 == 0x01 {
		b = append(b, s.PID)
	}
	if s.PI&0x02 == 0x02 {
		b = append(b, s.DCS)
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

// UnmarshalBinary unmarshals an SMS-Submit-Report TPDU.
func (s *SubmitReport) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	s.FirstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	s.FCS = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	s.PI = src[ri]
	ri++
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err := s.SCTS.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
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
		s.DCS = src[ri]
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

func decodeSubmitReport(src []byte) (interface{}, error) {
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
