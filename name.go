package df2014

import "strconv"

type NameForm uint16

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
	NameNounSingular:      "singular noun",
	NameNounPlural:        "plural noun",
	NameAdjective:         "adjective",
	NamePrefix:            "prefix",
	NamePresent1st:        "present first person",
	NamePresent3rd:        "present third person",
	NamePreterite:         "preterite",
	NamePastParticiple:    "past participle",
	NamePresentParticiple: "present participle",
}

func (i NameForm) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if int(i) < len(nameForms) {
		buf = append(buf, " ("...)
		buf = append(buf, nameForms[i]...)
		buf = append(buf, ')')
	}

	return buf
}

type WordIndex int32

func (i WordIndex) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendInt(buf, int64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if i >= 0 && int(i) < len(w.StringTables.Word) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Word[i]...)
		buf = append(buf, ')')
	}

	return buf
}

type TranslationIndex int32

func (i TranslationIndex) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendInt(buf, int64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if i >= 0 && int(i) < len(w.StringTables.Translation) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Translation[i]...)
		buf = append(buf, ')')
	}

	return buf
}

type Name struct {
	First    string
	Nick     string
	Index    NameIndices
	Form     NameForms
	Language TranslationIndex `df2014_assert_gte:"0"`
	Unknown  int16
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
