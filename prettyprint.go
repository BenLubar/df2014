package df2014

import (
	"reflect"
	"sort"
	"strconv"

	"github.com/BenLubar/df2014/versions"
)

type prettyPrinter interface {
	prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte
}

func prettyPrint(w *WorldDat, v reflect.Value, buf, indent []byte, outerTag reflect.StructTag) []byte {
	bufPreType := buf
	buf = append(buf, v.Kind().String()...)
	buf = append(buf, ' ')
	if v.Type().String() != v.Kind().String() {
		buf = append(buf, v.Type().String()...)
		buf = append(buf, ' ')
	}

	if p, ok := v.Interface().(prettyPrinter); ok {
		return p.prettyPrint(w, buf, indent, outerTag)
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			buf = append(buf, "(nil)"...)
		} else {
			buf = append(bufPreType, '&')
			buf = prettyPrint(w, v.Elem(), buf, indent, outerTag)
		}

	case reflect.Struct:
		buf = append(buf, '{')
		indent = append(indent, '\t')

		for i, l := 0, v.NumField(); i < l; i++ {
			fieldTag := v.Type().Field(i).Tag
			if tag := fieldTag.Get("df2014_version_min"); tag != "" {
				expected, err := strconv.ParseUint(tag, 0, 32)
				if err != nil {
					panic(err)
				}

				if w.Version < versions.Version(expected) {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_version_max"); tag != "" {
				expected, err := strconv.ParseUint(tag, 0, 32)
				if err != nil {
					panic(err)
				}

				if w.Version > versions.Version(expected) {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_type"); tag != "" {
				expected, err := strconv.ParseInt(tag, 0, 64)
				if err != nil {
					panic(err)
				}

				if v.FieldByName("Type").Int() != expected {
					continue
				}
			}
			buf = append(buf, indent...)
			buf = append(buf, v.Type().Field(i).Name...)
			buf = append(buf, ": "...)
			buf = prettyPrint(w, v.Field(i), buf, indent, fieldTag)
		}

		buf = append(buf, indent[:len(indent)-1]...)
		buf = append(buf, '}')

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf = strconv.AppendInt(buf, v.Int(), 10)
		buf = append(buf, " (0x"...)
		buf = strconv.AppendUint(buf, uint64(v.Int())&(1<<uint(v.Type().Bits())-1), 16)
		buf = append(buf, ')')

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf = strconv.AppendUint(buf, v.Uint(), 10)
		buf = append(buf, " (0x"...)
		buf = strconv.AppendUint(buf, v.Uint(), 16)
		buf = append(buf, ')')

	case reflect.String:
		buf = append(buf, "(len = "...)
		buf = strconv.AppendInt(buf, int64(len([]rune(v.String()))), 10)
		buf = append(buf, ") "...)
		buf = strconv.AppendQuote(buf, v.String())

	case reflect.Bool:
		buf = strconv.AppendBool(buf, v.Bool())

	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice {
			buf = append(buf, "(len = "...)
			buf = strconv.AppendInt(buf, int64(v.Len()), 10)
			buf = append(buf, ") {"...)
		} else {
			buf = append(buf, '{')
		}
		indent = append(indent, '\t')

		var names []string
		if tag := outerTag.Get("df2014_key_is_string"); tag != "" {
			names = reflect.ValueOf(w.StringTables).FieldByName(tag).Interface().([]string)
		}
		for i, l := 0, v.Len(); i < l; i++ {
			buf = append(buf, indent...)
			buf = strconv.AppendInt(buf, int64(i), 10)
			if i >= 0 && int(i) < len(names) {
				buf = append(buf, " ("...)
				buf = append(buf, names[i]...)
				buf = append(buf, ')')
			}
			buf = append(buf, ": "...)
			buf = prettyPrint(w, v.Index(i), buf, indent, "")
		}

		buf = append(buf, indent[:len(indent)-1]...)
		buf = append(buf, '}')

	case reflect.Map:
		buf = append(buf, "(len = "...)
		buf = strconv.AppendInt(buf, int64(v.Len()), 10)
		buf = append(buf, ") {"...)
		indent = append(indent, '\t')

		var elements mapElements
		keys := v.MapKeys()
		for _, key := range keys {
			elements = append(elements, mapElement{key: key, value: v.MapIndex(key)})
		}

		sort.Sort(elements)

		for _, e := range elements {
			buf = append(buf, indent...)
			buf = prettyPrint(w, e.key, buf, indent, "")
			buf = append(buf, ": "...)
			buf = prettyPrint(w, e.value, buf, indent, "")
		}

		buf = append(buf, indent[:len(indent)-1]...)
		buf = append(buf, '}')

	default:
		panic("unhandled type: " + v.Kind().String() + " " + v.Type().String())
	}

	return buf
}

type mapElement struct {
	key, value reflect.Value
}

type mapElements []mapElement

func (m mapElements) Len() int {
	return len(m)
}
func (m mapElements) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m mapElements) Less(i, j int) bool {
	switch m[i].key.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return m[i].key.Int() < m[j].key.Int()

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return m[i].key.Uint() < m[j].key.Uint()
	}
	panic("unreachable")
}
