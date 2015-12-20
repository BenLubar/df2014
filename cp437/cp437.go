package cp437

var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\u00A0")

var reverse = func(r []rune) map[rune]byte {
	m := make(map[rune]byte)
	for i, c := range r {
		m[c] = byte(i)
	}
	return m
}(cp437)

func Rune(b byte) rune {
	return cp437[b]
}

func String(b []byte) string {
	r := make([]rune, len(b))

	for i := range r {
		r[i] = Rune(b[i])
	}

	return string(r)
}

func Byte(r rune) byte {
	if b, ok := reverse[r]; ok {
		return b
	}
	return '?'
}

func Bytes(s string) []byte {
	r := []rune(s)
	b := make([]byte, len(r))

	for i := range b {
		b[i] = Byte(r[i])
	}

	return b
}
