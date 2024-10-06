package util

import "log"

const debug = true

func Log(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}
