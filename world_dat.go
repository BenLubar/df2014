package df2014

import (
	"fmt"
)

type WorldDat struct {
	Version     uint32
	Compression uint32

	Unk000 uint16
	Unk001 [28]int32
	Name   Name
	Unk002 uint8
	Unk003 int16
	Unk004 [3]int32
	Name2  string

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
	Unk007 [19][]uint32
}

func (r *Reader) WorldDat() (w WorldDat, err error) {
	w.Version, w.Compression, err = r.Header()
	if err != nil {
		return
	}

	w.Unk000, err = r.Uint16()
	if err != nil {
		return
	}
	for i := range w.Unk001 {
		w.Unk001[i], err = r.Int32()
		if err != nil {
			return
		}
	}
	ok, err := r.Bool()
	if err != nil {
		return
	}
	if ok {
		w.Name, err = r.Name()
		if err != nil {
			return
		}
	}

	w.Unk002, err = r.Uint8()
	if err != nil {
		return
	}

	w.Unk003, err = r.Int16()
	if err != nil {
		return
	}

	for i := range w.Unk004 {
		w.Unk004[i], err = r.Int32()
		if err != nil {
			return
		}
	}

	w.Name2, err = r.String()
	if err != nil {
		return
	}

	stringList := func() (l []string, err error) {
		length, err := r.Uint32()
		if err != nil {
			return
		}

		l = make([]string, length)
		for i := range l {
			l[i], err = r.String()
			if err != nil {
				return
			}
		}

		return
	}
	stringListList := func() (l [][]string, err error) {
		length, err := r.Uint32()
		if err != nil {
			return
		}

		l = make([][]string, length)
		for i := range l {
			l[i], err = stringList()
			if err != nil {
				return
			}
		}

		return
	}

	w.GeneratedRaws.Inorganic, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.Item, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.Creature, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.Interaction, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.EntityLayer, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.InteractionLayer, err = stringListList()
	if err != nil {
		return
	}
	w.GeneratedRaws.LanguageLayer, err = stringListList()
	if err != nil {
		return
	}

	w.StringTables.Inorganic, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Plant, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Body, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.BodyGloss, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Creature, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Item, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Building, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Entity, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Word, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Symbol, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Translation, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Color, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Shape, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Pattern, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Reaction, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.MaterialTemplate, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.TissueTemplate, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.BodyDetailPlan, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.CreatureVariation, err = stringList()
	if err != nil {
		return
	}
	w.StringTables.Interaction, err = stringList()
	if err != nil {
		return
	}

	length, err := r.Uint32()
	w.Unk005 = make([][14]uint32, length)
	for i := range w.Unk005 {
		for j := range w.Unk005[i] {
			w.Unk005[i][j], err = r.Uint32()
			if err != nil {
				return
			}
		}
	}

	uint32Uint32Map := func() (m map[uint32]uint32, err error) {
		length, err := r.Uint32()
		if err != nil {
			return
		}

		m = make(map[uint32]uint32, length)
		for i := uint32(0); i < length; i++ {
			var k, v uint32
			k, err = r.Uint32()
			if err != nil {
				return
			}

			if _, ok := m[k]; ok {
				err = fmt.Errorf("df2014: map already contains key %d (index %d)", k, i)
				return
			}

			v, err = r.Uint32()
			if err != nil {
				return
			}

			m[k] = v
		}

		return
	}

	w.Unk006, err = uint32Uint32Map()
	if err != nil {
		return
	}

	uint32OrderedList := func() (l []uint32, err error) {
		length, err := r.Uint32()
		if err != nil {
			return
		}

		l = make([]uint32, length)
		for i := range l {
			l[i], err = r.Uint32()
			if err != nil {
				return
			}
			if i > 0 && l[i-1] > l[i] {
				err = fmt.Errorf("df2014: list out of order at index %d: %d > %d", i, l[i-1], l[i])
				return
			}
			if i > 0 && l[i-1] == l[i] {
				err = fmt.Errorf("df2014: list contains duplicate element at index %d: %d", i, l[i])
				return
			}
		}

		return
	}

	for i := range w.Unk007 {
		w.Unk007[i], err = uint32OrderedList()
		if err != nil {
			return
		}
	}

	return
}
