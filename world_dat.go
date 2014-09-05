package df2014

type WorldDat struct {
	Header

	Unk000         uint16 `df2014_assert_equals:"0x0"`
	Unk001         int32  `df2014_assert_gte:"-1"`
	Unk002         int32  `df2014_assert_gte:"-1"`
	Unk003         int32  `df2014_assert_gte:"-1"`
	Unk004         int32  `df2014_assert_gte:"-1"`
	Unk005         int32  `df2014_assert_gte:"-1"`
	Unk006         int32  `df2014_assert_gte:"-1"`
	Unk007         int32  `df2014_assert_gte:"-1"`
	Unk008         int32  `df2014_assert_gte:"-1"`
	Unk009         int32  `df2014_assert_gte:"-1"`
	Unk010         int32  `df2014_assert_gte:"-1"`
	Unk011         int32  `df2014_assert_gte:"-1"`
	Unk012         int32  `df2014_assert_gte:"-1"`
	Unk013         int32  `df2014_assert_gte:"-1"`
	Unk014         int32  `df2014_assert_gte:"-1"`
	Unk015         int32  `df2014_assert_gte:"-1"`
	Unk016         int32  `df2014_assert_gte:"-1"`
	Unk017         int32  `df2014_assert_gte:"-1"`
	Unk018         int32  `df2014_assert_gte:"-1"`
	Unk019         int32  `df2014_assert_gte:"-1"`
	Unk020         int32  `df2014_assert_gte:"-1"`
	Unk021         int32  `df2014_assert_gte:"-1"`
	Unk022         int32  `df2014_assert_gte:"-1"`
	Unk023         int32  `df2014_assert_gte:"-1"`
	Unk024         int32  `df2014_assert_gte:"-1"`
	Unk025         int32  `df2014_assert_gte:"-1"`
	Unk026         int32  `df2014_assert_gte:"-1"`
	Unk027         int32  `df2014_assert_gte:"-1"`
	Unk028         int32  `df2014_assert_gte:"-1"`
	Name           *Name
	Unk029         uint8
	Unk030         int16
	Unk031         int32
	Unk032         int32
	Unk033         int32
	TranslatedName string

	GeneratedRaws WorldGeneratedRaws
	StringTables  WorldStringTables

	Unk034 []WorldDatUnk034
	Unk035 map[uint32]uint32
	Unk036 map[uint32]bool
	Unk037 map[uint32]bool
	Unk038 map[uint32]bool
	Unk039 map[uint32]bool `df2014_assert_same_length_as:"Unk035"`
	Unk040 map[uint32]bool
	Unk041 map[uint32]bool
	Unk042 map[uint32]bool
	Unk043 map[uint32]bool
	Unk044 map[uint32]bool
	Unk045 map[uint32]bool
	Unk046 map[uint32]bool
	Unk047 map[uint32]bool
	Unk048 map[uint32]bool
	Unk049 map[uint32]bool
	Unk050 map[uint32]bool
	Unk051 map[uint32]bool
	Unk052 map[uint32]bool
	Unk053 map[uint32]bool
	Unk054 map[uint32]bool

	Books    []Book `df2014_assert_same_length_as:"Unk035"`
	Entities Entity
}

type WorldGeneratedRaws struct {
	Inorganic   [][]string
	Unk000      [][]string
	Item        [][]string
	Creature    [][]string
	Entity      [][]string
	Interaction [][]string
	Language    [][]string
}

type WorldStringTables struct {
	Inorganic         []string
	Plant             []string
	Body              []string
	BodyGloss         []string
	Creature          []string
	Item              []string
	Building          []string
	Entity            []string
	Word              []string
	Symbol            []string
	Translation       []string
	Color             []string
	Shape             []string
	Pattern           []string
	Reaction          []string
	MaterialTemplate  []string
	TissueTemplate    []string
	BodyDetailPlan    []string
	CreatureVariation []string
	Interaction       []string
}

type WorldDatUnk034 struct {
	Unk000 uint32
	Unk001 uint32
	Unk002 uint32
	Unk003 uint32
	Unk004 uint32
	Unk005 uint32
	Unk006 uint32
	Unk007 uint32
	Unk008 uint32
	Unk009 uint32
	Unk010 uint32
	Unk011 uint32
	Unk012 uint32
	Unk013 uint32
}
