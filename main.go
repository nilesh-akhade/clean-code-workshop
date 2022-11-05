package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var err error
	dir := flag.String("path", "", "the path to traverse searching for duplicates")
	flag.Parse()

	if *dir == "" {
		*dir, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error traversing directory: %v", err)
			return
		}
	}

	var dupeSize int64
	dupeDetails := &DuplicateDetails{
		Hashes:     map[string]string{},
		Duplicates: map[string]string{},
		DupeSize:   &dupeSize,
	}

	err = dupeDetails.traverseDir(*dir)
	if err != nil {
		fmt.Printf("Error traversing directory: %v", err)
		return
	}

	dupeDetails.printResult()
}

// running into problems of not being able to open directories inside .app folders
