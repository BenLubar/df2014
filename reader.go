package df2014

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/BenLubar/df2014/versions"
	"github.com/BenLubar/df2014/wtf23a"
)

const MaxAlloc = 1 << 30 // allocations over 1 gigabyte are probably mistakes.

type NestedError struct {
	Message string
	Inner   error
}

func (err NestedError) Error() string {
	return err.Message + "\n" + err.Inner.Error()
}

type Reader struct {
	io.Reader
}

func (r *Reader) DecodeSimple(v interface{}) error {
	return r.DecodeValue(nil, reflect.ValueOf(v).Elem())
}

func (r *Reader) Decode(v *WorldDat) error {
	return r.DecodeValue(v, reflect.ValueOf(v).Elem())
}

func (r *Reader) DecodeValue(world *WorldDat, v reflect.Value) (err error) {
	switch v.Kind() {
	case reflect.Struct:
		if _, ok := v.Interface().(Header); ok {
			h, err := r.header()
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(h))
			return nil
		}
		if _, ok := v.Interface().(Book); ok {
			b, err := r.book()
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(b))
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			fieldTag := v.Type().Field(i).Tag
			if tag := fieldTag.Get("df2014_version_min"); tag != "" {
				expected, err := strconv.ParseUint(tag, 0, 32)
				if err != nil {
					return err
				}

				if world.Version < versions.Version(expected) {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_version_max"); tag != "" {
				expected, err := strconv.ParseUint(tag, 0, 32)
				if err != nil {
					return err
				}

				if world.Version > versions.Version(expected) {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_type"); tag != "" {
				expected, err := strconv.ParseInt(tag, 0, 64)
				if err != nil {
					return err
				}

				if v.FieldByName("Type").Int() != expected {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_get_length_from"); tag != "" {
				l := v.FieldByName(tag).Len()
				v.Field(i).Set(reflect.MakeSlice(v.Type().Field(i).Type, l, l))

				for j := 0; j < l; j++ {
					err := r.DecodeValue(world, v.Field(i).Index(j))
					if err != nil {
						v.Field(i).SetLen(j + 1) // hide the entries that we didn't get to
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), err}}
					}
				}
			} else {
				err := r.DecodeValue(world, v.Field(i))
				if err != nil {
					return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), err}
				}
			}

			if tag := fieldTag.Get("df2014_assert_same_length_as"); tag != "" {
				expected, actual := v.FieldByName(tag).Len(), v.Field(i).Len()
				if expected != actual {
					return fmt.Errorf("df2014: len(%s) %d != len(%s) %d", v.Type().Field(i).Name, actual, tag, expected)
				}
			}
			if tag := fieldTag.Get("df2014_assert_same_as"); tag != "" {
				actual := fmt.Sprintf("%#v", v.Field(i).Interface())
				expected := fmt.Sprintf("%#v", v.FieldByName(tag).Interface())
				if expected != actual {
					return fmt.Errorf("df2014: %s %q != %s %q", v.Type().Field(i).Name, actual, tag, expected)
				}
			}
			if expected := fieldTag.Get("df2014_assert_equals"); expected != "" {
				actual := fmt.Sprintf("%#v", v.Field(i).Interface())
				if expected != actual {
					return fmt.Errorf("df2014: %s: %q != %q", v.Type().Field(i).Name, actual, expected)
				}
			}
			if tag := fieldTag.Get("df2014_assert_gte"); tag != "" {
				switch v.Field(i).Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					actual := v.Field(i).Int()
					expected, err := strconv.ParseInt(tag, 0, v.Field(i).Type().Bits())
					if err != nil {
						return err
					}
					if actual < expected {
						return fmt.Errorf("df2014: %s: %d ≱ %d", v.Type().Field(i).Name, actual, expected)
					}

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					actual := v.Field(i).Uint()
					expected, err := strconv.ParseUint(tag, 0, v.Field(i).Type().Bits())
					if err != nil {
						return err
					}
					if actual < expected {
						return fmt.Errorf("df2014: %s: %d ≱ %d", v.Type().Field(i).Name, actual, expected)
					}
				}
			}
			if tag := fieldTag.Get("df2014_assert_lte"); tag != "" {
				switch v.Field(i).Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					actual := v.Field(i).Int()
					expected, err := strconv.ParseInt(tag, 0, v.Field(i).Type().Bits())
					if err != nil {
						return err
					}
					if actual > expected {
						return fmt.Errorf("df2014: %s: %d ≰ %d", v.Type().Field(i).Name, actual, expected)
					}

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					actual := v.Field(i).Uint()
					expected, err := strconv.ParseUint(tag, 0, v.Field(i).Type().Bits())
					if err != nil {
						return err
					}
					if actual > expected {
						return fmt.Errorf("df2014: %s: %d ≰ %d", v.Type().Field(i).Name, actual, expected)
					}
				}
			}
			if tag := fieldTag.Get("df2014_assert_id_set"); tag != "" {
				for j, l := 0, v.Field(i).Len(); j < l; j++ {
					expected := v.FieldByName(tag).Index(j).Int()
					actual := v.Field(i).Index(j).FieldByName("ID").Int()

					if expected != actual {
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), fmt.Errorf("df2014: id in set (%d) does not match id (%d)", expected, actual)}}
					}
				}
			}
			if tag := fieldTag.Get("df2014_assert_id_parent"); tag != "" {
				expected := v.FieldByName("ID").Uint()
				for j, l := 0, v.Field(i).Len(); j < l; j++ {
					actual := v.Field(i).Index(j).FieldByName(tag).Uint()
					if expected != actual {
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), fmt.Errorf("df2014: id in parent (%d) does not match id in %s (%d)", expected, tag, actual)}}
					}
				}
			}
			if tag := fieldTag.Get("df2014_assert_next_id"); tag != "" {
				expected := reflect.ValueOf(world.NextIDs).FieldByName(tag).Int()
				set := make(map[int64]bool, v.Field(i).Len())
				for j, l := 0, v.Field(i).Len(); j < l; j++ {
					actual := v.Field(i).Index(j).Int()
					if actual < 0 || actual > expected {
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), fmt.Errorf("df2014: next %s id (%d) is invalid for id (%d)", tag, expected, actual)}}
					}
					if set[actual] {
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), fmt.Errorf("df2014: duplicate %s id (%d)", tag, actual)}}
					}
					set[actual] = true
				}
			}
		}
		return nil

	case reflect.Ptr:
		flag, err := r.bool()
		if err != nil {
			return err
		}
		if !flag {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		v.Set(reflect.New(v.Type().Elem()))
		return r.DecodeValue(world, v.Elem())

	case reflect.String:
		s, err := r.string()
		if err != nil {
			return err
		}
		v.SetString(s)
		return nil

	case reflect.Slice:
		var length int32 // signed so huge numbers cause errors instead of allocating tons of memory
		err := binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return err
		}
		if length < 0 {
			return fmt.Errorf("df2014: negative length (%d)", length)
		}

		if size := v.Type().Size() * uintptr(length); size > MaxAlloc {
			return fmt.Errorf("df2014: huge alloc (%d bytes)", size)
		}
		v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))

		fallthrough
	case reflect.Array:
		for i, l := 0, v.Len(); i < l; i++ {
			err := r.DecodeValue(world, v.Index(i))
			if err != nil {
				return NestedError{fmt.Sprintf("at index %d", i), err}
			}
		}
		return nil

	case reflect.Bool:
		b, err := r.bool()
		if err != nil {
			return err
		}
		v.SetBool(b)
		return nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return binary.Read(r, binary.LittleEndian, v.Addr().Interface())

	case reflect.Map:
		var length int32 // signed so huge numbers cause errors instead of allocating tons of memory
		err := binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return err
		}
		if length < 0 {
			return fmt.Errorf("df2014: negative length (%d)", length)
		}

		v.Set(reflect.MakeMap(v.Type()))

		var prev reflect.Value
		check := func(next reflect.Value) error {
			if prev.IsValid() {
				switch prev.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if prev.Int() > next.Int() {
						return fmt.Errorf("df2014: values not in order: %v > %v", prev.Interface(), next.Interface())
					}
					if prev.Int() == next.Int() {
						return fmt.Errorf("df2014: duplicate value: %v", prev.Interface())
					}

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if prev.Uint() > next.Uint() {
						return fmt.Errorf("df2014: values not in order: %v > %v", prev.Interface(), next.Interface())
					}
					if prev.Uint() == next.Uint() {
						return fmt.Errorf("df2014: duplicate value: %v", prev.Interface())
					}
				}
			}
			prev = next
			return nil
		}

		for i := int32(0); i < length; i++ {
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()

			err = r.DecodeValue(world, key)
			if err != nil {
				return err
			}
			err = r.DecodeValue(world, val)
			if err != nil {
				return err
			}

			v.SetMapIndex(key, val)

			if err = check(key); err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("df2014: unexpected value type %v: %v", v.Kind(), v.Type())
}

