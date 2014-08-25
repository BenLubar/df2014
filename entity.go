package df2014

import (
	"strconv"
)

type EntityType uint16

const (
	Civilization EntityType = iota
	SiteGovernment
	VesselCrew
	MigratingGroup
	NomadicGroup
	Religion
	MilitaryUnit
	Outcast
)

type EntityCreatureIndex uint16

func (i EntityCreatureIndex) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if int(i) < len(w.StringTables.Creature) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Creature[i]...)
		buf = append(buf, ')')
	}

	return buf
}

type Entity struct {
	Type     EntityType `df2014_assert_equals:"0x0"`
	ID       uint32
	Class    string
	Unk000   uint16 `df2014_assert_equals:"0x19"`
	Unk001   uint16 `df2014_assert_equals:"0x4b"`
	Unk002   uint32
	Unk003   uint16 `df2014_assert_equals:"0x0"`
	Name     *Name
	Creature EntityCreatureIndex
	Unk004   uint32 `df2014_assert_equals:"0x0"`

	Unk006 []uint16
	Unk007 []uint32 `df2014_assert_same_length_as:"Unk006"`
	Unk008 uint32   `df2014_assert_equals:"0x0"`
	Unk009 uint32   `df2014_assert_equals:"0x0"`
	Unk010 uint32   `df2014_assert_equals:"0x0"`
	Unk011 uint32   `df2014_assert_equals:"0x0"`
	Unk012 uint32   `df2014_assert_equals:"0x0"`
	Unk013 uint32   `df2014_assert_equals:"0x0"`
	Unk014 MaterialList
	Unk016 MaterialList
	Unk018 MaterialList
	Unk020 MaterialList
	Unk022 MaterialList
	Unk024 MaterialList
	Unk026 MaterialList
	Unk028 MaterialList
	Unk030 MaterialList
	Unk032 MaterialList
	Unk034 MaterialList
	Unk036 MaterialList
	Unk038 MaterialList
	Unk040 MaterialList
	Unk042 uint32 `df2014_assert_equals:"0x0"`
	Unk043 uint32 `df2014_assert_equals:"0x0"`
	Unk044 []uint32
	Unk045 []uint16 `df2014_assert_same_length_as:"Unk044"`
	Unk046 uint32   `df2014_assert_equals:"0x0"`
	Unk047 uint32   `df2014_assert_equals:"0x0"`
	Unk048 uint32   `df2014_assert_equals:"0x0"`
	Unk049 uint32   `df2014_assert_equals:"0x0"`
	Unk050 uint32   `df2014_assert_equals:"0x0"`
	Unk051 uint32   `df2014_assert_equals:"0x0"`
	Unk052 uint32   `df2014_assert_equals:"0x0"`
	Unk053 uint32   `df2014_assert_equals:"0x0"`
	Unk054 uint32   `df2014_assert_equals:"0x0"`
	Unk055 uint32   `df2014_assert_equals:"0x0"`
	Unk056 uint32   `df2014_assert_equals:"0x0"`
	Unk057 uint32   `df2014_assert_equals:"0x0"`
	Unk058 []uint32
	Unk059 []uint16 `df2014_assert_same_length_as:"Unk058"`
	Unk060 []uint32
	Unk061 []uint16 `df2014_assert_same_length_as:"Unk060"`
	Unk062 uint32   `df2014_assert_equals:"0x0"`
	Unk063 uint32   `df2014_assert_equals:"0x0"`
	Unk064 uint32   `df2014_assert_equals:"0x0"`
	Unk065 uint32   `df2014_assert_equals:"0x0"`
	Unk066 []uint32
	Unk067 []uint16 `df2014_assert_same_length_as:"Unk066"`
	Unk068 []uint32
	Unk069 []uint16 `df2014_assert_same_length_as:"Unk068"`
	Unk070 []uint32
	Unk071 []uint16 `df2014_assert_same_length_as:"Unk070"`

	Unk072 []uint16
	Unk073 []uint32 `df2014_assert_same_length_as:"Unk072"`

	Unk074 uint32 `df2014_assert_equals:"0x0"`
	Unk075 uint32 `df2014_assert_equals:"0x0"`
	Unk076 uint32 `df2014_assert_equals:"0x0"`

	Unk077 []uint16
	Unk078 []uint32 `df2014_assert_same_length_as:"Unk077"`
	Unk079 uint32   `df2014_assert_equals:"0x0"`
	Unk080 uint32   `df2014_assert_equals:"0x0"`
	Unk081 uint32   `df2014_assert_equals:"0x0"`
	Unk082 uint32   `df2014_assert_equals:"0x0"`
	Unk083 []uint16
	Unk084 []uint32 `df2014_assert_same_length_as:"Unk083"`
	Unk085 []uint16
	Unk086 []uint32 `df2014_assert_same_length_as:"Unk085"`
	Unk087 uint32   `df2014_assert_equals:"0x0"`
	Unk088 uint32   `df2014_assert_equals:"0x0"`
	Unk089 []uint16
	Unk090 []uint32 `df2014_assert_same_length_as:"Unk089"`
	Unk091 []uint16
	Unk092 []uint32 `df2014_assert_same_length_as:"Unk091"`
	Unk093 []uint16
	Unk094 []uint32 `df2014_assert_same_length_as:"Unk093"`

	Unk095 int16 `df2014_assert_equals:"-1"`
	Unk096 uint32
	Unk097 int16 `df2014_assert_equals:"-1"`
	Unk098 uint32
	Unk099 int16 `df2014_assert_equals:"-1"`
	Unk100 uint32

	Unk101 []uint16
	Unk102 []uint16
	Unk103 []uint16
	Unk104 []uint16
	Unk105 []uint16
	Unk106 []uint16
	Unk107 []uint16
	Unk108 []uint16
	Unk109 []uint16
	Unk110 []uint16
	Unk111 []uint16
	Unk112 []uint16
	Unk113 []uint16
	Unk114 []uint16
	Unk115 []uint16

	Unk116 []uint16
	Unk117 []uint32 `df2014_assert_same_length_as:"Unk116"`
	Unk118 []uint16
	Unk119 []uint32 `df2014_assert_same_length_as:"Unk118"`
	Unk120 []uint16
	Unk121 []uint32 `df2014_assert_same_length_as:"Unk120"`
	Unk122 []uint16
	Unk123 []uint32 `df2014_assert_same_length_as:"Unk122"`
	Unk124 []uint16
	Unk125 []uint32 `df2014_assert_same_length_as:"Unk124"`
	Unk126 []uint16
	Unk127 []uint32 `df2014_assert_same_length_as:"Unk126"`
	Unk128 []uint16
	Unk129 []uint32 `df2014_assert_same_length_as:"Unk128"`
	Unk130 []uint16
	Unk131 []uint32 `df2014_assert_same_length_as:"Unk130"`
	Unk132 []uint16
	Unk133 []uint32 `df2014_assert_same_length_as:"Unk132"`

	Unk134 uint32 `df2014_assert_equals:"0x0"`
	Unk135 uint32 `df2014_assert_equals:"0x0"`
	Unk136 uint32 `df2014_assert_equals:"0x0"`
	Unk137 uint32 `df2014_assert_equals:"0x0"`
	Unk138 uint32 `df2014_assert_equals:"0x0"`

	Unk139 []uint64
	Unk140 []uint16 `df2014_assert_same_length_as:"Unk139"`
	Unk141 []uint32 `df2014_assert_same_length_as:"Unk139"`
	Unk142 []uint16 `df2014_assert_same_length_as:"Unk139"`
}
