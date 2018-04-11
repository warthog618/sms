// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// DeliverReport represents a SMS-Deliver-Report PDU as defined in 3GPP TS 23.038 Section 9.2.2.1a.
type DeliverReport struct {
	TPDU
	FCS byte
	PI  byte
}

// NewDeliverReport creates a DeliverReport TPDU and initialises non-zero fields.
func NewDeliverReport() *DeliverReport {
	return &DeliverReport{TPDU: TPDU{FirstOctet: byte(MtDeliver)}}
}

// SetDCS sets the DeliverReport dcs field and the corresponding bit of the pi.
func (d *DeliverReport) SetDCS(dcs byte) {
	d.PI = d.PI | 0x02
	d.TPDU.DCS = dcs
}

// SetPID sets the DeliverReport pid field and the corresponding bit of the pi.
func (d *DeliverReport) SetPID(pid byte) {
	d.PI = d.PI | 0x01
	d.TPDU.PID = pid
}

// SetUD sets the DeliverReport ud field and the corresponding bit of the pi.
func (d *DeliverReport) SetUD(ud UserData) {
	d.PI = d.PI | 0x04
	d.TPDU.UD = ud
}

// SetUDH sets the User Data Header of the DeliverReport and the corresponding bit of the pi.
func (d *DeliverReport) SetUDH(udh UserDataHeader) {
	d.PI = d.PI | 0x04
	d.TPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Deliver-Report TPDU.
func (d *DeliverReport) MarshalBinary() ([]byte, error) {
	b := []byte{d.FirstOctet, d.FCS, d.PI}
	if d.PI&0x01 == 0x01 {
		b = append(b, d.PID)
	}
	if d.PI&0x02 == 0x02 {
		b = append(b, d.DCS)
	}
	if d.PI&0x4 == 0x4 {
		ud, err := d.encodeUserData()
		if err != nil {
			return nil, EncodeError("ud", err)
		}
		b = append(b, ud...)
	}
	return b, nil
}

// UnmarshalBinary unmarshals an SMS-Deliver-Report TPDU.
func (d *DeliverReport) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	d.FirstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	d.FCS = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	d.PI = src[ri]
	ri++
	if d.PI&0x01 == 0x01 {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		d.PID = src[ri]
		ri++
	}
	if d.PI&0x02 == 0x02 {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		d.DCS = src[ri]
		ri++
	}
	if d.PI&0x04 == 0x04 {
		err := d.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func decodeDeliverReport(src []byte) (interface{}, error) {
	d := NewDeliverReport()
	if err := d.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return d, nil
}

// RegisterDeliverReportDecoder registers a decoder for the DeliverReport TPDU.
func RegisterDeliverReportDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtDeliver, MO, decodeDeliverReport)
}
