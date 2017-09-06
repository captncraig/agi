package logic

type Logic struct {
	Text         map[byte]string
	Instructions []byte
}

func Decode(d []byte, idx int) *Logic {
	r16 := func(i int) uint16 {
		return uint16(d[i]) | uint16(d[i+1])<<8
	}
	instructions := []byte{}

	startPos := int(r16(0)) + 2
	for i := 2; i < startPos; i++ {
		instructions = append(instructions, d[i])
	}
	lo := &Logic{
		Instructions: instructions,
		Text: decodeStrings(d[startPos:], idx),
	}
	return lo
}

func decodeStrings(d []byte, idx int) map[byte]string {
	r16 := func(i int) uint16 {
		return uint16(d[i]) | uint16(d[i+1])<<8
	}
	n := int(d[0])
	offsets := make([]int, n)
	for i := range offsets {
		offsets[i] = int(r16(3 + 2*i))
	}
	d = d[3+2*n:]
	key := "Avis Durgan"
	for i := range d {
		d[i] ^= key[i%len(key)]
	}
	m := map[byte]string{}
	for i, off := range offsets {
		if off == 0 {
			continue
		}
		off -= (2 + 2*n)
		str := ""
		for j := off; j < len(d); j++ {
			c := d[j]
			if c == 0 {
				break
			}
			str += string(c)
		}
		m[byte(i)] = str
	}
	return m
}