func (r *Reader) bool() (b bool, err error) {
	var n uint8
	err = binary.Read(r, binary.LittleEndian, &n)
	if err == nil {
		switch n {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = fmt.Errorf("df2014: unexpected value for bool: %d", n)
		}
	}
	return
}

var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\u00A0")

func (r *Reader) string() (string, error) {
	var length int16 // signed so huge numbers cause errors instead of allocating tons of memory
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("df2014: negative length (%d)", length)
	}

	buf, s := make([]byte, length), make([]rune, length)

	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	for i, b := range buf {
		s[i] = cp437[b]
	}

	return string(s), nil
}

type CompressionType uint32

const (
	Uncompressed CompressionType = iota
	ZLib
	// the rest aren't used in files
	Special23a
	Special40d
)

var compressionTypeNames = []string{
	Uncompressed: "Uncompressed",
	ZLib:         "ZLib",
	Special23a:   "special: 23a",
	Special40d:   "special: 40d",
}

func (i CompressionType) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	return prettyPrintIndex(int64(i), uint64(i), compressionTypeNames, buf)
}

func (i CompressionType) String() string {
	if i < CompressionType(len(compressionTypeNames)) {
		return compressionTypeNames[i]
	}

	return strconv.FormatUint(uint64(i), 10)
}

