package df2014

type WorldDat struct {
	Header

	Unk000         uint16
	Unk001         [28]int32
	Name           *Name
	Unk002         uint8
	Unk003         int16
	Unk004         [3]int32
	TranslatedName string

	GeneratedRaws struct {
		Inorganic        [][]string
		Item             [][]string
		Creature         [][]string
		Interaction      [][]string
		EntityLayer      [][]string
		InteractionLayer [][]string
		LanguageLayer    [][]string
	}

	StringTables struct {
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

	Unk005 [][14]uint32
	Unk006 map[uint32]uint32
	Unk007 [19]map[uint32]bool

	Books []Book
}
