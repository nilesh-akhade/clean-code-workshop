package main

import (
	"io/ioutil"
	"testing"
)

func TestTraverseDir(t *testing.T) {

	tests := []struct {
		Name      string
		Directory string
		Want      struct {
			Total     int
			Duplicate int
			DupSize   int64
		}
	}{
		{
			Name:      "No duplicates",
			Directory: "./testdata/testno-dupes",
			Want: struct {
				Total     int
				Duplicate int
				DupSize   int64
			}{
				Total:     2,
				Duplicate: 0,
				DupSize:   0,
			},
		},
		{
			Name:      "4 duplicates",
			Directory: "./testdata/some-duplicates",
			Want: struct {
				Total     int
				Duplicate int
				DupSize   int64
			}{
				Total:     4,
				Duplicate: 1,
				DupSize:   11,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			entries, err := ioutil.ReadDir(tt.Directory)
			if err != nil {
				panic(err)
			}
			var dupeSize int64
			hashes := make(map[string]string)
			duplicates := make(map[string]string)
			traverseDir(hashes, duplicates, &dupeSize, entries, tt.Directory)
			if dupeSize != tt.Want.DupSize {
				t.Errorf("size: expected:%v, got: %v", tt.Want.DupSize, dupeSize)
			}
			if len(hashes) != tt.Want.Total {
				t.Errorf("total: expected:%v, got: %v", tt.Want.Total, len(hashes))
			}
			if len(duplicates) != tt.Want.Duplicate {
				t.Errorf("duplicates: expected:%v, got: %v", tt.Want.Duplicate, len(duplicates))
			}
		})
	}
}

func TestTraverseDir_NoDir(t *testing.T) {

	tests := []struct {
		Name      string
		Directory string
	}{
		{
			Name:      "Not found",
			Directory: "./testdata/notfound-dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			entries, err := ioutil.ReadDir("./testdata/testno-dupes")
			if err != nil {
				panic(err)
			}
			var dupeSize int64
			hashes := make(map[string]string)
			duplicates := make(map[string]string)
			traverseDir(hashes, duplicates, &dupeSize, entries, tt.Directory)
		})
	}
}

func TestToReadableSize(t *testing.T) {
	tests := []struct {
		Name       string
		InBytes    int64
		InReadable string
	}{
		{
			Name:       "Bytes conversion",
			InBytes:    100,
			InReadable: "100 B",
		}, {
			Name:       "KB conversion",
			InBytes:    1001,
			InReadable: "1 KB",
		}, {
			Name:       "MB conversion",
			InBytes:    2005 * 1000,
			InReadable: "2 MB",
		}, {
			Name:       "GB conversion",
			InBytes:    9005 * 1000 * 1000,
			InReadable: "9 GB",
		}, {
			Name:       "TB conversion",
			InBytes:    580 * 1000 * 1000 * 1000 * 1000,
			InReadable: "580 TB",
		}, {
			Name:       "0 bytes",
			InBytes:    0,
			InReadable: "0 B",
		},
	}

	for _, tt := range tests {
		szReadable := toReadableSize(tt.InBytes)
		if szReadable != tt.InReadable {
			t.Errorf("toReadable size: expected:%v, got: %v", tt.InReadable, szReadable)
		}
	}
}
