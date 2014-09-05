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

	Unk006 MaterialList
	Unk008 uint32 `df2014_assert_equals:"0x0"`
	Unk009 uint32 `df2014_assert_equals:"0x0"`
	Unk010 MaterialList
	Unk012 uint32 `df2014_assert_equals:"0x0"`
	Unk013 uint32 `df2014_assert_equals:"0x0"`
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
	Unk042 []uint32
	Unk043 []uint16 `df2014_assert_same_length_as:"Unk042"`
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

	Unk143 uint32 `df2014_assert_equals:"0x1"`
	Unk144 uint32
	Unk145 uint16 // flags?

	Unk146 uint32 `df2014_assert_equals:"0x1"`
	Unk147 uint8

	Unk148 uint32 `df2014_assert_equals:"0x1"`
	Unk149 bool

	Unk150 uint32 `df2014_assert_equals:"0x1"`
	Unk151 uint16

	Unk152 uint32 `df2014_assert_equals:"0x1"`
	Unk153 uint32

	Unk154 uint32 `df2014_assert_equals:"0x0"`
	Unk155 uint32 `df2014_assert_equals:"0x0"`
	Unk156 uint32 `df2014_assert_equals:"0x0"`
	Unk157 uint32 `df2014_assert_equals:"0x0"`
	Unk158 uint32 `df2014_assert_equals:"0x0"`
	Unk159 uint32 `df2014_assert_equals:"0x0"`
	Unk160 uint32 `df2014_assert_equals:"0x0"`
	Unk161 uint32 `df2014_assert_equals:"0x0"`
	Unk162 uint32 `df2014_assert_equals:"0x0"`

	Unk163 uint16 `df2014_assert_equals:"0xd"`
	Unk164 uint16 `df2014_assert_equals:"0x10"`
	Unk165 uint16 `df2014_assert_equals:"0x10"`
	Unk166 uint16 `df2014_assert_equals:"0x1"`
	Unk167 uint16 `df2014_assert_equals:"0x1"`
	Unk168 uint16 `df2014_assert_equals:"0xf"`
	Unk169 uint16 `df2014_assert_equals:"0x0"`
	Unk170 uint16 `df2014_assert_equals:"0x1"`
	Unk171 uint16 `df2014_assert_equals:"0xf"`
	Unk172 uint16 `df2014_assert_equals:"0xf"`
	Unk173 uint16 `df2014_assert_equals:"0x0"`
	Unk174 uint16 `df2014_assert_equals:"0x0"`
	Unk175 uint16 `df2014_assert_equals:"0x0"`
	Unk176 uint16 `df2014_assert_equals:"0x0"`
	Unk177 uint16 `df2014_assert_equals:"0x0"`
	Unk178 uint16 `df2014_assert_equals:"0x2"`
	Unk179 uint16 `df2014_assert_equals:"0xf"`
	Unk180 uint16 `df2014_assert_equals:"0xf"`
	Unk181 uint16 `df2014_assert_equals:"0xf"`
	Unk182 uint16 `df2014_assert_equals:"0xf"`
	Unk183 uint16 `df2014_assert_equals:"0xf"`
	Unk184 uint16 `df2014_assert_equals:"0xf"`

	Unk185 uint32 `df2014_assert_equals:"0x0"`
	Unk186 uint32 `df2014_assert_equals:"0x0"`
	Unk187 uint32 `df2014_assert_equals:"0x0"`
	Unk188 uint32 `df2014_assert_equals:"0x0"`
	Unk189 uint32 `df2014_assert_equals:"0x0"`
	Unk190 uint32 `df2014_assert_equals:"0x0"`
	Unk191 uint32 `df2014_assert_equals:"0x0"`
	Unk192 uint32 `df2014_assert_equals:"0x0"`
	Unk193 uint32 `df2014_assert_equals:"0x0"`
	Unk194 uint32 `df2014_assert_equals:"0x0"`
	Unk195 uint32 `df2014_assert_equals:"0x0"`
	Unk196 uint32 `df2014_assert_equals:"0x0"`
	Unk197 uint32 `df2014_assert_equals:"0x0"`
	Unk198 uint32 `df2014_assert_equals:"0x0"`
	Unk199 uint32 `df2014_assert_equals:"0x0"`
	Unk200 uint32 `df2014_assert_equals:"0x0"`
	Unk201 uint32 `df2014_assert_equals:"0x0"`
	Unk202 uint32 `df2014_assert_equals:"0x0"`
	Unk203 uint32 `df2014_assert_equals:"0x0"`
	Unk204 uint32 `df2014_assert_equals:"0x0"`
	Unk205 uint32 `df2014_assert_equals:"0x0"`
	Unk206 uint32 `df2014_assert_equals:"0x0"`
	Unk207 uint32 `df2014_assert_equals:"0x0"`
	Unk208 uint32 `df2014_assert_equals:"0x0"`
	Unk209 uint32 `df2014_assert_equals:"0x0"`
	Unk210 uint32 `df2014_assert_equals:"0x0"`
	Unk211 uint32 `df2014_assert_equals:"0x0"`
	Unk212 uint32 `df2014_assert_equals:"0x0"`
	Unk213 uint32 `df2014_assert_equals:"0x0"`
	Unk214 uint32 `df2014_assert_equals:"0x0"`
	Unk215 uint32 `df2014_assert_equals:"0x0"`
	Unk216 uint32 `df2014_assert_equals:"0x0"`
	Unk217 uint32 `df2014_assert_equals:"0x0"`
	Unk218 uint32 `df2014_assert_equals:"0x0"`
	Unk219 uint32 `df2014_assert_equals:"0x0"`
	Unk220 uint32 `df2014_assert_equals:"0x0"`
	Unk221 uint32 `df2014_assert_equals:"0x0"`
	Unk222 uint32 `df2014_assert_equals:"0x0"`
	Unk223 uint32 `df2014_assert_equals:"0x0"`
	Unk224 uint32 `df2014_assert_equals:"0x0"`
	Unk225 uint32 `df2014_assert_equals:"0x0"`
	Unk226 uint32 `df2014_assert_equals:"0x0"`
	Unk227 uint32 `df2014_assert_equals:"0x0"`
	Unk228 uint32 `df2014_assert_equals:"0x0"`
	Unk229 uint32 `df2014_assert_equals:"0x0"`
	Unk230 uint32 `df2014_assert_equals:"0x0"`
	Unk231 int32  `df2014_assert_equals:"-1"`
	Unk232 uint32 `df2014_assert_equals:"0x0"`
	Unk233 int32  `df2014_assert_equals:"-1"`
	Unk234 uint32 `df2014_assert_equals:"0x0"`
	Unk235 uint32 `df2014_assert_equals:"0x0"`
	Unk236 uint32 `df2014_assert_equals:"0x0"`
	Unk237 uint32 `df2014_assert_equals:"0x0"`
	Unk238 uint32 `df2014_assert_equals:"0x0"`
	Unk239 uint32 `df2014_assert_equals:"0x0"`
	Unk240 uint32 `df2014_assert_equals:"0x0"`
	Unk241 uint32 `df2014_assert_equals:"0x0"`
	Unk242 uint32 `df2014_assert_equals:"0x0"`
	Unk243 uint32 `df2014_assert_equals:"0x0"`
	Unk244 uint32 `df2014_assert_equals:"0x0"`
	Unk245 uint32 `df2014_assert_equals:"0x0"`
	Unk246 uint32 `df2014_assert_equals:"0x0"`
	Unk247 uint32 `df2014_assert_equals:"0x0"`
	Unk248 uint32 `df2014_assert_equals:"0x0"`
	Unk249 uint32 `df2014_assert_equals:"0x0"`
	Unk250 uint8  `df2014_assert_equals:"0x0"`

	Unk251 []EntityUnk251

	Unk252 uint32 `df2014_assert_equals:"0x0"`
	Unk253 uint32 `df2014_assert_equals:"0x0"`
	Unk254 uint32 `df2014_assert_equals:"0x0"`
	Unk255 uint32 `df2014_assert_equals:"0x0"`
	Unk256 uint32 `df2014_assert_equals:"0x0"`
	Unk257 uint32 `df2014_assert_equals:"0x0"`
}

type EntityUnk251 struct {
	Unk000 uint32 `df2014_assert_equals:"0x4"`
	Unk001 uint32
	Unk002 uint32
	Unk003 int32 `df2014_assert_equals:"-1"`
	Unk004 uint32
	Unk005 int32  `df2014_assert_equals:"-1"`
	Unk006 uint32 `df2014_assert_same_as:"Unk004"`
	Unk007 int32  `df2014_assert_equals:"-1"`
	Unk008 uint32 `df2014_assert_equals:"0x0"`
}
