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

func (m MaterialBuiltin) prettyPrint(w *WorldDat, buf, indent []byte) []byte {
	buf = append(buf, '{')
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "T: "...)
	buf = prettyPrint(w, reflect.ValueOf(m.T), buf, indent)

	// keep in sync with builtin_mats in https://github.com/DFHack/df-structures/blob/master/df.materials.xml
	var builtin = [...]string{"INORGANIC", "AMBER", "CORAL", "GLASS_GREEN", "GLASS_CLEAR", "GLASS_CRYSTAL", "WATER", "COAL", "POTASH", "ASH", "PEARLASH", "LYE", "MUD", "VOMIT", "SALT", "FILTH_B", "FILTH_Y", "UNKNOWN_SUBSTANCE", "GRIME"}

	if m.T >= 0 && int(m.T) < len(builtin) {
		buf = append(buf, " ("...)
		buf = append(buf, builtin[m.T]...)
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
