package df2014

import (
	"reflect"
	"strconv"
)

func prettyPrintIndex(i int64, u uint64, definitions []string, buf []byte) []byte {
	buf = strconv.AppendInt(buf, i, 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, u, 16)
	buf = append(buf, ')')

	if i >= 0 && int(i) < len(definitions) {
		buf = append(buf, " ("...)
		buf = append(buf, definitions[i]...)
		buf = append(buf, ')')
	}

	return buf
}

func prettyPrintFlags(f uint64, definitions map[uint8]string, buf []byte) []byte {
	buf = strconv.AppendUint(buf, f, 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, f, 16)
	buf = append(buf, ')')

	for i := uint8(0); i < 64; i++ {
		if f&(1<<i) != 0 {
			buf = append(buf, " ("...)
			if d, ok := definitions[i]; ok {
				buf = append(buf, d...)
			} else {
				buf = append(buf, "unknown"...)
			}
			buf = append(buf, ')')
		}
	}

	return buf
}

type BodyIndex int32

func (i BodyIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Body, buf)
}

type BodyDetailPlanIndex int32

func (i BodyDetailPlanIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.BodyDetailPlan, buf)
}

type BodyGlossIndex int32

func (i BodyGlossIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.BodyGloss, buf)
}

type BuildingIndex int32

func (i BuildingIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Building, buf)
}

type ColorIndex int32

func (i ColorIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Color, buf)
}

type CreatureIndex int32

func (i CreatureIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Creature, buf)
}

type CreatureIndex16 uint16

func (i CreatureIndex16) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), w.StringTables.Creature, buf)
}

type CreatureVariationIndex int32

func (i CreatureVariationIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.CreatureVariation, buf)
}

type EntityIndex int32

func (i EntityIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Entity, buf)
}

type GemIndex int32

func (i GemIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Gem, buf)
}

type InorganicIndex int32

func (i InorganicIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Inorganic, buf)
}

type InteractionIndex int32

func (i InteractionIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Interaction, buf)
}

type ItemIndex int32

func (i ItemIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Item, buf)
}

type MaterialTemplateIndex int32

func (i MaterialTemplateIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.MaterialTemplate, buf)
}

type PatternIndex int32

func (i PatternIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Pattern, buf)
}

type PlantIndex int32

func (i PlantIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Plant, buf)
}

type ReactionIndex int32

func (i ReactionIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Reaction, buf)
}

type ShapeIndex int32

func (i ShapeIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Shape, buf)
}

type SymbolIndex int32

func (i SymbolIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Symbol, buf)
}

type TissueTemplateIndex int32

func (i TissueTemplateIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.TissueTemplate, buf)
}

type TranslationIndex int32

func (i TranslationIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Translation, buf)
}

type TreeIndex int32

func (i TreeIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Tree, buf)
}

type WordIndex int32

func (i WordIndex) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(uint32(i)), w.StringTables.Word, buf)
}
