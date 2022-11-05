package main

import "errors"

var (
	ErrReadDir  = errors.New("failed to read directory")
	ErrReadFile = errors.New("failed to read file")
)

const (
	BYTES_TB = BYTES_GB * 1000
	BYTES_GB = BYTES_MB * 1000
	BYTES_MB = BYTES_KB * 1000
	BYTES_KB = 1000
)
