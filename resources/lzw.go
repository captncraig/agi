package resources

// import "fmt"

// //adapted from https://github.com/scummvm/scummvm/blob/master/engines/agi/lzw.cpp

// func lzwExpand(dat []byte, expectedLength int) []byte {
// 	// class in a function to contain things
// 	const TableSize = 18041
// 	const MaxBits uint16 = 12
// 	//initialization of class level things: #L69
// 	decodeStack := []byte{}
// 	//prefixCodes := make([]uint32, TableSize)
// 	//appendCharacter := make([]byte, TableSize)
// 	var bits uint16
// 	var maxValue uint32
// 	var maxCode uint32

// 	curIdx := 0
// 	curBit := byte(7)
// 	readBit := func() uint16 {
// 		b := (dat[curIdx] >> curBit) & 1
// 		curBit--
// 		if curBit == 0xff {
// 			curBit = 7
// 			curIdx++
// 		}
// 		return uint16(b)
// 	}
// 	readCode := func() uint16 {
// 		d := uint16(0)
// 		for i := uint16(0); i < bits; i++ {
// 			d <<= 1
// 			d |= readBit()
// 		}
// 		return d
// 	}
// 	setBits := func(value uint16) {
// 		bits = value
// 		maxValue = (1 << bits) - 1
// 		maxCode = maxValue - 1
// 	}
// 	decodeString := func(code uint16) []byte {
// 		for i := 0; code > 255; {

// 		}
// 	}

// 	// finally expand function:
// 	var c, lzwnext, lzwnew, lzwold uint16
// 	//var s []byte
// 	out := []byte{}
// 	setBits(9)
// 	lzwnext = 257
// 	lzwold = readCode()
// 	c = lzwold
// 	lzwnew = readCode()
// 	count := 0
// 	countMax := 1
// 	for len(out) < expectedLength && lzwnew != 0x101 {
// 		fmt.Println("old", lzwold, "new", lzwnew)
// 		if count >= countMax {
// 			break
// 		}
// 		count++
// 		if lzwnew == 0x100 {
// 			// Code to "start over"
// 			lzwnext = 258
// 			setBits(9)
// 			lzwold = readCode()
// 			c = lzwold
// 			out = append(out, byte(c))
// 			lzwnew = readCode()
// 		} else {
// 			if lzwnew >= lzwnext {
// 				// Handles special LZW scenario
// 				decodeStack = append(decodeStack, byte(c))
// 				s = decodeString(decodeStack, lzwold)
// 			} else {
// 				s = decodeString(decodeStack, lzwnew)
// 			}
// 		}
// 	}

// 	return nil
// }
