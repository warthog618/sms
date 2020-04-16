// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

// smsdecode provides an example of unmarshalling and displaying a SMS TPDU.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/pdumode"
	"github.com/warthog618/sms/encoding/tpdu"
)

func main() {
	pm := flag.Bool("p", false, "PDU is prefixed with SCA (PDU mode)")
	orig := flag.Bool("o", false, "PDU is mobile originated")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	tp, smsc, err := decode(flag.Arg(0), *pm, *orig)
	if err != nil {
		log.Fatal(err)
	}
	if smsc != nil {
		dumpSMSC(os.Stdout, smsc)
	}
	dumpTPDU(os.Stdout, tp)
}

func decode(s string, pm, mo bool) (p *tpdu.TPDU, a *pdumode.SMSCAddress, err error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return
	}
	if pm {
		var pdu *pdumode.PDU
		pdu, err = pdumode.UnmarshalHexString(s)
		if err != nil {
			return
		}
		a = &pdu.SMSC
		b = pdu.TPDU
	}
	if mo {
		p, err = sms.Unmarshal(b, sms.AsMO)
		return
	}
	p, err = sms.Unmarshal(b)
	return
}

func dumpSMSC(w io.Writer, smsc *pdumode.SMSCAddress) {
	n := smsc.Number()
	fmt.Fprintf(w, "SMSC: %s\n", n)
}

func dumpTPDU(w io.Writer, t *tpdu.TPDU) {
	var st string
	var dump func(w io.Writer, t *tpdu.TPDU)
	switch t.SmsType() {
	case tpdu.SmsCommand:
		st = "SMS-COMMAND"
		dump = dumpCommand
	case tpdu.SmsDeliver:
		st = "SMS-DELIVER"
		dump = dumpDeliver
	case tpdu.SmsDeliverReport:
		st = "SMS-DELIVER-REPORT"
		dump = dumpDeliverReport
	case tpdu.SmsStatusReport:
		st = "SMS-STATUS-REPORT"
		dump = dumpStatusReport
	case tpdu.SmsSubmit:
		st = "SMS-SUBMIT"
		dump = dumpSubmit
	case tpdu.SmsSubmitReport:
		st = "SMS-SUBMIT-REPORT"
		dump = dumpSubmitReport
	}
	fmt.Fprintf(w, "TPDU: %s\n", st)
	dump(w, t)
}

func dumpCommand(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-SRR: %t\n", t.FirstOctet.SRR())
	fmt.Fprintf(w, "TP-MR: %d\n", t.MR)
	fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	fmt.Fprintf(w, "TP-CT: 0x%02x\n", t.CT)
	fmt.Fprintf(w, "TP-MN: %d\n", t.MN)
	fmt.Fprintf(w, "TP-DA: %s\n", t.DA.Number())
	fmt.Fprintf(w, "TP-SCTS: %s\n", t.SCTS)
	fmt.Fprintf(w, "TP-CDL: %d\n", len(t.UD))
	dumpCD(w, t.UD)
}

func dumpDeliver(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-MMS: %t\n", t.FirstOctet.MMS())
	fmt.Fprintf(w, "TP-LP: %t\n", t.FirstOctet.LP())
	fmt.Fprintf(w, "TP-RP: %t\n", t.FirstOctet.RP())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-SRI: %t\n", t.FirstOctet.SRI())
	fmt.Fprintf(w, "TP-OA: %s\n", t.OA.Number())
	fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	fmt.Fprintf(w, "TP-DCS: %s\n", t.DCS)
	fmt.Fprintf(w, "TP-SCTS: %s\n", t.SCTS)
	if t.UDH != nil {
		dumpUDH(w, t.UDH)
	}
	dumpUD(w, t.UD)
}

func dumpDeliverReport(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-FCS: 0x%02x\n", t.FCS)
	fmt.Fprintf(w, "TP-PI: %s\n", t.PI)
	if t.PI.PID() {
		fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	}
	if t.PI.DCS() {
		fmt.Fprintf(w, "TP-DCS: %s\n", t.DCS)
	}
	if t.UDH != nil {
		dumpUDH(w, t.UDH)
	}
	dumpUD(w, t.UD)
}

