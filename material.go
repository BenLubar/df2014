package df2014

import (
	"reflect"
	"strconv"
)

type MaterialList struct {
	Type     []MaterialType
	Index    []int32 `df2014_assert_same_length_as:"Type" df2014_version_min:"1205"`
	Index23a []int16 `df2014_assert_same_length_as:"Type" df2014_version_max:"1169"`
}

func (ml MaterialList) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	convert := MaterialType.Convert
	if w.Header.Version < 1205 {
		convert = MaterialType.Convert23a
		ml.Index = make([]int32, len(ml.Index23a))
		for i, v := range ml.Index23a {
			ml.Index[i] = int32(v)
		}
	}
	if len(ml.Type) != len(ml.Index) {
		return prettyPrint(w, reflect.ValueOf(struct {
			Type  []MaterialType
			Index []int32
		}{ml.Type, ml.Index}), buf, indent, "")
	}

	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(ml.Type)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, mt := range ml.Type {
		m := convert(mt, ml.Index[i])

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(m), buf, indent, "")
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialType uint16

func (mt MaterialType) Convert(index int32) interface{} {
	switch {
	case mt == 0:
		return MaterialInorganic{uint16(mt), InorganicIndex(index)}
	case index < 0:
		return MaterialBuiltin{BuiltinIndex(mt), index}
	case 19 <= mt && mt < 219:
		return MaterialCreature{uint16(mt - 19), CreatureIndex(index)}
	case 219 <= mt && mt < 419:
		return MaterialFigure{uint16(mt - 219), index}
	case 419 <= mt && mt < 619:
		return MaterialPlant{uint16(mt - 419), PlantIndex(index)}
	default:
		return MaterialBuiltin{BuiltinIndex(mt), index}
	}
}

func (mt MaterialType) Convert23a(index int32) interface{} {
	switch mt {
	case 0: // WOOD
		return MaterialWood23a{MaterialType23a(mt), TreeIndex(index)}
	case 1: // STONE_GRAY
		return MaterialStone23a{MaterialType23a(mt), InorganicIndex(index)}
	case 2: // STONE_LIGHT
		return MaterialStone23a{MaterialType23a(mt), InorganicIndex(index)}
	case 3: // STONE_DARK
		return MaterialStone23a{MaterialType23a(mt), InorganicIndex(index)}
	case 4: // GOLD
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 5: // IRON
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 6: // SILVER
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 7: // COPPER
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 8: // GEM_ORNAMENTAL
		return MaterialGem23a{MaterialType23a(mt), GemIndex(index)}
	case 9: // GEM_SEMI
		return MaterialGem23a{MaterialType23a(mt), GemIndex(index)}
	case 10: // GEM_PRECIOUS
		return MaterialGem23a{MaterialType23a(mt), GemIndex(index)}
	case 11: // GEM_RARE
		return MaterialGem23a{MaterialType23a(mt), GemIndex(index)}
	case 12: // BONE
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 13: // IVORY
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 14: // JADE
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 15: // HORN
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 16: // AMBER
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 17: // CORAL
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 18: // PEARL
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 19: // SHELL
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 20: // LEATHER
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 21: // ADAMANTINE
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 22: // SILK
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 23: // PLANT
		return MaterialPlant23a{MaterialType23a(mt), PlantIndex(index)}
	case 24: // GLASS_GREEN
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 25: // GLASS_CLEAR
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 26: // GLASS_CRYSTAL
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 27: // SAND
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 28: // WATER
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 29: // ZINC
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 30: // TIN
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 31: // COAL
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 32: // BRONZE
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 33: // BRASS
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 34: // STEEL
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 35: // PIGIRON
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 36: // PLATINUM
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 37: // ELECTRUM
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 38: // POTASH
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 39: // ASH
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 40: // PEARLASH
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 41: // LYE
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 42: // RENDERED_FAT
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 43: // SOAP_ANIMAL
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 44: // FAT
		return MaterialCreature23a{MaterialType23a(mt), CreatureIndex(index)}
	case 45: // MUD
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 46: // VOMIT
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 47: // BLOOD
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	case 48: // SLIME
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	default:
		return MaterialBuiltin23a{MaterialType23a(mt), int16(index)}
	}
}

type MaterialType23a int16

