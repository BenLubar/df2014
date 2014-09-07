package df2014

import (
	"reflect"
	"strconv"
)

type MaterialList struct {
	Type  []MaterialType
	Index []int32 `df2014_assert_same_length_as:"Type"`
}

func (ml MaterialList) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(ml.Type)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, mt := range ml.Type {
		m := mt.Convert(ml.Index[i])

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(m), buf, indent)
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialType uint16

func (mt MaterialType) Convert(index int32) interface{} {
	switch {
	case mt == 0:
		return MaterialInorganic{uint16(mt), index}
	case index < 0:
		return MaterialBuiltin{uint16(mt), index}
	case 19 <= mt && mt < 219:
		return MaterialCreature{uint16(mt - 19), index}
	case 219 <= mt && mt < 419:
		return MaterialFigure{uint16(mt - 219), index}
	case 419 <= mt && mt < 619:
		return MaterialPlant{uint16(mt - 419), index}
	default:
		return MaterialBuiltin{uint16(mt), index}
	}
}

type MaterialBuiltin struct {
	T uint16
	I int32
}

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

func (m MaterialBuiltin) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	if int(m.T) < len(builtinMaterials) {
		buf = append(buf, " ("...)
		buf = append(buf, builtinMaterials[m.T]...)
		buf = append(buf, ')')
	}

	if m.I >= 0 {
		buf = append(buf, indent...)
		buf = append(buf, "I: "...)
		buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent)
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialInorganic struct {
	T uint16
	I int32
}

func (m MaterialInorganic) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent)

	if m.I >= 0 && int(m.I) < len(w.StringTables.Inorganic) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Inorganic[m.I]...)
		buf = append(buf, ')')
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialCreature struct {
	T uint16
	I int32
}

func (m MaterialCreature) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent)

	if m.I >= 0 && int(m.I) < len(w.StringTables.Creature) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Creature[m.I]...)
		buf = append(buf, ')')
	}

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialFigure struct {
	T uint16
	I int32
}

func (m MaterialFigure) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent)

	// TODO: human-readable indices

	buf = append(buf, indent[:len(indent)-1]...)
	buf = append(buf, '}')

	return buf
}

type MaterialPlant struct {
	T uint16
	I int32
}

func (m MaterialPlant) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	// TODO: human-readable types

	buf = append(buf, indent...)
	buf = append(buf, "I: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.I), buf, indent)

	if m.I >= 0 && int(m.I) < len(w.StringTables.Plant) {
		buf = append(buf, " ("...)
		buf = append(buf, w.StringTables.Plant[m.I]...)
		buf = append(buf, ')')
	}

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

func (iml ItemMaterialList) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(iml.Type)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, t := range iml.Type {
		m := iml.Material[i].Convert(iml.Index[i])

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(t), buf, indent)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(iml.Sub[i]), buf, indent)
		buf = append(buf, ": "...)
		buf = prettyPrint(w, reflect.ValueOf(m), buf, indent)
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

func (i ItemType) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = strconv.AppendInt(buf, int64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if i >= 0 && int(i) < len(itemTypes) {
		buf = append(buf, " ("...)
		buf = append(buf, itemTypes[i]...)
		buf = append(buf, ')')
	}

	return buf
}
