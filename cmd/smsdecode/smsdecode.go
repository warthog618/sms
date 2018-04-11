// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/pdumode"
)

func main() {
	pm := flag.Bool("p", false, "PDU is prefixed with SCA (PDU mode)")
	orig := flag.Bool("o", false, "PDU is mobile originated")
	drn := tpdu.MT
	flag.Usage = usage
	flag.Parse()
	if *orig {
		drn = tpdu.MO
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	b, err := hex.DecodeString(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	tb := b
	if *pm {
		pd := pdumode.Decoder{}
		smsc, ntb, err := pd.Decode(b)
		if err != nil {
			log.Fatal(err)
		}
		tb = ntb
		spew.Dump(smsc)
	}
	td, err := tpdu.NewDecoder(
		tpdu.RegisterCommandDecoder,
		tpdu.RegisterDeliverDecoder,
		tpdu.RegisterDeliverReportDecoder,
		tpdu.RegisterReservedDecoder,
		tpdu.RegisterSubmitDecoder,
		tpdu.RegisterSubmitReportDecoder,
		tpdu.RegisterStatusReportDecoder,
	)
	if err != nil {
		log.Fatal(err)
	}
	tp, err := td.Decode(tb, drn)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(tp)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: smsdecode [-p] [-o] <sms>\n")
	flag.PrintDefaults()
}
