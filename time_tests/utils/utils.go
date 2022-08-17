package utils

import (
	"log"
	"os"
)

func ParseArgs() (pattern, file string) {
	if len(os.Args) < 3 {
		log.Fatal("usage: petergrep <pattern> <file_name>")
	}
	pattern = os.Args[1]
	file = os.Args[2]
	return
}
