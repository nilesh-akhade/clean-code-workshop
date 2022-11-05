package main

import (
	"io/ioutil"
	"os"
)

type EntryHandler interface {
	Handle(*DuplicateDetails) error
}

type FileEntry struct {
	fullPath string
	size     int64
}

func (fileEntry *FileEntry) Handle(dupDetails *DuplicateDetails) error {
	file, err := ioutil.ReadFile(fileEntry.fullPath)
	if err != nil {
		return ErrReadFile
	}
	hashString := generateHash(file)
	dupDetails.AddEntry(hashString, fileEntry)
	return nil
}

type DirEntry struct {
	fullPath string
}

func (dirEntry *DirEntry) Handle(dupDetails *DuplicateDetails) error {
	err := dupDetails.traverseDir(dirEntry.fullPath)
	if err != nil {
		return err
	}
	return nil
}

type NilEntry struct{}

func (NilEntry) Handle(dupDetails *DuplicateDetails) error {
	return nil
}

func NewEntryHandler(entry os.FileInfo, fullPath string) EntryHandler {
	switch {
	case entry.IsDir():
		return &DirEntry{
			fullPath: fullPath,
		}
	case entry.Mode().IsRegular():
		return &FileEntry{
			fullPath: fullPath,
			size:     entry.Size(),
		}
	default:
		return &NilEntry{}
	}
}