type Header struct {
	Version     versions.Version
	Compression CompressionType
}

func (r *Reader) header() (h Header, err error) {
	err = binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return
	}

	switch h.Compression {
	case Uncompressed:
		// nothing to be done
	case ZLib:
		r.Reader = &compression1Reader{r: r.Reader}
	default:
		var undo bytes.Buffer
		err = binary.Write(&undo, binary.LittleEndian, &h)
		if err != nil {
			panic(err)
		}
		if h.Version < 10000 && h.Compression&0xffff == 0xda78 {
			// guess 40d compression
			r.Reader = &compression1Reader{r: io.MultiReader(&undo, r.Reader)}
			h.Compression = Special40d
		} else {
			// guess 23a compression
			r.Reader = wtf23a.NewReader(io.MultiReader(&undo, r.Reader))
			h.Compression = Special23a
		}
		err = binary.Read(r, binary.LittleEndian, &h.Version)
		if err != nil {
			return
		}
	}

	if h.Version == 0 && h.Compression == 0 {
		// guess 23a compression with modded-out encryption
		var undo bytes.Buffer
		err = binary.Write(&undo, binary.LittleEndian, &h)
		if err != nil {
			panic(err)
		}
		r.Reader = wtf23a.NewReader(io.MultiReader(&undo, r.Reader))
		h.Compression = Special23a
		err = binary.Read(r, binary.LittleEndian, &h.Version)
		if err != nil {
			return
		}
	}

	if !h.Version.IsKnown() {
		err = fmt.Errorf("df2014: unhandled version %d", h.Version)
		return
	}

	return
}

type compression1Reader struct {
	r    io.Reader
	buf  bytes.Buffer
	nlen []byte
	next []byte
}

func (r *compression1Reader) Read(b []byte) (n int, err error) {
	if r.buf.Len() == 0 {
		err = r.fill()
		if err != nil {
			return
		}
	}

	return r.buf.Read(b)
}

func (r *compression1Reader) fill() (err error) {
	var n int
	if r.next == nil {
		if r.nlen == nil {
			r.nlen = make([]byte, 0, 4)
		}
		n, err = io.ReadFull(r.r, r.nlen[len(r.nlen):cap(r.nlen)])
		r.nlen = r.nlen[:len(r.nlen)+n]
		if err != nil {
			return
		}
		length := int32(binary.LittleEndian.Uint32(r.nlen))
		if length < 0 {
			return fmt.Errorf("df2014: negative length (%d)", length)
		}
		r.next = make([]byte, 0, length)
		r.nlen = r.nlen[:0]
	}
	if len(r.next) != cap(r.next) {
		n, err = io.ReadFull(r.r, r.next[len(r.next):cap(r.next)])
		r.next = r.next[:len(r.next)+n]
		if err != nil {
			return
		}
	}

	z, err := zlib.NewReader(bytes.NewReader(r.next))
	if err != nil {
		return
	}
	defer func() {
		e := z.Close()
		if err == nil {
			err = e
		}
	}()

	_, err = io.Copy(&r.buf, z)
	if err != nil {
		return
	}
	r.next = nil
	return
}
