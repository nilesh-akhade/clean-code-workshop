package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

func traverseDir(hashes, duplicates map[string]string, dupeSize *int64, entries []os.FileInfo, directory string) {
	for _, entry := range entries {
		fullpath := (path.Join(directory, entry.Name()))

		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}

		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				panic(err)
			}
			traverseDir(hashes, duplicates, dupeSize, dirFiles, fullpath)
			continue
		}
		file, err := ioutil.ReadFile(fullpath)
		if err != nil {
			panic(err)
		}
		hashString := generateHash(file)
		if hashEntry, ok := hashes[hashString]; ok {
			duplicates[hashEntry] = fullpath
			atomic.AddInt64(dupeSize, entry.Size())
		} else {
			hashes[hashString] = fullpath
		}
	}
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
			panic(err)
		}
	}

	hashes := map[string]string{}
	duplicates := map[string]string{}
	var dupeSize int64

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		panic(err)
	}

	traverseDir(hashes, duplicates, &dupeSize, entries, *dir)

	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(hashes))
	fmt.Println("DUPLICATES:", len(duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(dupeSize))
}

// running into problems of not being able to open directories inside .app folders
