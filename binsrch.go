package gown

import (
	"fmt"
	"io"
	"os"
)

const (
	KEY_LEN = 1024
	LINE_LEN = 1024*25
)

var (
//	line [LINE_LEN]byte
	last_bin_search_offset int64 = 0
)

func read_index(offset int64, fh io.ReaderAt) []byte { // io.ReaderAt?? as a pointer??
	line := make([]byte, LINE_LEN)
	n, err := fh.ReadAt(offset) // We may be reading too much, no?
	if err != nil && err != os.EOF {
		fmt.Fprintf(os.Stderr, "An error occurred while trying to read from file: %s\n", err)
		os.Exit(1)
	}
	return line
}

func bin_search(searchkey []byte, fh io.Reader) []byte {
	var c int
	var top, mid, bot, diff int64
	var length int

	diff = 666
	top = 0
	fileStats, err := fh.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the stats of file: %s\n", err)
		os.Exit(1)
	}
	bot = fileStats.Size
	mid = (bot - top) / 2
	
	line := make([]byte, LINE_LEN)

	for {
		newpos, err := fh.Seek(mid - 1, os.SEEK_SET)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while Seeking: %s\n", err)
			os.Exit(1)
		}
		if newpos != mid - 1 {
			fmt.Fprintf(os.Stderr, "Something happened while Seeking: %d != %d\n", newpos, mid-1)
			os.Exit(1)
		}
		
	}
}
