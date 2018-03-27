// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package charset

// DefaultDecoder returns the default mapping table from GSM7 to UTF8.
func DefaultDecoder() Decoder {
	if defaultDecoder == nil {
		defaultDecoder = generateDefaultDecoder()
	}
	return defaultDecoder
}

// NewDecoder returns the mapping table from GSM7 to UTF8 for the given language.
func NewDecoder(nli NationalLanguageIdentifier) Decoder {
	switch nli {
	case Turkish:
		if turkishDecoder == nil {
			turkishDecoder = generateTurkishDecoder()
		}
		return turkishDecoder
	case Portuguese:
		if portugueseDecoder == nil {
			portugueseDecoder = generatePortugueseDecoder()
		}
		return portugueseDecoder
	case Bengali:
		return bengaliDecoder
	case Gujaranti:
		return gujaratiDecoder
	case Hindi:
		return hindiDecoder
	case Kannada:
		return kannadaDecoder
	case Malayalam:
		return malayalamDecoder
	case Oriya:
		return oriyaDecoder
	case Punjabi:
		return punjabiDecoder
	case Tamil:
		return tamilDecoder
	case Telugu:
		return teluguDecoder
	case Urdu:
		return urduDecoder
	case Spanish:
		fallthrough
	default:
		return DefaultDecoder()
	}
}

// DefaultExtDecoder returns the default extension mapping table from GSM7 to UTF8.
func DefaultExtDecoder() Decoder {
	return defaultExtDecoder
}

// NewExtDecoder returns the extension mapping table from GSM7 to UTF8 for the given language.
func NewExtDecoder(nli NationalLanguageIdentifier) Decoder {
	switch nli {
	case Turkish:
		return turkishExtDecoder
	case Spanish:
		return spanishExtDecoder
	case Portuguese:
		return portugueseExtDecoder
	case Bengali:
		return bengaliExtDecoder
	case Gujaranti:
		return gujaratiExtDecoder
	case Hindi:
		return hindiExtDecoder
	case Kannada:
		return kannadaExtDecoder
	case Malayalam:
		return malayalamExtDecoder
	case Oriya:
		return oriyaExtDecoder
	case Punjabi:
		return punjabiExtDecoder
	case Tamil:
		return tamilExtDecoder
	case Telugu:
		return teluguExtDecoder
	case Urdu:
		return urduExtDecoder
	default:
		return defaultExtDecoder
	}
}

// DefaultEncoder returns the default mapping table from UTF8 to GSM7.
func DefaultEncoder() Encoder {
	if defaultEncoder == nil {
		defaultEncoder = generateDefaultEncoder()
	}
	return defaultEncoder
}

// NewEncoder returns the mapping table from UTF8 to GSM7 for the given language.
func NewEncoder(nli NationalLanguageIdentifier) Encoder {
	switch nli {
	case Turkish:
		if turkishEncoder == nil {
			turkishEncoder = generateTurkishEncoder()
		}
		return turkishEncoder
	case Portuguese:
		if portugueseEncoder == nil {
			portugueseEncoder = generatePortugueseEncoder()
		}
		return portugueseEncoder
	case Bengali:
		if bengaliEncoder == nil {
			bengaliEncoder = generateBengaliEncoder()
		}
		return bengaliEncoder
	case Gujaranti:
		if gujaratiEncoder == nil {
			gujaratiEncoder = generateGujaratiEncoder()
		}
		return gujaratiEncoder
	case Hindi:
		if hindiEncoder == nil {
			hindiEncoder = generateHindiEncoder()
		}
		return hindiEncoder
	case Kannada:
		if kannadaEncoder == nil {
			kannadaEncoder = generateKannadaEncoder()
		}
		return kannadaEncoder
	case Malayalam:
		if malayalamEncoder == nil {
			malayalamEncoder = generateMalayalamEncoder()
		}
		return malayalamEncoder
	case Oriya:
		if oriyaEncoder == nil {
			oriyaEncoder = generateOriyaEncoder()
		}
		return oriyaEncoder
	case Punjabi:
		if punjabiEncoder == nil {
			punjabiEncoder = generatePunjabiEncoder()
		}
		return punjabiEncoder
	case Tamil:
		if tamilEncoder == nil {
			tamilEncoder = generateTamilEncoder()
		}
		return tamilEncoder
	case Telugu:
		if teluguEncoder == nil {
			teluguEncoder = generateTeluguEncoder()
		}
		return teluguEncoder
	case Urdu:
		if urduEncoder == nil {
			urduEncoder = generateUrduEncoder()
		}
		return urduEncoder
	case Spanish:
		fallthrough
	default:
		return DefaultEncoder()
	}
}

