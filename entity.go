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

var entityTypes = [...]string{
	"civilization",
	"site government",
	"vessel crew",
	"migrating group",
	"nomadic group",
	"religion",
	"military unit",
	"outcast",
}

func (i EntityType) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if int(i) < len(entityTypes) {
		buf = append(buf, " ("...)
		buf = append(buf, entityTypes[i]...)
		buf = append(buf, ')')
	}

	return buf
}

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

type RaceCasteList struct {
	Race  []uint32
	Caste []uint16 `df2014_assert_same_length_as:"Race"`
}

func (rcl RaceCasteList) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(rcl.Race)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, r := range rcl.Race {
		c := rcl.Caste[i]

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = strconv.AppendUint(buf, uint64(r), 10)
		buf = append(buf, " (0x"...)
		buf = strconv.AppendUint(buf, uint64(r), 16)
		buf = append(buf, ')')
		if r >= 0 && int(r) < len(w.StringTables.Creature) {
			buf = append(buf, " ("...)
			buf = append(buf, w.StringTables.Creature[r]...)
			buf = append(buf, ')')
		}
		buf = append(buf, ": "...)
		buf = strconv.AppendUint(buf, uint64(c), 10)
		buf = append(buf, " (0x"...)
		buf = strconv.AppendUint(buf, uint64(c), 16)
		buf = append(buf, ')')
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type Entity struct {
	Type         EntityType
	ID           uint32
	Class        string
	Unk000       uint16 //`df2014_assert_equals:"0x19"`
	Unk001       uint16 //`df2014_assert_equals:"0x4b"`
	SaveFileID   uint32
	NextMemberID uint16
	Name         *Name
	Creature     EntityCreatureIndex
	Flags        uint32

	Materials EntityMaterials

	Fish         RaceCasteList
	Egg          RaceCasteList
	Plant        MaterialList
	Unk048       uint32 `df2014_assert_equals:"0x0"`
	Unk049       uint32 `df2014_assert_equals:"0x0"`
	Unk050       []uint32
	Unk051       []uint16 `df2014_assert_same_length_as:"Unk050"`
	Seed         MaterialList
	WoodProducts ItemMaterialList
	Pet          RaceCasteList
	Wagon        RaceCasteList
	PackAnimal   RaceCasteList
	WagonPuller  RaceCasteList
	Mount        RaceCasteList
	Minion       RaceCasteList
	ExoticPet    RaceCasteList
	Wood         MaterialList

	Unk074 []uint32
	Unk075 []uint32
	Unk076 []uint32

	Unk077 MaterialList
	Unk079 MaterialList
	Unk081 uint32 `df2014_assert_equals:"0x0"`
	Unk082 uint32 `df2014_assert_equals:"0x0"`
	Unk083 MaterialList
	Unk085 MaterialList
	Unk087 uint32 `df2014_assert_equals:"0x0"`
	Unk088 uint32 `df2014_assert_equals:"0x0"`
	Unk089 MaterialList
	Unk091 MaterialList
	Unk093 MaterialList

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
	Unk136 []EntityUnk136
	Unk137 []EntityUnk137
	Unk138 []uint32

	Unk139   uint32 `df2014_assert_equals:"0x1"`
	Unk139_1 uint32
	Unk139_2 uint32 `df2014_assert_equals:"0x0"`
	Unk140   uint32 `df2014_assert_equals:"0x1"`
	Unk140_1 uint16
	Unk141   uint32 `df2014_assert_equals:"0x1"`
	Unk141_1 uint32
	Unk142   uint32 `df2014_assert_equals:"0x1"`
	Unk142_1 uint16

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

type EntityMaterials struct {
	Leather  MaterialList
	Fiber    MaterialList
	Silk     MaterialList
	Wool     MaterialList
	Craft    MaterialList
	Unk002   MaterialList
	Barrel   MaterialList
	Flask    MaterialList
	Quiver   MaterialList
	Backpack MaterialList
	Cage     MaterialList
	Wood     MaterialList
	Ore      MaterialList
	Drink    MaterialList
	Cheese   MaterialList
	Powder   MaterialList
	Extract  MaterialList
	Meat     MaterialList
}

type EntityUnk136 struct {
	Unk000 uint16 `df2014_assert_equals:"0x1"`
	Unk001 uint32
	Unk002 uint16 // flags?
}

type EntityUnk137 struct {
	Unk000 uint32
	Unk001 uint32
	Unk002 uint32
	Unk003 uint32
	Unk004 int32 `df2014_assert_equals:"-1"`
	Unk005 int32 `df2014_assert_equals:"-1"`
	Unk006 uint32
	Unk007 uint32
	Unk008 uint32
	Unk009 uint32
	Unk010 uint32
	Unk011 uint32
	Unk012 uint32
	Unk013 uint32
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
