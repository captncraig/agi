package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/captncraig/agi"
	"github.com/captncraig/agi/dasm"
	"github.com/captncraig/agi/resources"
)

var dir = flag.String("d", "games/sq1", "directory containing resource files")
var print = flag.Bool("p", false, "print a bunch of stuff")

func main() {
	flag.Parse()
	loader := agi.NewFromDir(*dir)
	rez, err := resources.Load(loader)
	if err != nil {
		log.Fatal(err)
	}
	if !*print {
		return
	}
	for _, l := range rez.Logics {
		fmt.Println(dasm.Decompile(l.Instructions))
		break
	}
}
