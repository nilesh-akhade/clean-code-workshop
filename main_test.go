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
				Total:     5,
				Duplicate: 2,
				DupSize:   0,
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