var materialTypes23a = [...]string{
	0:  "WOOD",
	1:  "STONE_GRAY",
	2:  "STONE_LIGHT",
	3:  "STONE_DARK",
	4:  "GOLD",
	5:  "IRON",
	6:  "SILVER",
	7:  "COPPER",
	8:  "GEM_ORNAMENTAL",
	9:  "GEM_SEMI",
	10: "GEM_PRECIOUS",
	11: "GEM_RARE",
	12: "BONE",
	13: "IVORY",
	14: "JADE",
	15: "HORN",
	16: "AMBER",
	17: "CORAL",
	18: "PEARL",
	19: "SHELL",
	20: "LEATHER",
	21: "ADAMANTINE",
	22: "SILK",
	23: "PLANT",
	24: "GLASS_GREEN",
	25: "GLASS_CLEAR",
	26: "GLASS_CRYSTAL",
	27: "SAND",
	28: "WATER",
	29: "ZINC",
	30: "TIN",
	31: "COAL",
	32: "BRONZE",
	33: "BRASS",
	34: "STEEL",
	35: "PIGIRON",
	36: "PLATINUM",
	37: "ELECTRUM",
	38: "POTASH",
	39: "ASH",
	40: "PEARLASH",
	41: "LYE",
	42: "RENDERED_FAT",
	43: "SOAP_ANIMAL",
	44: "FAT",
	45: "MUD",
	46: "VOMIT",
	47: "BLOOD",
	48: "SLIME",
}

func (i MaterialType23a) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), materialTypes23a[:], buf)
}

type MaterialBuiltin23a struct {
	Type  MaterialType23a
	Index int16
}

type MaterialWood23a struct {
	Type  MaterialType23a
	Index TreeIndex
}

type MaterialStone23a struct {
	Type  MaterialType23a
	Index InorganicIndex
}

type MaterialGem23a struct {
	Type  MaterialType23a
	Index GemIndex
}

type MaterialCreature23a struct {
	Type  MaterialType23a
	Index CreatureIndex
}

type MaterialPlant23a struct {
	Type  MaterialType23a
	Index PlantIndex
}

type BuiltinIndex uint16

// keep in sync with builtin_mats in https://github.com/DFHack/df-structures/blob/master/df.materials.xml
var builtinMaterials = [...]string{
	"INORGANIC",
	"AMBER",
	"CORAL",
	"GLASS_GREEN",
	"GLASS_CLEAR",
	"GLASS_CRYSTAL",
	"WATER",
	"COAL",
	"POTASH",
	"ASH",
	"PEARLASH",
	"LYE",
	"MUD",
	"VOMIT",
	"SALT",
	"FILTH_B",
	"FILTH_Y",
	"UNKNOWN_SUBSTANCE",
	"GRIME",
}

func (i BuiltinIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), builtinMaterials[:], buf)
}

type MaterialBuiltin struct {
	T BuiltinIndex
	I int32
}

func (m MaterialBuiltin) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent, "")

	if int(m.T) < len(builtinMaterials) {
		buf = append(buf, " ("...)
		buf = append(buf, builtinMaterials[m.T]...)
		buf = append(buf, ')')
	}

	if m.I >= 0 {
		buf = append(buf, indent...)
		buf = append(buf, "I: "...)
		buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent, "")
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialInorganic struct {
	T uint16
	I InorganicIndex
}

type MaterialCreature struct {
	T uint16
	I CreatureIndex
}

func (m MaterialCreature) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent, "")

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent, "")

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialFigure struct {
	T uint16
	I int32
}

func (m MaterialFigure) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent, "")

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent, "")

	// TODO: human-readable indices

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialPlant struct {
	T uint16
	I PlantIndex
}

func (m MaterialPlant) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent, "")

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent, "")

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type ItemMaterialList struct {
	Type     []ItemType
	Sub      []int16        `df2014_assert_same_length_as:"Type"`
	Material []MaterialType `df2014_assert_same_length_as:"Type"`
	Index    []int32        `df2014_assert_same_length_as:"Type"`
}

func (iml ItemMaterialList) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	if len(iml.Type) != len(iml.Sub) || len(iml.Sub) != len(iml.Material) || len(iml.Material) != len(iml.Index) {
		return prettyPrint(w, reflect.ValueOf(struct {
			Type     []ItemType
			Sub      []int16
			Material []MaterialType
			Index    []int32
		}{iml.Type, iml.Sub, iml.Material, iml.Index}), buf, indent, outerTag)
	}

	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(iml.Type)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, t := range iml.Type {
		m := iml.Material[i].Convert(iml.Index[i])

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(t), buf, indent, "")
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(iml.Sub[i]), buf, indent, "")
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(m), buf, indent, "")
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type ItemType uint16

