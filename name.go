package df2014

// indices into Name.Index and Name.Form
const (
	NameFrontCompound = iota
	NameRearCompound
	NameAdjective1
	NameAdjective2
	NameHyphenCompound
	NameTheX
	NameOfX
)

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

type Name struct {
	First    string
	Nick     string
	Index    [7]int32
	Form     [7]NameForm
	Language uint32
	Unknown  int16
}
