package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/captncraig/agi/logic"
	"github.com/captncraig/agi/picture"

	"github.com/captncraig/agi/directory"
)

func main() {
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
			lo := logic.Decode(l.Data, i)
			fmt.Println("LOGIC", i)
			// for j := byte(0); j < 255; j++ {
			// 	s, ok := lo.Text[j]
			// 	if !ok {
			// 		continue
			// 	}
			// 	d, _ := json.Marshal(s)
			// 	fmt.Println("\t\t", j, string(d))
			// }
			for j := 0; j < 15 && j < len(lo.Instructions); j++ {
				fmt.Printf("0x%x\n", lo.Instructions[j])
			}
		}
	}
}
