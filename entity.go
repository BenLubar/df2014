package df2014

import (
	"reflect"
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

var entityTypes = []string{
	"civilization",
	"site government",
	"vessel crew",
	"migrating group",
	"nomadic group",
	"religion",
	"military unit",
	"outcast",
}

func (i EntityType) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), entityTypes, buf)
}

type EntityFlags uint32

var entityFlags = map[uint8]string{
	4: "ruin",
}

func (f EntityFlags) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintFlags(uint64(f), entityFlags, buf)
}

type RaceCasteList struct {
	Race  []uint32
	Caste []uint16 `df2014_assert_same_length_as:"Race"`
}

func (rcl RaceCasteList) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = append(buf, "(len = "...)
	buf = strconv.AppendInt(buf, int64(len(rcl.Race)), 10)
	buf = append(buf, ") {"...)

	indent = append(indent, '\t')

	for i, r := range rcl.Race {
		c := rcl.Caste[i]

		buf = append(buf, indent...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ": "...)
		buf = CreatureIndex(r).prettyPrint(w, buf, indent, "")
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

type EntityEthicResponse uint16

const (
	EERNotApplicable EntityEthicResponse = iota
	EERAcceptable
	EERPersonalMatter
	EERJustifiedIfNoRepercussions
	EERJustifiedIfGoodReason
	EERJustifiedIfExtremeReason
	EERJustifiedIfSelfDefense
	EEROnlyIfSanctioned
	EERMisguided
	EERShun
	EERAppalling
	EERPunishReprimand
	EERPunishSerious
	EERPunishExile
	EERPunishCapital
	EERUnthinkable
	EERRequired
)

var entityEthicResponses = [...]string{
	"not applicable",
	"acceptable",
	"personal matter",
	"justified if no repercussions",
	"justified if good reason",
	"justified if extreme reason",
	"justified if self defense",
	"only if sanctioned",
	"misguided",
	"shun",
	"appalling",
	"punish reprimand",
	"punish serious",
	"punish exile",
	"punish capital",
	"unthinkable",
	"required",
}

func (i EntityEthicResponse) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), entityEthicResponses[:], buf)
}

