// Package directory handles loading of a resource index and reading associated data from VOL files
// It will load ALL data into memory immediately, assuming any modern system has ample memory availible.
package directory

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/captncraig/agi/agilog"
)

type ResourceFile struct {
	Offset    uint32
	Signature uint16
	VolNumber byte
	Length    uint16
	Data      []byte
}

// Load will read a directory file, and load all records it finds into memory from the referenced VOL files.
// If no error is returned, the result will always contain 256 entries, but some may be nil of the resource does not exist
func Load(filename string) ([]*ResourceFile, error) {
	resources := make([]*ResourceFile, 256)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	vols := [][]byte{}
	numRecs := len(dat) / 3
	for i := 0; i < numRecs; i++ {
		idx := i * 3
		vol := int(dat[idx]>>4) & 0x0f
		offset := uint32(dat[idx]&0xf)<<16 | uint32(dat[idx+1])<<8 | uint32(dat[idx+2])
		if offset == 0x0fffff && vol == 0x0f {
			continue
		}
		agilog.Debugf("record %d VOL.%d@0x%x", i, vol, offset)
		for len(vols) < vol+1 {
			vols = append(vols, nil)
		}
		if vols[vol] == nil {
			fname := fmt.Sprintf("VOL.%d", vol)
			agilog.Debugf("LOADING %s", fname)
			dat, err := ioutil.ReadFile(fname)
			if err != nil {
				return nil, err
			}
			agilog.Debugf("%d bytes read", len(dat))
			vols[vol] = dat
		}
		v := vols[vol]
		if len(v) < int(offset+5) {
			log.Printf("VOLUME %d not long enough to read offset 0x%x", vol, offset)
			continue
		}
		rec := &ResourceFile{
			Offset:    offset,
			Signature: uint16(v[offset])<<8 | uint16(v[offset+1]),
			VolNumber: v[offset+2],
			Length:    uint16(v[offset+3]) | uint16(v[offset+4])<<8,
		}
		if rec.Signature != 0x1234 {
			log.Fatal("Sig No Match")
		}
		if len(v) < int(offset+5)+int(rec.Length) {
			// Length is larger than remaining data?
			// Take the rest of the file, and hope it is valid
			log.Println("EXTRA DATA!!!?!?!?")
			rec.Data = v[offset+5:]
		} else {
			rec.Data = v[offset+5 : offset+5+uint32(rec.Length)]
		}
		resources[i] = rec
		agilog.Debugf("Signature: 0x%x Vol: %d Length %d(%d read)\n\n", rec.Signature, rec.VolNumber, rec.Length, len(rec.Data))
	}
	return resources, err
}
