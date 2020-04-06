// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/ms/pdumode"
)

func main() {
	pm := flag.Bool("p", false, "PDU is prefixed with SCA (PDU mode)")
	orig := flag.Bool("o", false, "PDU is mobile originated")
	var dirOpt sms.DecodeOption
	flag.Usage = usage
	flag.Parse()
	if *orig {
		dirOpt = sms.AsMO
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
		smsc, ntb, err := pdumode.Decode(b)
		if err != nil {
			log.Fatal(err)
		}
		tb = ntb
		spew.Dump(smsc)
	}
	tp, err := sms.Decode(tb, dirOpt)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(tp)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: smsdecode [-p] [-o] <sms>\n")
	flag.PrintDefaults()
}