func dumpStatusReport(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-MMS: %t\n", t.FirstOctet.MMS())
	fmt.Fprintf(w, "TP-LP: %t\n", t.FirstOctet.LP())
	fmt.Fprintf(w, "TP-SRQ: %t\n", t.FirstOctet.SRQ())
	fmt.Fprintf(w, "TP-MR: %d\n", t.MR)
	fmt.Fprintf(w, "TP-RA: %s\n", t.RA.Number())
	fmt.Fprintf(w, "TP-SCTS: %s\n", t.SCTS)
	fmt.Fprintf(w, "TP-DT: %s\n", t.DT)
	fmt.Fprintf(w, "TP-ST: 0x%02x\n", t.ST)
	fmt.Fprintf(w, "TP-PI: %s\n", t.PI)
	if t.PI.PID() {
		fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	}
	if t.PI.DCS() {
		fmt.Fprintf(w, "TP-DCS: %s\n", t.DCS)
	}
	if t.UDH != nil {
		dumpUDH(w, t.UDH)
	}
	dumpUD(w, t.UD)
}

func dumpSubmit(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-RD: %t\n", t.FirstOctet.RD())
	fmt.Fprintf(w, "TP-VPF: 0x%02x %s\n", int(t.FirstOctet.VPF()), t.FirstOctet.VPF())
	fmt.Fprintf(w, "TP-RP: %t\n", t.FirstOctet.RP())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-SRR: %t\n", t.FirstOctet.SRR())
	fmt.Fprintf(w, "TP-MR: %d\n", t.MR)
	fmt.Fprintf(w, "TP-DA: %s\n", t.DA.Number())
	fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	fmt.Fprintf(w, "TP-DCS: %s\n", t.DCS)
	dumpVP(w, t.VP)
	if t.UDH != nil {
		dumpUDH(w, t.UDH)
	}
	dumpUD(w, t.UD)
}

func dumpSubmitReport(w io.Writer, t *tpdu.TPDU) {
	fmt.Fprintf(w, "TP-MTI: 0x%02x %s\n", int(t.SmsType().MTI()), t.SmsType().MTI())
	fmt.Fprintf(w, "TP-UDHI: %t\n", t.FirstOctet.UDHI())
	fmt.Fprintf(w, "TP-FCS: 0x%02x\n", t.FCS)
	fmt.Fprintf(w, "TP-PI: %s\n", t.PI)
	fmt.Fprintf(w, "TP-SCTS: %s\n", t.SCTS)
	if t.PI.PID() {
		fmt.Fprintf(w, "TP-PID: 0x%02x\n", t.PID)
	}
	if t.PI.DCS() {
		fmt.Fprintf(w, "TP-DCS: %s\n", t.DCS)
	}
	if t.UDH != nil {
		dumpUDH(w, t.UDH)
	}
	dumpUD(w, t.UD)
}

func dumpCD(w io.Writer, ud []byte) {
	lines := strings.Split(hex.Dump(ud), "\n")
	fmt.Fprintf(w, "TP-CD: %s\n", lines[0])
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "       %s\n", l)
	}
}

func dumpVP(w io.Writer, vp tpdu.ValidityPeriod) {
	switch vp.Format {
	case tpdu.VpfNotPresent:
		fmt.Fprintf(w, "TP-VP: Not Present\n")
	case tpdu.VpfAbsolute:
		fmt.Fprintf(w, "TP-VP: Absolute - %s\n", vp.Time)
	case tpdu.VpfEnhanced:
		fmt.Fprintf(w, "TP-VP: Enhanced %s - %s\n",
			tpdu.EnhancedFormat(vp.EFI), vp.Duration)
	case tpdu.VpfRelative:
		fmt.Fprintf(w, "TP-VP: Relative - %s\n", vp.Duration)
	}
}

func dumpUDH(w io.Writer, udh tpdu.UserDataHeader) {
	ie := udh[0]
	fmt.Fprintf(w, "TP-UDH: ID: %d  Data: %v\n", ie.ID, ie.Data)
	for _, ie = range udh[1:] {
		fmt.Fprintf(w, "       ID: %d  Data: %v\n", ie.ID, ie.Data)
	}
}

func dumpUD(w io.Writer, ud []byte) {
	lines := strings.Split(strings.TrimSpace(hex.Dump(ud)), "\n")
	fmt.Fprintf(w, "TP-UD: %s\n", lines[0])
	for _, l := range lines[1:] {
		fmt.Fprintf(w, "       %s\n", l)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: smsdecode [-p] [-o] <sms>\n")
	flag.PrintDefaults()
}
