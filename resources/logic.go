package resources

type Logic struct {
	Messages     map[byte]string
	Instructions []byte
}

func DecodeLogic(dat []byte) *Logic {
	//todo: bounds checking
	l := &Logic{}
	textStart := int(read16(dat, 0)) + 2
	l.Instructions = dat[2 : 2+textStart-2]
	l.Messages = decodeMessages(dat[textStart:])
	return l
}

func decodeMessages(d []byte) map[byte]string {
	n := int(d[0])
	offsets := make([]int, n)
	for i := range offsets {
		offsets[i] = int(read16(d, 3+2*i))
	}
	d = d[3+2*n:]
	key := "Avis Durgan"
	// decode all
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
