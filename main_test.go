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
			Err       error
		}
	}{
		{
			Name:      "No duplicates",
			Directory: "./testdata/testno-dupes",
			Want: struct {
				Total     int
				Duplicate int
				DupSize   int64
				Err       error
			}{
				Total:     2,
				Duplicate: 0,
				DupSize:   0,
				Err:       nil,
			},
		},
		{
			Name:      "4 duplicates",
			Directory: "./testdata/some-duplicates",
			Want: struct {
				Total     int
				Duplicate int
				DupSize   int64
				Err       error
			}{
				Total:     4,
				Duplicate: 1,
				DupSize:   11,
				Err:       nil,
			},
		},
		{
			Name:      "Dir not found",
			Directory: "./testdata/notfound",
			Want: struct {
				Total     int
				Duplicate int
				DupSize   int64
				Err       error
			}{
				Total:     0,
				Duplicate: 0,
				DupSize:   0,
				Err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var dupeSize int64
			dupeDetails := &DuplicateDetails{
				Hashes:     map[string]string{},
				Duplicates: map[string]string{},
				DupeSize:   &dupeSize,
			}
			entries, _ := ioutil.ReadDir(tt.Directory)
			// if err != nil {
			// 	t.Error(err)
			// }
			err := traverseDir(dupeDetails, entries, tt.Directory)
			if err != tt.Want.Err {
				t.Errorf("error: expected:%v, got: %v", tt.Want.Err, err)
			}
			if dupeSize != tt.Want.DupSize {
				t.Errorf("size: expected:%v, got: %v", tt.Want.DupSize, dupeSize)
			}
			if len(dupeDetails.Hashes) != tt.Want.Total {
				t.Errorf("total: expected:%v, got: %v", tt.Want.Total, len(dupeDetails.Hashes))
			}
			if len(dupeDetails.Duplicates) != tt.Want.Duplicate {
				t.Errorf("duplicates: expected:%v, got: %v", tt.Want.Duplicate, len(dupeDetails.Duplicates))
			}
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
