package df2014

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
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

				if world.Version < SaveVersion(expected) {
					continue
				}
			}
			if tag := fieldTag.Get("df2014_version_max"); tag != "" {
				expected, err := strconv.ParseUint(tag, 0, 32)
				if err != nil {
					return err
				}

				if world.Version > SaveVersion(expected) {
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

			if tag := v.Type().Field(i).Tag.Get("df2014_assert_same_length_as"); tag != "" {
				expected, actual := v.FieldByName(tag).Len(), v.Field(i).Len()
				if expected != actual {
					return fmt.Errorf("df2014: len(%s) %d != len(%s) %d", v.Type().Field(i).Name, actual, tag, expected)
				}
			}
			if tag := v.Type().Field(i).Tag.Get("df2014_assert_same_as"); tag != "" {
				actual := fmt.Sprintf("%#v", v.Field(i).Interface())
				expected := fmt.Sprintf("%#v", v.FieldByName(tag).Interface())
				if expected != actual {
					return fmt.Errorf("df2014: %s %q != %s %q", v.Type().Field(i).Name, actual, tag, expected)
				}
			}
			if expected := v.Type().Field(i).Tag.Get("df2014_assert_equals"); expected != "" {
				actual := fmt.Sprintf("%#v", v.Field(i).Interface())
				if expected != actual {
					return fmt.Errorf("df2014: %s: %q != %q", v.Type().Field(i).Name, actual, expected)
				}
			}
			if tag := v.Type().Field(i).Tag.Get("df2014_assert_gte"); tag != "" {
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
			if tag := v.Type().Field(i).Tag.Get("df2014_assert_lte"); tag != "" {
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
			if tag := v.Type().Field(i).Tag.Get("df2014_assert_id_set"); tag != "" {
				want := make(map[interface{}]bool, v.FieldByName(tag).Len())
				for _, id := range v.FieldByName(tag).MapKeys() {
					want[id.Interface()] = true
				}

				have := make(map[interface{}]bool, v.Field(i).Len())
				for j, l := 0, v.Field(i).Len(); j < l; j++ {
					have[v.Field(i).Index(j).FieldByName("ID").Interface()] = true
				}

				missing := make(map[interface{}]bool)
				for id := range want {
					if !have[id] {
						missing[id] = true
					}
				}
				unexpected := make(map[interface{}]bool)
				for id := range have {
					if !want[id] {
						unexpected[id] = true
					}
				}

				if len(missing) > 0 || len(unexpected) > 0 {
					return fmt.Errorf("df2014: %s: ids missing=%v unexpected=%v", v.Type().Field(i).Name, missing, unexpected)
				}
			}
			if tag := v.Type().Field(i).Tag.Get("df2014_assert_id_parent"); tag != "" {
				expected := v.FieldByName("ID").Uint()
				for j, l := 0, v.Field(i).Len(); j < l; j++ {
					actual := v.Field(i).Index(j).FieldByName(tag).Uint()
					if expected != actual {
						return NestedError{fmt.Sprintf("in struct field %q", v.Type().Field(i).Name), NestedError{fmt.Sprintf("at index %d", j), fmt.Errorf("df2014: id in parent (%d) does not match id in %s (%d)", expected, tag, actual)}}
					}
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

		if v.Type().Elem().Kind() == reflect.Bool {
			// it's a set
			trueVal := reflect.New(v.Type().Elem()).Elem()
			trueVal.SetBool(true)

			for i := int32(0); i < length; i++ {
				key := reflect.New(v.Type().Key()).Elem()

				r.DecodeValue(world, key)
				if err != nil {
					return err
				}

				v.SetMapIndex(key, trueVal)

				if err = check(key); err != nil {
					return err
				}
			}
		} else {
			// it's a mapping

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

var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\xA0")

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

type SaveVersion uint32

var saveVersions = map[SaveVersion]string{
	1107: "0.21.93.19a",
	1108: "0.21.93.19c",
	// => "0.21.95.19a",
	// => "0.21.95.19b",
	1110: "0.21.95.19c",
	// => "0.21.100.19a",
	1113: "0.21.101.19a",
	1114: "0.21.101.19d",
	1117: "0.21.102.19a",
	1119: "0.21.104.19b",
	1121: "0.21.104.19d",
	1123: "0.21.104.21a",
	1125: "0.21.104.21b",
	1128: "0.21.105.21a",
	// => "0.22.107.21a",
	1134: "0.22.110.22e",
	1137: "0.22.110.22f",
	1148: "0.22.110.23a",
	1151: "0.22.120.23a",
	1161: "0.22.121.23b",
	1165: "0.22.123.23a",
	1169: "0.23.130.23a",

	1205: "0.27.169.32a",
	1206: "0.27.169.33a",
	1209: "0.27.169.33b",
	1211: "0.27.169.33c",
	1212: "0.27.169.33d",
	1213: "0.27.169.33e",
	1215: "0.27.169.33f",
	1216: "0.27.169.33g",
	1223: "0.27.173.38a",
	1231: "0.27.176.38a",
	1234: "0.27.176.38b",
	1235: "0.27.176.38c",
	1254: "0.28.181.39a",
	1255: "0.28.181.39b",
	1256: "0.28.181.39c",
	1259: "0.28.181.39d",
	1260: "0.28.181.39e",
	1261: "0.28.181.39f",
	1265: "0.28.181.40a",
	1266: "0.28.181.40b",
	1267: "0.28.181.40c",
	1268: "0.28.181.40d",

	1287: "0.31.01",
	1288: "0.31.02",
	1289: "0.31.03",
	1292: "0.31.04",
	1295: "0.31.05",
	1297: "0.31.06",
	1300: "0.31.08",
	1304: "0.31.09",
	1305: "0.31.10",
	1310: "0.31.11",
	1311: "0.31.12",
	1323: "0.31.13",
	1325: "0.31.14",
	1326: "0.31.15",
	1327: "0.31.16",
	1340: "0.31.17",
	1341: "0.31.18",
	1351: "0.31.19",
	1353: "0.31.20",
	1354: "0.31.21",
	1359: "0.31.22",
	1360: "0.31.23",
	1361: "0.31.24",
	1362: "0.31.25",

	1372: "0.34.01",
	1374: "0.34.02",
	1376: "0.34.03",
	1377: "0.34.04",
	1378: "0.34.05",
	1382: "0.34.06",
	1383: "0.34.07",
	1400: "0.34.08",
	1402: "0.34.09",
	1403: "0.34.10",
	1404: "0.34.11",

	1441: "0.40.01",
	1442: "0.40.02",
	1443: "0.40.03",
	1444: "0.40.04",
	1445: "0.40.05",
	1446: "0.40.06",
	1448: "0.40.07",
	1449: "0.40.08",
	1451: "0.40.09",
	1452: "0.40.10",
	1456: "0.40.11",
	1459: "0.40.12",
	1462: "0.40.13",
	1469: "0.40.14",
	1470: "0.40.15",
	1471: "0.40.16",
	1472: "0.40.17",
	1473: "0.40.18",
	1474: "0.40.19",
}

func (i SaveVersion) prettyPrint(w *WorldDat, buf, indent []byte, outerTag reflect.StructTag) []byte {
	buf = strconv.AppendInt(buf, int64(i), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(i), 16)
	buf = append(buf, ')')

	if v, ok := saveVersions[i]; ok {
		buf = append(buf, " ("...)
		buf = append(buf, v...)
		buf = append(buf, ')')
	}

	return buf
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

type Header struct {
	Version     SaveVersion
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
			r.Reader = &wtf23aReader{r: io.MultiReader(&undo, r.Reader)}
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
		r.Reader = &wtf23aReader{r: io.MultiReader(&undo, r.Reader)}
		h.Compression = Special23a
		err = binary.Read(r, binary.LittleEndian, &h.Version)
		if err != nil {
			return
		}
	}

	if _, ok := saveVersions[h.Version]; !ok {
		err = fmt.Errorf("df2014: unhandled version %d", h.Version)
		return
	}

	return
}

type compression1Reader struct {
	r   io.Reader
	buf bytes.Buffer
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
	var length int32 // signed so huge numbers cause errors instead of allocating tons of memory
	err = binary.Read(r.r, binary.LittleEndian, &length)
	if err != nil {
		return
	}
	if length < 0 {
		return fmt.Errorf("df2014: negative length (%d)", length)
	}

	z, err := zlib.NewReader(io.LimitReader(r.r, int64(length)))
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
	return
}