// DefaultExtEncoder returns the default extension mapping table from UTF8 to GSM7.
func DefaultExtEncoder() Encoder {
	if defaultExtEncoder == nil {
		defaultExtEncoder = generateDefaultExtEncoder()
	}
	return defaultExtEncoder
}

// NewExtEncoder returns the extension mapping table from UTF8 to GSM7 for the given language.
func NewExtEncoder(nli NationalLanguageIdentifier) Encoder {
	switch nli {
	case Turkish:
		if turkishExtEncoder == nil {
			turkishExtEncoder = generateTurkishExtEncoder()
		}
		return turkishExtEncoder
	case Spanish:
		if spanishExtEncoder == nil {
			spanishExtEncoder = generateSpanishExtEncoder()
		}
		return spanishExtEncoder
	case Portuguese:
		if portugueseExtEncoder == nil {
			portugueseExtEncoder = generatePortugueseExtEncoder()
		}
		return portugueseExtEncoder
	case Bengali:
		if bengaliExtEncoder == nil {
			bengaliExtEncoder = generateBengaliExtEncoder()
		}
		return bengaliExtEncoder
	case Gujaranti:
		if gujaratiExtEncoder == nil {
			gujaratiExtEncoder = generateGujaratiExtEncoder()
		}
		return gujaratiExtEncoder
	case Hindi:
		if hindiExtEncoder == nil {
			hindiExtEncoder = generateHindiExtEncoder()
		}
		return hindiExtEncoder
	case Kannada:
		if kannadaExtEncoder == nil {
			kannadaExtEncoder = generateKannadaExtEncoder()
		}
		return kannadaExtEncoder
	case Malayalam:
		if malayalamExtEncoder == nil {
			malayalamExtEncoder = generateMalayalamExtEncoder()
		}
		return malayalamExtEncoder
	case Oriya:
		if oriyaExtEncoder == nil {
			oriyaExtEncoder = generateOriyaExtEncoder()
		}
		return oriyaExtEncoder
	case Punjabi:
		if punjabiExtEncoder == nil {
			punjabiExtEncoder = generatePunjabiExtEncoder()
		}
		return punjabiExtEncoder
	case Tamil:
		if tamilExtEncoder == nil {
			tamilExtEncoder = generateTamilExtEncoder()
		}
		return tamilExtEncoder
	case Telugu:
		if teluguExtEncoder == nil {
			teluguExtEncoder = generateTeluguExtEncoder()
		}
		return teluguExtEncoder
	case Urdu:
		if urduExtEncoder == nil {
			urduExtEncoder = generateUrduExtEncoder()
		}
		return urduExtEncoder
	default:
		return DefaultExtEncoder()
	}
}

// Decoder provides a mapping from GSM7 byte to UTF8 rune.
type Decoder map[byte]rune

// Encoder provides a mapping from UTF8 rune to GSM7 byte.
type Encoder map[rune]byte

// NationalLanguageIdentifier indicates the character set in use, as defined in
// 3GPP TS 23.038 Section 6.2.1.2.4.
type NationalLanguageIdentifier int

const (
	// Default character set.
	Default NationalLanguageIdentifier = iota
	// Turkish character set.
	Turkish
	// Spanish character set
	Spanish
	// Portuguese character set
	Portuguese
	// Bengali character set
	Bengali
	// Gujaranti character set
	Gujaranti
	// Hindi character set
	Hindi
	// Kannada character set
	Kannada
	// Malayalam character set
	Malayalam
	// Oriya character set
	Oriya
	// Punjabi character set
	Punjabi
	// Tamil character set
	Tamil
	// Telugu character set
	Telugu
	// Urdu character set
	Urdu
)

func generateEncoder(d Decoder) Encoder {
	e := make(Encoder, len(d))
	for k, v := range d {
		if ko, ok := e[v]; !ok || ko > k {
			e[v] = k
		}
	}
	return e
}

func generateEncoderFromRunes(runes []rune) Encoder {
	e := make(Encoder, len(runes))
	for i, r := range runes {
		e[r] = byte(i)
	}
	return e
}

func generateDecoderFromRunes(runes []rune) Decoder {
	dset := make(Decoder, len(runes))
	for i, r := range runes {
		dset[byte(i)] = r
	}
	return dset
}
