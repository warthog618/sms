// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// DeliverReport represents a SMS-Deliver-Report PDU as defined in 3GPP TS 23.038 Section 9.2.2.1a.
type DeliverReport struct {
	BaseTPDU
	fcs byte
	pi  byte
}

// NewDeliverReport creates a DeliverReport TPDU and initialises non-zero fields.
func NewDeliverReport() *DeliverReport {
	return &DeliverReport{BaseTPDU: BaseTPDU{firstOctet: byte(MtDeliver), udhiMask: 0x04}}
}

// FCS returns the DeliverReport fcs.
func (d *DeliverReport) FCS() byte {
	return d.fcs
}

// PI returns the DeliverReport pi.
func (d *DeliverReport) PI() byte {
	return d.pi
}

// SetDCS sets the DeliverReport dcs field and the corresponding bit of the pi.
func (d *DeliverReport) SetDCS(dcs DCS) {
	d.pi = d.pi | 0x02
	d.BaseTPDU.SetDCS(dcs)
}

// SetFCS sets the DeliverReport fcs field.
func (d *DeliverReport) SetFCS(fcs byte) {
	d.fcs = fcs
}

// SetPI sets the DeliverReport pi field.
func (d *DeliverReport) SetPI(pi byte) {
	d.pi = pi
}

// SetPID sets the DeliverReport pid field and the corresponding bit of the pi.
func (d *DeliverReport) SetPID(pid byte) {
	d.pi = d.pi | 0x01
	d.BaseTPDU.SetPID(pid)
}

// SetUD sets the DeliverReport ud field and the corresponding bit of the pi.
func (d *DeliverReport) SetUD(ud UserData) {
	d.pi = d.pi | 0x04
	d.BaseTPDU.SetUD(ud)
}

// SetUDH sets the User Data Header of the DeliverReport and the corresponding bit of the pi.
func (d *DeliverReport) SetUDH(udh UserDataHeader) {
	d.pi = d.pi | 0x04
	d.BaseTPDU.SetUDH(udh)
}

// MarshalBinary marshals an SMS-Deliver-Report TPDU.
func (d *DeliverReport) MarshalBinary() ([]byte, error) {
	b := []byte{d.firstOctet, d.fcs, d.pi}
	if d.pi&0x01 == 0x01 {
		b = append(b, d.pid)
	}
	if d.pi&0x02 == 0x02 {
		b = append(b, d.dcs)
	}
	if d.pi&0x4 == 0x4 {
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
	d.firstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	d.fcs = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	d.pi = src[ri]
	ri++
	if d.pi&0x01 == 0x01 {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		d.pid = src[ri]
		ri++
	}
	if d.pi&0x02 == 0x02 {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		d.dcs = src[ri]
		ri++
	}
	d.udhiMask = 0x04
	if d.pi&0x04 == 0x04 {
		err := d.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func decodeDeliverReport(src []byte) (TPDU, error) {
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
