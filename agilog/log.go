package agilog

import "log"

var Debug = true

func Debugf(f string, args ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(f, args...)
}