const (
	ITBar ItemType = iota
	ITSmallGem
	ITBlock
	ITRoughGem
	ITBoulder
	ITLog
	ITDoor
	ITFloodgate
	ITBed
	ITChair
	ITChain
	ITFlask
	ITGoblet
	ITInstrument
	ITToy
	ITWindow
	ITCage
	ITBarrel
	ITBucket
	ITAnimalTrap
	ITTable
	ITCoffin
	ITStatue
	ITCorpse
	ITWeapon
	ITArmor
	ITShoes
	ITShield
	ITHelm
	ITGloves
	ITBox
	ITBin
	ITArmorStand
	ITWeaponRack
	ITCabinet
	ITFigurine
	ITAmulet
	ITScepter
	ITAmmo
	ITCrown
	ITRing
	ITEarring
	ITBracelet
	ITLargeGem
	ITAnvil
	ITCorpsePiece
	ITRemains
	ITMeat
	ITFish
	ITFishRaw
	ITVermin
	ITPet
	ITSeeds
	ITPlant
	ITLeather
	ITPlantGrowth
	ITThread
	ITCloth
	ITTotem
	ITPants
	ITBackpack
	ITQuiver
	ITCatapultParts
	ITBallistaParts
	ITSiegeAmmo
	ITBallistaArrowHead
	ITMechanism
	ITTrapComponent
	ITDrink
	ITPowder
	ITCheese
	ITFood
	ITLiquid
	ITCoin
	ITGlob
	ITSmallRock
	ITPipeSection
	ITHatchCover
	ITGrate
	ITQuern
	ITMillstone
	ITSplint
	ITCrutch
	ITTractionBench
	ITOrthopedicCast
	ITTool
	ITSlab
	ITEgg
	ITBook
)

var itemTypes = [...]string{
	"bar",
	"small gem",
	"block",
	"rough gem",
	"boulder",
	"log",
	"door",
	"floodgate",
	"bed",
	"chair",
	"chain",
	"flask",
	"goblet",
	"instrument",
	"toy",
	"window",
	"cage",
	"barrel",
	"bucket",
	"animal trap",
	"table",
	"coffin",
	"statue",
	"corpse",
	"weapon",
	"armor",
	"shoes",
	"shield",
	"helm",
	"gloves",
	"box",
	"bin",
	"armor stand",
	"weapon rack",
	"cabinet",
	"figurine",
	"amulet",
	"scepter",
	"ammo",
	"crown",
	"ring",
	"earring",
	"bracelet",
	"large gem",
	"anvil",
	"corpse piece",
	"remains",
	"meat",
	"fish",
	"fish raw",
	"vermin",
	"pet",
	"seeds",
	"plant",
	"leather",
	"plant growth",
	"thread",
	"cloth",
	"totem",
	"pants",
	"backpack",
	"quiver",
	"catapult parts",
	"ballista parts",
	"siege ammo",
	"ballista arrow head",
	"mechanism",
	"trap component",
	"drink",
	"powder",
	"cheese",
	"food",
	"liquid",
	"coin",
	"glob",
	"small rock",
	"pipe section",
	"hatch cover",
	"grate",
	"quern",
	"millstone",
	"splint",
	"crutch",
	"traction bench",
	"orthopedic cast",
	"tool",
	"slab",
	"egg",
	"book",
}

func (i ItemType) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), itemTypes[:], buf)
}

type PlantGrowthList struct {
	Plant  []PlantIndex
	Growth []int16 `df2014_assert_same_length_as:"Plant"`
}

func (pgl PlantGrowthList) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	if len(pgl.Plant) != len(pgl.Growth) {
		return prettyPrint(w, reflect.ValueOf(struct {
			Plant  []PlantIndex
			Growth []int16
		}{pgl.Plant, pgl.Growth}), buf, indent, outerTag)
	}

	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(pgl.Plant)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, p := range pgl.Plant {
		g := pgl.Growth[i]

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(p), buf, indent, "")
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(g), buf, indent, "")
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}
