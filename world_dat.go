package df2014

type WorldDat struct {
	Header

	Unused000 int16 `df2014_assert_equals:"0" df2014_version_min:"1205"`
	Unused001 int32 `df2014_assert_equals:"0" df2014_version_max:"1169"`

	// SaveSlot is the original region number of this world. A value of 0 means
	// the save will be titled region1. 1 means region2, and so on.
	SaveSlot int32 `df2014_assert_gte:"0" df2014_version_max:"1169"`

	NextIDs WorldNextID

	Name *Name

	// all I know about the next five variables is that they add up to 15 bytes.
	Unk100 int8  `df2014_assert_equals:"1"`
	Unk101 int16 `df2014_assert_equals:"0"`
	Unk102 int32 `df2014_assert_equals:"1"`
	Unk103 int32 `df2014_assert_equals:"0"`
	Unk104 int32 `df2014_assert_equals:"0"`

	Title string `df2014_version_min:"1110"`

	GeneratedRaws WorldGeneratedRaws `df2014_version_min:"1287"`
	StringTables  WorldStringTables

	STOP struct{} `df2014_assert_equals:"STOP" df2014_version_min:"1205"`

	ItemIDs     []int32 `df2014_assert_next_id:"Item"`
	BuildingIDs []int32 `df2014_assert_next_id:"Building"`
	EntityIDs   []int32 `df2014_assert_next_id:"Entity"`
	NemesisIDs  []int32 `df2014_assert_next_id:"Nemesis"`
	ArtifactIDs []int32 `df2014_assert_next_id:"Artifact"`
	CoinBatches int32   `df2014_assert_gte:"0"`
	TaskTypes   []int16

	Items    []Item   `df2014_get_length_from:"ItemIDs" df2014_assert_id_set:"ItemIDs"`
	Entities []Entity `df2014_get_length_from:"EntityIDs" df2014_assert_id_set:"EntityIDs"`
}

type WorldNextID struct {
	Unit                int32 `df2014_assert_gte:"-1"`
	Unit2               int32 `df2014_assert_gte:"-1" df2014_version_min:"1287"`
	Item                int32 `df2014_assert_gte:"-1"`
	Entity              int32 `df2014_assert_gte:"-1"`
	Nemesis             int32 `df2014_assert_gte:"-1"`
	Artifact            int32 `df2014_assert_gte:"-1"`
	Building            int32 `df2014_assert_gte:"-1"`
	Unk007              int32 `df2014_assert_gte:"-1" df2014_version_min:"1205"`
	HistFigure          int32 `df2014_assert_gte:"-1"`
	HistEvent           int32 `df2014_assert_gte:"-1"`
	HistEventCollection int32 `df2014_assert_gte:"-1" df2014_version_min:"1254"`
	UnitChunk           int32 `df2014_assert_gte:"-1"`
	ArtImageChunk       int32 `df2014_assert_gte:"-1"`
	Task                int32 `df2014_assert_gte:"-1"`
	Squad               int32 `df2014_assert_gte:"-1" df2014_version_min:"1287"`
	Schedule            int32 `df2014_assert_gte:"-1" df2014_version_min:"1287"`
	Activity            int32 `df2014_assert_gte:"-1" df2014_version_min:"1287"`
	InteractionInstance int32 `df2014_assert_gte:"-1" df2014_version_min:"1372"`
	WrittenContent      int32 `df2014_assert_gte:"-1" df2014_version_min:"1372"`
	Identity            int32 `df2014_assert_gte:"-1" df2014_version_min:"1372"`
	Incident            int32 `df2014_assert_gte:"-1" df2014_version_min:"1372"`
	Crime               int32 `df2014_assert_gte:"-1" df2014_version_min:"1372"`
	Vehicle             int32 `df2014_assert_gte:"-1" df2014_version_min:"1400"`
	Unk023              int32 `df2014_assert_gte:"-1" df2014_version_min:"1441"`
	Unk024              int32 `df2014_assert_gte:"-1" df2014_version_min:"1441"`
	Unk025              int32 `df2014_assert_gte:"-1" df2014_version_min:"1441"`
	Unk026              int32 `df2014_assert_gte:"-1" df2014_version_min:"1441"`
	Unk027              int32 `df2014_assert_gte:"-1" df2014_version_min:"1441"`
}

type WorldGeneratedRaws struct {
	Inorganic   [][]string `df2014_version_min:"1372"`
	Unk000      [][]string `df2014_version_min:"1400"`
	Item        [][]string `df2014_version_min:"1441"`
	Creature    [][]string
	Entity      [][]string `df2014_version_min:"1441"`
	Interaction [][]string `df2014_version_min:"1372"`
	Language    [][]string `df2014_version_min:"1441"`
}

type WorldStringTables struct {
	Tree              []string `df2014_version_max:"1268"`
	Inorganic         []string
	Gem               []string `df2014_version_max:"1169"`
	Metal             []string `df2014_version_min:"1205" df2014_version_max:"1268"`
	Plant             []string
	Body              []string
	BodyGloss         []string
	Creature          []string
	Item              []string
	Building          []string `df2014_version_min:"1287"`
	Entity            []string
	Word              []string
	Symbol            []string
	Translation       []string
	Color             []string `df2014_version_min:"1139"`
	Shape             []string `df2014_version_min:"1139"`
	Pattern           []string `df2014_version_min:"1205"`
	Reaction          []string `df2014_version_min:"1287"`
	MaterialTemplate  []string `df2014_version_min:"1287"`
	TissueTemplate    []string `df2014_version_min:"1287"`
	BodyDetailPlan    []string `df2014_version_min:"1287"`
	CreatureVariation []string `df2014_version_min:"1287"`
	Interaction       []string `df2014_version_min:"1287"`
}
