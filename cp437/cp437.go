// Package cp437 provides conversions between IBM code page 437 and UTF-8.
package cp437

var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\u00A0")

var reverse = func(r []rune) map[rune]byte {
	m := make(map[rune]byte)
	for i, c := range r {
		m[c] = byte(i)
	}
	return m
}(cp437)

// Rune takes a CP437 byte and returns a Unicode code point.
func Rune(b byte) rune {
	return cp437[b]
}

// String takes a CP437-encoded byte slice and returns a UTF-8-encoded string.
func String(b []byte) string {
	r := make([]rune, len(b))

	for i := range r {
		r[i] = Rune(b[i])
	}

	return string(r)
}

// Byte takes a Unicode code point and returns the corresponding CP437 byte,
// or '?' if there is no corresponding CP437 byte.
func Byte(r rune) byte {
	if b, ok := reverse[r]; ok {
		return b
	}
	return '?'
}

// Bytes takes a UTF-8-encoded string and returns a CP437-encoded byte slice.
func Bytes(s string) []byte {
	r := []rune(s)
	b := make([]byte, len(r))

	for i := range b {
		b[i] = Byte(r[i])
	}

	return b
}
