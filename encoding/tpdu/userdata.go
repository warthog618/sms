// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// UserData represents the User Data field as defined in 3GPP TS 23.040 Section 9.2.3.24.
// The UserData is comprised of an optional User Data Header and a short message field.
type UserData []byte

// UserDataHeader represents the header section of the User Data as defined in
// 3GPP TS 23.040 Section 9.2.3.24.
type UserDataHeader []InformationElement

// InformationElement represents one of the information elements contained in
// the User Data Header.
type InformationElement struct {
	ID   byte
	Data []byte
}

// MarshalBinary marshals the User Data Header, including the UDHL, into binary.
func (udh UserDataHeader) MarshalBinary() ([]byte, error) {
	if len(udh) == 0 {
		return nil, nil
	}
	udhl := 0
	for _, ie := range udh {
		udhl += (2 + len(ie.Data))
	}
	b := make([]byte, 0, udhl+1)
	b = append(b, byte(udhl))
	for _, ie := range udh {
		b = append(b, ie.ID, byte(len(ie.Data)))
		b = append(b, ie.Data...)
	}
	return b, nil
}

// UnmarshalBinary reads the InformationElements from the binary User Data Haeder.
// The src contains the complete UDH, including the UDHL and all IEs.
// The function returns the number of bytes read from src, and any error detected
// while unmarshalling.
func (udh *UserDataHeader) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 1 {
		return 0, DecodeError("udhl", 0, ErrUnderflow)
	}
	udhl := int(src[0])
	udhl++ // so it includes itself
	ri := 1
	if len(src) < udhl {
		return ri, DecodeError("ie", ri, ErrUnderflow)
	}
	ies := []InformationElement(nil)
	for ri < udhl {
		if udhl < ri+2 {
			return ri, DecodeError("ie", ri, ErrUnderflow)
		}
		var ie InformationElement
		ie.ID = src[ri]
		ri++
		iedl := int(src[ri])
		ri++
		if len(src) < ri+iedl {
			return ri, DecodeError("ied", ri, ErrUnderflow)
		}
		ie.Data = append([]byte(nil), src[ri:ri+iedl]...)
		ri += iedl
		ies = append(ies, ie)
	}
	*udh = ies
	return udhl, nil
}

// IE returns the last instance of the GetIE with the given id in the UDH.
// If no such GetIE is found then the function returns false.
func (udh UserDataHeader) IE(id byte) (InformationElement, bool) {
	for i := len(udh) - 1; i >= 0; i-- {
		if udh[i].ID == id {
			return udh[i], true
		}
	}
	return InformationElement{}, false
}

// IEs returns all instances of the GetIEs with the given id in the UDH.
func (udh UserDataHeader) IEs(id byte) []InformationElement {
	ies := []InformationElement(nil)
	for _, ie := range udh {
		if ie.ID == id {
			ies = append(ies, ie)
		}
	}
	return ies
}