type Entity struct {
	Type  EntityType
	ID    uint32
	Class string

	MinTemperature int16 `df2014_version_min:"1131"`
	MaxTemperature int16 `df2014_version_min:"1131"`

	UnitChunkID    uint32
	UnitChunkSubID uint16

	Unused_41C []uint32 `df2014_assert_equals:"[]uint32{}"`

	Name  *Name
	Race  CreatureIndex16
	Flags EntityFlags

	Leather      []int16
	Cloth        []int16
	Silk         []int16
	Crafts       MaterialList
	Glass        MaterialList
	Wood         MaterialList
	Pets         []int16
	Cages        MaterialList
	Drinks       MaterialList
	Cheese       MaterialList
	Mill         MaterialList
	Extract      MaterialList
	MeatFish     []int16
	Plants       []int16
	Wagons       []int16
	PackAnimals  []int16
	WagonPullers []int16
	Mounts       []int16
	SiegeMinions []int16

	LyeWood     []int16      `df2014_version_min:"1147"`
	CraftMetals []int16      `df2014_version_min:"1147"`
	Stone       MaterialList `df2014_version_min:"1147"`
	Gem         MaterialList `df2014_version_min:"1147"`
	Bones       []int16      `df2014_version_min:"1147"`
	Shells      []int16      `df2014_version_min:"1147"`
	Pearls      []int16      `df2014_version_min:"1147"`
	Ivory       []int16      `df2014_version_min:"1147"`
	Horn        []int16      `df2014_version_min:"1147"`
	Other       MaterialList `df2014_version_min:"1147"`
	SpecialMat0 Material
	SpecialMat1 Material
	SpecialMat2 Material

	Diggers     []int16 `df2014_version_min:"1130"`
	Weapons     []int16 `df2014_version_min:"1130"`
	Armor       []int16 `df2014_version_min:"1130"`
	Ammo        []int16 `df2014_version_min:"1130"`
	Helms       []int16 `df2014_version_min:"1130"`
	Gloves      []int16 `df2014_version_min:"1130"`
	Shoes       []int16 `df2014_version_min:"1130"`
	Pants       []int16 `df2014_version_min:"1130"`
	Shields     []int16 `df2014_version_min:"1130"`
	TrapComps   []int16 `df2014_version_min:"1130"`
	Toys        []int16 `df2014_version_min:"1130"`
	Instruments []int16 `df2014_version_min:"1130"`
	SiegeAmmo   []int16 `df2014_version_min:"1130"`

	WeaponMaterials MaterialList
	ArmorMaterials  MaterialList

	/*
		ActivityStats *struct {
			Stats                 EntityActivityStatistics
			LastCommunicateSeason int16
			LastCommunicateYear   int32
		}
		Imports, Exports, Offerings int32
		OfferingsHistory            [10]int32
		HostileLevel, SiegeCount    int32
	*/

	Discoveries     EntityDiscoveredMatGlosses
	MeatFishRecipes []EntityRecipe
	OtherRecipes    []EntityRecipe
	OwnedItems      []int32

	/*
		MeetingEvents []EntityMeetingEvent
	*/

	EntityLinks []EntityEntityLink
	SiteLinks   []EntitySiteLink
	FigureIDs   []int32
	NemesisIDs  []int32

	ArtImageTypes  []int16 `df2014_version_min:"1147"`
	ArtImageIDs    []int32 `df2014_version_min:"1147"`
	ArtImageSubIDs []int16 `df2014_version_min:"1147"`

	Uniforms []EntityUniform `df2014_version_min:"1147"`
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

type EntityResources struct {
	Fish         RaceCasteList
	Egg          RaceCasteList
	Plant        MaterialList
	Orchard      PlantGrowthList
	Garden       PlantGrowthList
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
	Metal        []InorganicIndex
	Stone        []InorganicIndex
	Gem          []InorganicIndex
	Bone         MaterialList
	Shell        MaterialList
	Pearl        MaterialList
	Ivory        MaterialList
	Horn         MaterialList
	Unk004       MaterialList
	Sand         MaterialList
	Glass        MaterialList
	Clay         MaterialList
}

type EntityEntityLinkType uint16

const (
	EntityParent EntityEntityLinkType = iota
	EntityChild
)

var entityEntityLinkTypes = [...]string{
	"parent",
	"child",
}

func (i EntityEntityLinkType) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = strconv.AppendUint(buf, uint64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if int(i) < len(entityEntityLinkTypes) {
		buf = append(buf, " ("...)
		buf = append(buf, entityEntityLinkTypes[i]...)
		buf = append(buf, ')')
	}

	return buf
}

type EntityEntityLink struct {
	Type     EntityEntityLinkType
	ID       uint32
	Strength uint16 `df2014_assert_lte:"100"`
}

type EntitySiteLinkType uint16

type EntitySiteLink struct {
	Type     EntitySiteLinkType
	ID       uint32
	Strength uint16 `df2014_assert_lte:"100"`
}

type EntityUnk143 struct {
	Unk000 uint32
	Unk001 uint16 // flags?
}

type EntityUnk153 struct {
	Unk000 uint32 `df2014_assert_equals:"0x0"`
	Unk001 uint16 `df2014_assert_equals:"0x3"`
}

type EntityUnk156 struct {
	Unk000 bool  `df2014_assert_equals:"true"`
	Unk001 int8  `df2014_assert_equals:"-1"`
	Unk002 int16 `df2014_assert_equals:"-1"`
}

type EntityUnk157 struct {
	Unk000 uint16 `df2014_assert_equals:"0x108"`
	Unk001 int16  `df2014_assert_equals:"-1"`
	Unk002 uint16 `df2014_assert_equals:"0x4"`
	Unk003 int32  `df2014_assert_equals:"-1"`
	Unk004 int32  `df2014_assert_equals:"-1"`
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

type EntityDiscoveredMatGlosses struct {
	CreatureFoods []bool `df2014_key_is_string:"Creature"`
	Creatures     []bool `df2014_key_is_string:"Creature"`
	PlantFoods    []bool `df2014_key_is_string:"Plant"`
	Plants        []bool `df2014_key_is_string:"Plant"`
	Trees         []bool `df2014_key_is_string:"Tree"`
}

type EntityRecipe struct {
	SubType      int16
	ItemTypes    []int16
	ItemSubtypes []int16
	Materials    []int16
	MatGlosses   []int16
}

type EntityUniform struct {
	ID           int16
	ItemTypes    []int16
	ItemSubtypes []int16
	ItemInfo     []EntityUniformItem
}

type EntityUniformItem struct {
	RandomDye          int8
	ArmorLevel         int8
	ItemColor          int16
	ArtImageChunk      int32
	ArtImageID         int16
	ImageThreadColor   int16
	ImageMaterialClass int16
	MakerRace          int16 `df2014_version_min:"1164"`
}
