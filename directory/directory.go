// Package directory handles loading of a resource index and reading associated data from VOL files
// It will load ALL data into memory immediately, assuming any modern system has ample memory availible.
package directory

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ResourcePointer struct {
	Offset    uint32
	Signature uint16
	VolNumber byte
	Length    uint16
	Data      []byte
}

type Directory struct {
	Vols                           [][]byte
	Pictures, Logic, Views, Sounds []*ResourcePointer
}

func (d *Directory) NumResources(rps []*ResourcePointer) int {
	count := 0
	for _, p := range rps {
		if p != nil {
			count++
		}
	}
	return count
}

func New(dir string) (*Directory, error) {
	d := &Directory{}
	for i := 0; ; i++ {
		fname := fmt.Sprintf("VOL.%d", i)
		dat, err := ioutil.ReadFile(filepath.Join(dir, fname))
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return nil, err
		}
		d.Vols = append(d.Vols, dat)
	}
	for _, file := range []string{"LOGDIR", "PICDIR", "VIEWDIR"} {
		rps, err := d.load(filepath.Join(dir, file))
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %s", file, err)
		}
		switch file[:3] {
		case "LOG":
			d.Logic = rps
		case "PIC":
			d.Pictures = rps
		case "SND":
			d.Sounds = rps
		case "VIE":
			d.Views = rps
		}
	}
	return d, nil
}

// Load will read a directory file, and load all records it finds into memory from the referenced VOL files.
// If no error is returned, the result will always contain 256 entries, but some may be nil of the resource does not exist
func (d *Directory) load(filename string) ([]*ResourcePointer, error) {
	resources := make([]*ResourcePointer, 256)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	numRecs := len(dat) / 3
	for i := 0; i < numRecs; i++ {
		idx := i * 3
		vol := int(dat[idx]>>4) & 0x0f
		offset := uint32(dat[idx]&0xf)<<16 | uint32(dat[idx+1])<<8 | uint32(dat[idx+2])
		if offset == 0x0fffff && vol == 0x0f {
			continue
		}
		if vol >= len(d.Vols) {
			return nil, fmt.Errorf("Volume %d not loaded", vol)
		}
		v := d.Vols[vol]
		if len(v) < int(offset+5) {
			log.Printf("VOLUME %d not long enough to read offset 0x%x", vol, offset)
			continue
		}
		rec := &ResourcePointer{
			Offset:    offset,
			Signature: uint16(v[offset])<<8 | uint16(v[offset+1]),
			VolNumber: v[offset+2],
			Length:    uint16(v[offset+3]) | uint16(v[offset+4])<<8,
		}
		if rec.Signature != 0x1234 {
			log.Fatal("Sig No Match")
		}
		if len(v) < int(offset+5)+int(rec.Length) {
			return nil, fmt.Errorf("Remaining size less than required for chunk")
		}
		rec.Data = v[offset+5 : offset+5+uint32(rec.Length)]
		resources[i] = rec
	}
	return resources, err
}
