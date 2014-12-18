package df2014

type WorldDat struct {
	Header

	Unused000 int16 `df2014_assert_equals:"0" df2014_version_min:"1205"`
	Unused001 int32 `df2014_assert_equals:"0" df2014_version_max:"1169"`

	// SaveSlot is the original region number of this world. A value of 0 means
	// the save will be titled region1. 1 means region2, and so on.
	SaveSlot int32 `df2014_assert_gte:"0"`

	NextIDs WorldNextID

	Name   *Name
	Unk050 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk051 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk052 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk054 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk055 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk056 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk057 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk058 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk059 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk060 int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk061 int16 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk062 int16 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk063 int16 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	Unk064 int16 `df2014_assert_gte:"-1" df2014_version_min:"1205"`

	// all I know about the next five variables is that they add up to 15 bytes.
	Unk100 int8  `df2014_assert_equals:"1"`
	Unk101 int16 `df2014_assert_equals:"0"`
	Unk102 int32 `df2014_assert_equals:"1"`
	Unk103 int32 `df2014_assert_equals:"0"`
	Unk104 int32 `df2014_assert_equals:"0"`

	Title string `df2014_version_min:"1110"`

	//GeneratedRaws WorldGeneratedRaws `df2014_version_min:"1205"`
	StringTables WorldStringTables

	ItemIDs     map[uint32]bool
	BuildingIDs map[uint32]bool
	EntityIDs   map[uint32]bool
	NemesisIDs  map[uint32]bool
	ArtifactIDs map[uint32]bool
	CoinBatches int32 `df2014_assert_gte:"0"`
	TaskTypes   []int16

	Items    []Item   `df2014_get_length_from:"ItemIDs" df2014_assert_id_set:"ItemIDs"`
	Entities []Entity `df2014_get_length_from:"EntityIDs" df2014_assert_id_set:"EntityIDs"`
}

type WorldNextID struct {
	Unit          int32 `df2014_assert_gte:"-1"`
	Item          int32 `df2014_assert_gte:"-1"`
	Entity        int32 `df2014_assert_gte:"-1"`
	Nemesis       int32 `df2014_assert_gte:"-1"`
	Artifact      int32 `df2014_assert_gte:"-1"`
	Building      int32 `df2014_assert_gte:"-1"`
	HistFigure    int32 `df2014_assert_gte:"-1"`
	HistEvent     int32 `df2014_assert_gte:"-1"`
	UnitChunk     int32 `df2014_assert_gte:"-1"`
	ArtImageChunk int32 `df2014_assert_gte:"-1"`
	Task          int32 `df2014_assert_gte:"-1"`
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
	Tree              []string `df2014_version_max:"1169"`
	Inorganic         []string
	Gem               []string `df2014_version_max:"1169"`
	Plant             []string
	Body              []string
	BodyGloss         []string
	Creature          []string
	Item              []string
	Building          []string `df2014_version_min:"1205"`
	Entity            []string
	Word              []string
	Symbol            []string
	Translation       []string
	Color             []string `df2014_version_min:"1139"`
	Shape             []string `df2014_version_min:"1139"`
	Pattern           []string `df2014_version_min:"1205"`
	Reaction          []string `df2014_version_min:"1205"`
	MaterialTemplate  []string `df2014_version_min:"1205"`
	TissueTemplate    []string `df2014_version_min:"1205"`
	BodyDetailPlan    []string `df2014_version_min:"1205"`
	CreatureVariation []string `df2014_version_min:"1205"`
	Interaction       []string `df2014_version_min:"1205"`
}
