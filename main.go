package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/captncraig/agi/logic"
	"github.com/captncraig/agi/picture"

	"github.com/captncraig/agi/directory"
)

func main() {
	log.SetFlags(log.Lshortfile)
	name := "sq1"
	d, err := directory.New(filepath.Join("games", name))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d VOL files", len(d.Vols))
	log.Printf("%d Pictures", d.NumResources(d.Pictures))
	log.Printf("%d Logic", d.NumResources(d.Logic))
	log.Printf("%d Views", d.NumResources(d.Views))
	log.Printf("%d Sounds", d.NumResources(d.Sounds))
	imgOutDir := filepath.Join("docs", name, "pictures")
	os.MkdirAll(imgOutDir, 0644)
	for i, p := range d.Pictures {
		break
		if p != nil {
			pic, err := picture.Decode(p.Data)
			if err != nil {
				log.Fatal(i, err)
			}
			if err = pic.RenderToFile(imgOutDir, i, false); err != nil {
				log.Fatal(err)
			}
			if err = pic.RenderToFile(imgOutDir, i, true); err != nil {
				log.Fatal(err)
			}
			pic.RenderGifs(imgOutDir, i)
		}
	}
	for i, l := range d.Logic {
		if l != nil {
			fmt.Println("LOGIC", i)
			lo := logic.Decode(l.Data, i)
			fmt.Println(lo.DASM)
			ioutil.WriteFile("dasm", []byte(lo.DASM),0666)
			break
			// for j := byte(0); j < 255; j++ {
			// 	s, ok := lo.Text[j]
			// 	if !ok {
			// 		continue
			// 	}
			// 	d, _ := json.Marshal(s)
			// 	fmt.Println("\t\t", j, string(d))
			// }
			// for j := 0; j < 15 && j < len(lo.Instructions); j++ {
			// 	fmt.Printf("0x%x\n", lo.Instructions[j])
			// }
		}
	}

	// dat, err := ioutil.ReadFile(filepath.Join("games", "sq1", "AGIDATA.OVL"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // search for possible header blocks for agi commands. n is number of commands to check from AgiCommands
	// const n = 7
	// tableStart := 0
	// for i := 0; i < len(dat)-(4*n); i++ {
	// 	ok := true
	// 	for j := 0; j < n; j++ {
	// 		baseIdx := i + (4 * j)
	// 		if int(dat[baseIdx+2]) != len(logic.TestCommands[j].ArgTypes) {
	// 			ok = false
	// 			break
	// 		}
	// 	}
	// 	if ok {
	// 		tableStart = i
	// 		break
	// 	}
	// }
	// r := binutils.NewReader(dat[tableStart:])
	// for i := 0; i < len(logic.TestCommands) && r.HasMore(); i++ {
	// 	fmt.Printf("%02x: %04x n=%d %02x\n", i, r.Take16(), r.Take(), r.Take())
	// }
}
