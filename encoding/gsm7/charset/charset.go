// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

// Package charset provides encoders and decoders for GSM character sets.
package charset

// DefaultDecoder returns the default mapping table from GSM7 to UTF8.
func DefaultDecoder() Decoder {
	if defaultDecoder == nil {
		defaultDecoder = generateDefaultDecoder()
	}
	return defaultDecoder
}

// NewDecoder returns the mapping table from GSM7 to UTF8 for the given language.
func NewDecoder(nli int) Decoder {
	if di, ok := decoder[nli]; ok {
		if *di.e == nil {
			*di.e = di.g()
		}
		return *di.e
	}
	return DefaultDecoder()
}

// DefaultExtDecoder returns the default extension mapping table from GSM7 to UTF8.
func DefaultExtDecoder() Decoder {
	return defaultExtDecoder
}

// NewExtDecoder returns the extension mapping table from GSM7 to UTF8 for the given language.
func NewExtDecoder(nli int) Decoder {
	if di, ok := extDecoder[nli]; ok {
		return *di.e
	}
	return DefaultExtDecoder()
}

// DefaultEncoder returns the default mapping table from UTF8 to GSM7.
func DefaultEncoder() Encoder {
	if defaultEncoder == nil {
		defaultEncoder = generateDefaultEncoder()
	}
	return defaultEncoder
}

// NewEncoder returns the mapping table from UTF8 to GSM7 for the given language.
func NewEncoder(nli int) Encoder {
	if ei, ok := encoder[nli]; ok {
		if *ei.e == nil {
			*ei.e = ei.g()
		}
		return *ei.e
	}
	return DefaultEncoder()
}

// DefaultExtEncoder returns the default extension mapping table from UTF8 to GSM7.
func DefaultExtEncoder() Encoder {
	if defaultExtEncoder == nil {
		defaultExtEncoder = generateDefaultExtEncoder()
	}
	return defaultExtEncoder
}

// NewExtEncoder returns the extension mapping table from UTF8 to GSM7 for the given language.
func NewExtEncoder(nli int) Encoder {
	if ei, ok := extEncoder[nli]; ok {
		if *ei.e == nil {
			*ei.e = ei.g()
		}
		return *ei.e
	}
	return DefaultExtEncoder()
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
	Default int = iota
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

	// Helper consts for creating and iterating over slices

	// End marker for loops (exclusive)
	End
	// Start point for loops (inclusive)
	Start = Turkish
	// Size is for array sizing (excluding default)
	Size = End - Start
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

type decoderGenerator func() Decoder

type decoderNGen struct {
	e *Decoder
	g decoderGenerator
}

type encoderGenerator func() Encoder

type encoderNGen struct {
	e *Encoder
	g encoderGenerator
}

var (
	decoder = map[int]decoderNGen{
		Turkish: {&turkishDecoder, generateTurkishDecoder},
		// Spanish uses default
		Portuguese: {&portugueseDecoder, generatePortugueseDecoder},
		Bengali:    {&bengaliDecoder, nil},
		Gujaranti:  {&gujaratiDecoder, nil},
		Hindi:      {&hindiDecoder, nil},
		Kannada:    {&kannadaDecoder, nil},
		Malayalam:  {&malayalamDecoder, nil},
		Oriya:      {&oriyaDecoder, nil},
		Punjabi:    {&punjabiDecoder, nil},
		Tamil:      {&tamilDecoder, nil},
		Telugu:     {&teluguDecoder, nil},
		Urdu:       {&urduDecoder, nil},
	}
	extDecoder = map[int]decoderNGen{
		Turkish:    {&turkishExtDecoder, nil},
		Spanish:    {&spanishExtDecoder, nil},
		Portuguese: {&portugueseExtDecoder, nil},
		Bengali:    {&bengaliExtDecoder, nil},
		Gujaranti:  {&gujaratiExtDecoder, nil},
		Hindi:      {&hindiExtDecoder, nil},
		Kannada:    {&kannadaExtDecoder, nil},
		Malayalam:  {&malayalamExtDecoder, nil},
		Oriya:      {&oriyaExtDecoder, nil},
		Punjabi:    {&punjabiExtDecoder, nil},
		Tamil:      {&tamilExtDecoder, nil},
		Telugu:     {&teluguExtDecoder, nil},
		Urdu:       {&urduExtDecoder, nil},
	}
	encoder = map[int]encoderNGen{
		Turkish: {&turkishEncoder, generateTurkishEncoder},
		// Spanish uses default
		Portuguese: {&portugueseEncoder, generatePortugueseEncoder},
		Bengali:    {&bengaliEncoder, generateBengaliEncoder},
		Gujaranti:  {&gujaratiEncoder, generateGujaratiEncoder},
		Hindi:      {&hindiEncoder, generateHindiEncoder},
		Kannada:    {&kannadaEncoder, generateKannadaEncoder},
		Malayalam:  {&malayalamEncoder, generateMalayalamEncoder},
		Oriya:      {&oriyaEncoder, generateOriyaEncoder},
		Punjabi:    {&punjabiEncoder, generatePunjabiEncoder},
		Tamil:      {&tamilEncoder, generateTamilEncoder},
		Telugu:     {&teluguEncoder, generateTeluguEncoder},
		Urdu:       {&urduEncoder, generateUrduEncoder},
	}
	extEncoder = map[int]encoderNGen{
		Turkish:    {&turkishExtEncoder, generateTurkishExtEncoder},
		Spanish:    {&spanishExtEncoder, generateSpanishExtEncoder},
		Portuguese: {&portugueseExtEncoder, generatePortugueseExtEncoder},
		Bengali:    {&bengaliExtEncoder, generateBengaliExtEncoder},
		Gujaranti:  {&gujaratiExtEncoder, generateGujaratiExtEncoder},
		Hindi:      {&hindiExtEncoder, generateHindiExtEncoder},
		Kannada:    {&kannadaExtEncoder, generateKannadaExtEncoder},
		Malayalam:  {&malayalamExtEncoder, generateMalayalamExtEncoder},
		Oriya:      {&oriyaExtEncoder, generateOriyaExtEncoder},
		Punjabi:    {&punjabiExtEncoder, generatePunjabiExtEncoder},
		Tamil:      {&tamilExtEncoder, generateTamilExtEncoder},
		Telugu:     {&teluguExtEncoder, generateTeluguExtEncoder},
		Urdu:       {&urduExtEncoder, generateUrduExtEncoder},
	}
)
