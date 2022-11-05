package main

import (
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

type DuplicateDetails struct {
	Hashes     map[string]string
	Duplicates map[string]string
	DupeSize   *int64
}

var (
	ErrReadDir  = errors.New("failed to read directory")
	ErrReadFile = errors.New("failed to read file")
)

func traverseDir(dupDetails *DuplicateDetails, entries []os.FileInfo, directory string) error {
	for _, entry := range entries {
		fullpath := (path.Join(directory, entry.Name()))

		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}

		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				return ErrReadDir
			}
			err = traverseDir(dupDetails, dirFiles, fullpath)
			if err != nil {
				return err
			}
			continue
		}
		file, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return ErrReadFile
		}
		hashString := generateHash(file)
		if hashEntry, ok := dupDetails.Hashes[hashString]; ok {
			dupDetails.Duplicates[hashEntry] = fullpath
			atomic.AddInt64(dupDetails.DupeSize, entry.Size())
		} else {
			dupDetails.Hashes[hashString] = fullpath
		}
	}
	return nil
}

func generateHash(bytes []byte) string {
	hash := sha1.New()
	if _, err := hash.Write(bytes); err != nil {
		panic(err)
	}
	hashSum := hash.Sum(nil)
	return fmt.Sprintf("%x", hashSum)
}

const (
	BYTES_TB = BYTES_GB * 1000
	BYTES_GB = BYTES_MB * 1000
	BYTES_MB = BYTES_KB * 1000
	BYTES_KB = 1000
)

func toReadableSize(nbytes int64) string {
	switch {
	case nbytes > BYTES_TB:
		return strconv.FormatInt(nbytes/(BYTES_TB), 10) + " TB"

	case nbytes > BYTES_GB:
		return strconv.FormatInt(nbytes/(BYTES_GB), 10) + " GB"

	case nbytes > BYTES_MB:
		return strconv.FormatInt(nbytes/(BYTES_MB), 10) + " MB"

	case nbytes > BYTES_KB:
		return strconv.FormatInt(nbytes/BYTES_KB, 10) + " KB"
	default:
		return strconv.FormatInt(nbytes, 10) + " B"
	}
}

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

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		fmt.Printf("Error traversing directory: %v", err)
		return
	}
	dupeDetails := &DuplicateDetails{
		Hashes:     map[string]string{},
		Duplicates: map[string]string{},
		DupeSize:   &dupeSize,
	}

	err = traverseDir(dupeDetails, entries, *dir)
	if err != nil {
		fmt.Printf("Error traversing directory: %v", err)
		return
	}

	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(dupeDetails.Hashes))
	fmt.Println("DUPLICATES:", len(dupeDetails.Duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(dupeSize))
}

// running into problems of not being able to open directories inside .app folders
