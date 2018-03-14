// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/gsm7/charset/basic"
	"github.com/warthog618/sms/encoding/gsm7/charset/bengali"
	"github.com/warthog618/sms/encoding/gsm7/charset/gujarati"
	"github.com/warthog618/sms/encoding/gsm7/charset/hindi"
	"github.com/warthog618/sms/encoding/gsm7/charset/kannada"
	"github.com/warthog618/sms/encoding/gsm7/charset/malayalam"
	"github.com/warthog618/sms/encoding/gsm7/charset/oriya"
	"github.com/warthog618/sms/encoding/gsm7/charset/portuguese"
	"github.com/warthog618/sms/encoding/gsm7/charset/punjabi"
	"github.com/warthog618/sms/encoding/gsm7/charset/spanish"
	"github.com/warthog618/sms/encoding/gsm7/charset/tamil"
	"github.com/warthog618/sms/encoding/gsm7/charset/telugu"
	"github.com/warthog618/sms/encoding/gsm7/charset/turkish"
	"github.com/warthog618/sms/encoding/gsm7/charset/urdu"
)

func main() {
	fmt.Println("Basic (default)")
	charset.Display(basic.NewDecoder())
	fmt.Println()
	fmt.Println("Basic (default) Extensions")
	charset.Display(basic.NewExtDecoder())
	fmt.Println()

	fmt.Println("Turkish")
	charset.Display(turkish.NewDecoder())
	fmt.Println()
	fmt.Println("Turkish Extensions")
	charset.Display(turkish.NewExtDecoder())
	fmt.Println()

	fmt.Println("Spanish")
	charset.Display(basic.NewDecoder())
	fmt.Println()
	fmt.Println("Spanish Extensions")
	charset.Display(spanish.NewExtDecoder())
	fmt.Println()

	fmt.Println("Portuguese")
	charset.Display(portuguese.NewDecoder())
	fmt.Println()
	fmt.Println("Portuguese Extensions")
	charset.Display(portuguese.NewExtDecoder())
	fmt.Println()

	fmt.Println("Bengali")
	charset.Display(bengali.NewDecoder())
	fmt.Println()
	fmt.Println("Bengali Extensions")
	charset.Display(bengali.NewExtDecoder())
	fmt.Println()

	fmt.Println("Gujarati")
	charset.Display(gujarati.NewDecoder())
	fmt.Println()
	fmt.Println("Gujarati Extensions")
	charset.Display(gujarati.NewExtDecoder())
	fmt.Println()

	fmt.Println("Hindi")
	charset.Display(hindi.NewDecoder())
	fmt.Println()
	fmt.Println("Hindi Extensions")
	charset.Display(hindi.NewExtDecoder())
	fmt.Println()

	fmt.Println("Kannada")
	charset.Display(kannada.NewDecoder())
	fmt.Println()
	fmt.Println("Kannada Extensions")
	charset.Display(kannada.NewExtDecoder())
	fmt.Println()

	fmt.Println("Malayalam")
	charset.Display(malayalam.NewDecoder())
	fmt.Println()
	fmt.Println("Malayalam Extensions")
	charset.Display(malayalam.NewExtDecoder())
	fmt.Println()

	fmt.Println("Oriya")
	charset.Display(oriya.NewDecoder())
	fmt.Println()
	fmt.Println("Oriya Extensions")
	charset.Display(oriya.NewExtDecoder())
	fmt.Println()

	fmt.Println("Punjabi")
	charset.Display(punjabi.NewDecoder())
	fmt.Println()
	fmt.Println("Punjabi Extensions")
	charset.Display(punjabi.NewExtDecoder())
	fmt.Println()

	fmt.Println("Tamil")
	charset.Display(tamil.NewDecoder())
	fmt.Println()
	fmt.Println("Tamil Extensions")
	charset.Display(tamil.NewExtDecoder())
	fmt.Println()

	fmt.Println("Telugu")
	charset.Display(telugu.NewDecoder())
	fmt.Println()
	fmt.Println("telugu Extensions")
	charset.Display(telugu.NewExtDecoder())
	fmt.Println()

	fmt.Println("Urdu")
	charset.Display(urdu.NewDecoder())
	fmt.Println()
	fmt.Println("Urdu Extensions")
	charset.Display(urdu.NewExtDecoder())
	fmt.Println()
}
