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

type Reader struct {
	io.Reader
}

func (r *Reader) Decode(v interface{}) error {
	return r.DecodeValue(reflect.ValueOf(v).Elem())
}

func (r *Reader) DecodeValue(v reflect.Value) error {
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
			if tag := v.Type().Field(i).Tag.Get("df2014_get_length_from"); tag != "" {
				l := v.FieldByName(tag).Len()
				v.Field(i).Set(reflect.MakeSlice(v.Type().Field(i).Type, l, l))

				for j := 0; j < l; j++ {
					err := r.DecodeValue(v.Field(i).Index(j))
					if err != nil {
						v.Field(i).SetLen(j + 1) // hide the entries that we didn't get to
						return err
					}
				}
			} else {
				err := r.DecodeValue(v.Field(i))
				if err != nil {
					return err
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
		}
		return nil

	case reflect.Ptr:
		flag, err := r.bool()
		if err != nil {
			return err
		}
		if flag {
			v.Set(reflect.New(v.Type().Elem()))
			return r.DecodeValue(v.Elem())
		} else {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}

	case reflect.String:
		s, err := r.string()
		if err != nil {
			return err
		}
		v.SetString(s)
		return nil

	case reflect.Slice:
		var length uint32
		err := binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return err
		}

		v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))

		fallthrough
	case reflect.Array:
		for i, l := 0, v.Len(); i < l; i++ {
			err := r.DecodeValue(v.Index(i))
			if err != nil {
				return err
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
		var length uint32
		err := binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return err
		}

		v.Set(reflect.MakeMap(v.Type()))

		var prev reflect.Value
		check := func(next reflect.Value) error {
			if prev.IsValid() {
				switch prev.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if prev.Int() > next.Int() {
						return fmt.Errorf("df2014: values not in order: %v > %v", prev, next)
					}
					if prev.Int() == next.Int() {
						return fmt.Errorf("df2014: duplicate value: %v", prev)
					}

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if prev.Uint() > next.Uint() {
						return fmt.Errorf("df2014: values not in order: %v > %v", prev, next)
					}
					if prev.Uint() == next.Uint() {
						return fmt.Errorf("df2014: duplicate value: %v", prev)
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

			for i := uint32(0); i < length; i++ {
				key := reflect.New(v.Type().Key()).Elem()

				r.DecodeValue(key)
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

			for i := uint32(0); i < length; i++ {
				key := reflect.New(v.Type().Key()).Elem()
				val := reflect.New(v.Type().Elem()).Elem()

				err = r.DecodeValue(key)
				if err != nil {
					return err
				}
				err = r.DecodeValue(val)
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
	err = r.Decode(&n)
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
	var length uint16
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	s := make([]rune, length)
	for i, b := range buf {
		s[i] = cp437[b]
	}

	return string(s), nil
}

type Header struct {
	Version, Compression uint32
}

func (r *Reader) header() (h Header, err error) {
	err = binary.Read(r, binary.LittleEndian, &h.Version)
	if err != nil {
		return
	}

	if h.Version != 1456 {
		err = fmt.Errorf("df2014: unhandled version %d", h.Version)
		return
	}

	err = binary.Read(r, binary.LittleEndian, &h.Compression)
	if err != nil {
		return
	}

	switch h.Compression {
	case 0:
		// nothing to be done
	case 1:
		r.Reader = &compression1Reader{r: r.Reader}
	default:
		err = fmt.Errorf("df2014: unhandled compression type %d", h.Compression)
		return
	}

	return
}

type compression1Reader struct {
	r   io.Reader
	buf []byte
}

func (r *compression1Reader) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 {
		err = r.fill()
		if err != nil {
			return
		}
	}

	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return
}

func (r *compression1Reader) fill() (err error) {
	var length uint32
	err = binary.Read(r.r, binary.LittleEndian, &length)
	if err != nil {
		return
	}

	var buf bytes.Buffer
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

	_, err = io.Copy(&buf, z)
	if err != nil {
		return
	}

	r.buf = buf.Bytes()
	return
}
