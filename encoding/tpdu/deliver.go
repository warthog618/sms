// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Deliver represents a SMS-Deliver PDU as defined in 3GPP TS 23.038 Section 9.2.2.1.
type Deliver struct {
	BaseTPDU
	oa Address
	// The scts timestamp indicates the time the SMS was sent.
	// The time is the originator's local time, the timezone of which may differ from the
	// receiver's.
	scts Timestamp
}

// NewDeliver creates a Deliver TPDU and initialises non-zero fields.
func NewDeliver() *Deliver {
	return &Deliver{BaseTPDU: BaseTPDU{firstOctet: byte(MtDeliver)}}
}

// OA returns the Deliver oa.
func (d *Deliver) OA() Address {
	return d.oa
}

// MaxUDL returns the maximum number of octets that can be encoded into the UD.
// Note that for 7bit encoding this can result in up to 160 septets.
func (d *Deliver) MaxUDL() int {
	return 140
}

// SCTS returns the Deliver scts.
func (d *Deliver) SCTS() Timestamp {
	return d.scts
}

// SetOA sets the Deliver oa field.
func (d *Deliver) SetOA(oa Address) {
	d.oa = oa
}

// SetSCTS sets the Deliver scts field.
func (d *Deliver) SetSCTS(scts Timestamp) {
	d.scts = scts
}

// MarshalBinary marshals a SMS-Deliver PDU into the corresponding byte array.
func (d *Deliver) MarshalBinary() ([]byte, error) {
	b := []byte{d.firstOctet}
	oa, err := d.oa.MarshalBinary()
	if err != nil {
		return nil, EncodeError("oa", err)
	}
	b = append(b, oa...)
	b = append(b, d.pid, d.dcs)
	scts, err := d.scts.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	ud, err := d.encodeUserData()
	if err != nil {
		return nil, EncodeError("ud", err)
	}
	b = append(b, ud...)
	return b, nil
}

// UnmarshalBinary unmarshals a SMS-Deliver PDU from the corresponding byte array.
// In the case of error the Deliver will be partially unmarshalled, up to
// the point that the decoding error was detected.
func (d *Deliver) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	d.firstOctet = src[0]
	ri := 1
	n, err := d.oa.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("oa", ri, err)
	}
	ri += n
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	d.pid = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("dcs", ri, ErrUnderflow)
	}
	d.dcs = src[ri]
	ri++
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err = d.scts.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	err = d.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func decodeDeliver(src []byte) (TPDU, error) {
	d := NewDeliver()
	if err := d.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return d, nil
}

// RegisterDeliverDecoder registers a decoder for the Deliver TPDU.
func RegisterDeliverDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtDeliver, MT, decodeDeliver)
}

// RegisterReservedDecoder registers a decoder for the Deliver TPDU for the Reserved message type.
func RegisterReservedDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtReserved, MT, decodeDeliver)
}
