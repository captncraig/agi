package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/captncraig/agi/directory"
	"github.com/captncraig/agi/picture"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dir, err := directory.Load("PICDIR")
	if err != nil {
		log.Fatal(err)
	}
	for i, pic := range dir {
		if pic == nil {
			continue
		}
		outFile := filepath.Join("pics", fmt.Sprintf("%d.png", i))
		logFile := filepath.Join("picdasm", fmt.Sprintf("%d.log", i))
		out2File := filepath.Join("pics", fmt.Sprintf("%d.priority.png", i))
		pic, pri := picture.DecodePicture(pic.Data, logFile)
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal(err)
		}
		err = png.Encode(f, pic.Image())
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		f, err = os.Create(out2File)
		if err != nil {
			log.Fatal(err)
		}
		err = png.Encode(f, pri.Image())
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}
