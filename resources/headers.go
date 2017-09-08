package resources

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/captncraig/agi"
)

type ResourceHeader struct {
	ID         byte
	FileNumber byte
	Offset     uint32

	LoadError error
	Raw       []byte
}

type View struct {
	ResourceHeader
}

type Picture struct {
	ResourceHeader
}

type Sound struct {
	ResourceHeader
}

type Directory struct {
	Logics   map[byte]*Logic
	Pictures map[byte]*Picture
	Views    map[byte]*Picture
	Sounds   map[byte]*Sound
}

func Load(l agi.Loader) (*Directory, error) {
	dirFiles, err := l.ListDIRs()
	if err != nil {
		return nil, err
	}
	d := &Directory{
		Logics:   map[byte]*Logic{},
		Pictures: map[byte]*Picture{},
		Views:    map[byte]*Picture{},
		Sounds:   map[byte]*Sound{},
	}
	var headerLists [][]*ResourceHeader
	if len(dirFiles) == 1 {
		headerLists, err = d.loadV3(l, dirFiles[0])
	} else {
		headerLists, err = d.loadV2(l)
	}
	if err != nil {
		return nil, err
	}
	for _, logic := range headerLists[0] {
		if logic.LoadError != nil {
			fmt.Printf("LOGIC %d: %s\n", logic.ID, logic.LoadError)
			continue
		}
		d.Logics[logic.ID] = DecodeLogic(logic.Raw)
	}
	for _, pic := range headerLists[1] {
		if pic.LoadError != nil {
			fmt.Printf("PICTURE %d: %s\n", pic.ID, pic.LoadError)
			continue
		}
	}
	for _, view := range headerLists[2] {
		if view.LoadError != nil {
			fmt.Printf("View %d: %s\n", view.ID, view.LoadError)
			continue
		}
	}
	for _, sound := range headerLists[3] {
		if sound.LoadError != nil {
			fmt.Printf("SOUND %d: %s\n", sound.ID, sound.LoadError)
			continue
		}
	}
	return d, nil
}

func (d *Directory) loadV2(l agi.Loader) ([][]*ResourceHeader, error) {
	headerLists := [][]*ResourceHeader{}
	for _, fname := range []string{"LOGDIR", "PICDIR", "VIEWDIR", "SNDDIR"} {
		f, err := l.GetFile(fname)
		if err != nil {
			return nil, err
		}
		hds, err := readV2Headers(l, f)
		if err != nil {
			return nil, err
		}
		headerLists = append(headerLists, hds)
	}
	return headerLists, nil

}

func readV2Headers(l agi.Loader, dat []byte) ([]*ResourceHeader, error) {
	if len(dat)%3 != 0 {
		return nil, fmt.Errorf("V2 DIR file should be multiple of three bytes")
	}
	if len(dat) > 256*3 {
		return nil, fmt.Errorf("V2 DIR file cannot have more than 256 entries")
	}
	hds := []*ResourceHeader{}
	for i := 0; i < len(dat)/3; i++ {
		hd := &ResourceHeader{
			ID:         byte(i),
			Offset:     uint32(dat[i*3]&0x0f)<<16 | uint32(read16BE(dat, i*3+1)), // last 4 bits of first byte, and second two bytes
			FileNumber: dat[i*3] >> 4,
		}
		if hd.Offset != 0x0fffff || hd.FileNumber != 0x0f {
			dat, err := readV2Data(l, hd)
			if err != nil {
				hd.LoadError = err
			}
			hd.Raw = dat
			hds = append(hds, hd)
		}
	}
	return hds, nil
}

func readV2Data(l agi.Loader, hd *ResourceHeader) ([]byte, error) {
	fname := fmt.Sprintf("VOL.%d", hd.FileNumber)
	dat, err := l.GetFile(fname)
	if err != nil {
		return nil, err
	}
	if len(dat) < int(hd.Offset)+5 {
		return nil, fmt.Errorf("%s is not long enough to seek offset 0x%x", fname, hd.Offset)
	}
	dat = dat[hd.Offset:]
	// every header starts with 0x1234
	if read16BE(dat, 0) != 0x1234 {
		return nil, fmt.Errorf("Resource signature not valid")
	}
	//dat[2] is volume number, don't care
	length := read16(dat, 3)
	dat = dat[5:]
	if len(dat) < int(length) {
		return nil, fmt.Errorf("%s is not long enough for resource with length 0x%x", fname, length)
	}
	return dat[:length], nil
}

func (d *Directory) loadV3(l agi.Loader, prefix string) ([][]*ResourceHeader, error) {
	headerLists := [][]*ResourceHeader{}
	fname := fmt.Sprintf("%sDIR", prefix)
	dat, err := l.GetFile(fname)
	if err != nil {
		return nil, err
	}
	typeNames := []string{"logic", "picture", "view", "sound"} // for errors
	// first 8 bytes are LE pointers to resource range for logic,picture,view,sound types
	for i := 0; i < 4; i++ {
		start := int(read16(dat, i*2))
		// end is first byte outside the range
		end := len(dat)
		if i < 3 {
			end = int(read16(dat, 2+i*2))
		}
		if (end-start)%3 != 0 {
			return nil, fmt.Errorf("resource header block for %s must be multiple of three bytes", typeNames[i])
		}
		hds := []*ResourceHeader{}
		for i2 := 0; i2 < int((end-start))/3; i2++ {
			s := start + (i2 * 3)
			hd := &ResourceHeader{
				ID:         byte(i2),
				Offset:     uint32(dat[s]&0x0f)<<16 | uint32(read16BE(dat, s+1)),
				FileNumber: dat[s] >> 4,
			}
			if hd.Offset != 0x0fffff || hd.FileNumber != 0x0f {
				dat, err := readV3Data(l, prefix, hd)
				if err != nil {
					fmt.Printf("ERROR LOADING %s %d %x: %s\n", typeNames[i], hd.ID, hd.Offset, err)
					hd.LoadError = err
				}
				hd.Raw = dat
				hds = append(hds, hd)
			}
		}
		headerLists = append(headerLists, hds)
	}
	return headerLists, nil
}

func readV3Data(l agi.Loader, prefix string, hd *ResourceHeader) ([]byte, error) {
	fname := fmt.Sprintf("%sVOL.%d", prefix, hd.FileNumber)
	dat, err := l.GetFile(fname)
	if err != nil {
		return nil, err
	}
	if len(dat) < int(hd.Offset)+7 {
		return nil, fmt.Errorf("%s is not long enough to seek offset 0x%x", fname, hd.Offset)
	}
	dat = dat[hd.Offset:]
	// every header starts with 0x1234
	if read16BE(dat, 0) != 0x1234 {
		return nil, fmt.Errorf("Resource signature not valid")
	}
	isPic := dat[2]>>7 == 1
	uncompressedLength := int(read16(dat, 3))
	compressedLength := int(read16(dat, 5))
	dat = dat[7:]
	if len(dat) < compressedLength {
		return nil, fmt.Errorf("%s is not long enough for resource with length 0x%x", fname, compressedLength)
	}
	rawDat := dat[:compressedLength]
	if uncompressedLength != compressedLength {
		if isPic {
			log.Fatal("PIC EXPAND")
		} else {
			log.Fatal("LZW EXPAND")
		}
		fmt.Println(len(rawDat), compressedLength, uncompressedLength)
	}
	return nil, nil
}

// binary helper functions

// little endian 16 bit read
func read16(dat []byte, i int) uint16 {
	return binary.LittleEndian.Uint16(dat[i:])
}

func read16BE(dat []byte, i int) uint16 {
	return binary.BigEndian.Uint16(dat[i:])
}
