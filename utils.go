package main

import (
	"crypto/sha1"
	"fmt"
	"strconv"
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

func generateHash(bytes []byte) string {
	hash := sha1.New()
	if _, err := hash.Write(bytes); err != nil {
		return fmt.Sprint(err)
	}
	hashSum := hash.Sum(nil)
	return fmt.Sprintf("%x", hashSum)
}
