package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync/atomic"
)

type DuplicateDetails struct {
	Hashes     map[string]string
	Duplicates map[string]string
	DupeSize   *int64
}

func (dupDetails *DuplicateDetails) AddEntry(hashString string, fileEntry *FileEntry) {
	if hashEntry, ok := dupDetails.Hashes[hashString]; ok {
		dupDetails.Duplicates[hashEntry] = fileEntry.fullPath
		atomic.AddInt64(dupDetails.DupeSize, fileEntry.size)
	} else {
		dupDetails.Hashes[hashString] = fileEntry.fullPath
	}
}

func (dupDetails *DuplicateDetails) traverseDir(directory string) error {

	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("Error traversing directory: %v", err)
	}

	for _, entry := range entries {
		fullPath := path.Join(directory, entry.Name())

		entryObject := NewEntryHandler(entry, fullPath)
		err := entryObject.Handle(dupDetails)
		if err != nil {
			return err
		}

	}
	return nil
}

func (dupDetails *DuplicateDetails) printResult() {
	fmt.Println("DUPLICATES")
	fmt.Println("TOTAL FILES:", len(dupDetails.Hashes))
	fmt.Println("DUPLICATES:", len(dupDetails.Duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(*dupDetails.DupeSize))
}
