// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/warthog618/sms/encoding/gsm7/charset"
)

var charsetName = []string{
	"Basic (default)",
	"Turkish",
	"Spanish",
	"Portuguese",
	"Bengali",
	"Gujaranti",
	"Hindi",
	"Kannada",
	"Malayalam",
	"Oriya",
	"Punjabi",
	"Tamil",
	"Telugu",
	"Urdu",
}

func main() {
	for nli := charset.Default; nli <= charset.Urdu; nli++ {
		fmt.Printf("%s Locking (NLI=%d)\n", charsetName[nli], nli)
		Display(charset.NewDecoder(nli))
		fmt.Println()
		fmt.Printf("%s Shift (NLI=%d)\n", charsetName[nli], nli)
		Display(charset.NewExtDecoder(nli))
		fmt.Println()
	}
}

// Display prints the character set for a given character set decoder.
func Display(m charset.Decoder) {
	specials := map[rune]string{
		'\n':   "LF",
		'\r':   "CR",
		'\f':   "FF",
		' ':    "SP",
		0x1b:   "ESC",
		0x20ac: " €",
	}
	fmt.Printf("      ")
	for c := 0; c < 8; c++ {
		fmt.Printf("0x%d_ ", c)
	}
	fmt.Println("")
	for r := 0; r < 0x10; r++ {
		fmt.Printf("0x_%x: ", r)
		for c := 0; c < 8; c++ {
			k := byte(c*0x10 + r)
			if v, ok := m[k]; ok {
				if s, ok := specials[v]; ok {
					fmt.Printf("%3s  ", s)
				} else if v >= 0x400 {
					fmt.Printf("%04x ", v)
				} else {
					fmt.Printf("  %c  ", v)
				}
			} else {
				fmt.Printf("     ")
			}
		}
		fmt.Println()
	}
}
