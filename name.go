package df2014

import (
	"reflect"
)

type NameForm uint16
type NameForm23a int32

const (
	NameNounSingular NameForm = iota
	NameNounPlural
	NameAdjective
	NamePrefix
	NamePresent1st
	NamePresent3rd
	NamePreterite
	NamePastParticiple
	NamePresentParticiple
)

var nameForms = [...]string{
	"singular noun",
	"plural noun",
	"adjective",
	"prefix",
	"present first person",
	"present third person",
	"preterite",
	"past participle",
	"present participle",
}

func (i NameForm) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), nameForms[:], buf)
}

func (i NameForm23a) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), nameForms[:], buf)
}

type Name struct {
	First    string
	Nick     string           `df2014_version_min:"1116"`
	Index    NameIndices      `df2014_version_min:"1205"`
	Form     NameForms        `df2014_version_min:"1205"`
	Parts23a NameParts23a     `df2014_version_max:"1169"`
	Language TranslationIndex `df2014_assert_gte:"-1"`
	Mode     int16
}

type NameIndices struct {
	FrontCompound  WordIndex `df2014_assert_gte:"-1"`
	RearCompound   WordIndex `df2014_assert_gte:"-1"`
	Adjective1     WordIndex `df2014_assert_gte:"-1"`
	Adjective2     WordIndex `df2014_assert_gte:"-1"`
	HyphenCompound WordIndex `df2014_assert_gte:"-1"`
	TheX           WordIndex `df2014_assert_gte:"-1"`
	OfX            WordIndex `df2014_assert_gte:"-1"`
}

type NameForms struct {
	FrontCompound  NameForm
	RearCompound   NameForm
	Adjective1     NameForm
	Adjective2     NameForm
	HyphenCompound NameForm
	TheX           NameForm
	OfX            NameForm
}

type NamePart23a struct {
	Index WordIndex `df2014_assert_gte:"-1"`
	Form  NameForm23a
}

type NameParts23a struct {
	FrontCompound  NamePart23a
	RearCompound   NamePart23a
	Adjective1     NamePart23a
	Adjective2     NamePart23a
	HyphenCompound NamePart23a
	TheX           NamePart23a
	OfX            NamePart23a
}
