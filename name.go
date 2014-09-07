package df2014

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

func (i NameForm) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	return prettyPrintIndex(int64(i), nameForms[:], buf)
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
