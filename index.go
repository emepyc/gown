package main

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"bufio"
	"strconv"
	"log"
)

const (
	LEMMA         = iota // 0
	POS                  // 1
	SYNSET_CNT           // 2
	PTR_CNT              // 3
	SENSE_CNT            // 4
	TAGSENSE_CNT         // 5
	SYNSET_OFFSET        // 6
)

// Error msgs
const (
	_ = iota
	UNKNOWN_WORD
)

const NINDEXRECS = 363000 // Current number of lines in index.* (wc -l index.*)
const BUFFSIZE = 3072     // for reading lines in data
const (
	NUMPARTS = 4
	//	NUMFRAMES = 35
)

var partnames []string = []string{"", "noun", "verb", "adj", "adv"}

type indexInfo struct {
	lemma        []byte
	pos          byte // 'n', 'v', 'a', 'r' ... should we have these in a map? probably yes
	offsets      []int64
	tagsense_cnt int
}

type ERR_MSG int

func (t ERR_MSG) Error() string {
	return errMsg(int(t))
}

type indexMap map[string][]int64 // TODO: Profile in-memory indexing
type tagSenseMap map[string]int
type InMemIndex struct {
	indexMaps    [NUMPARTS + 1]indexMap
	tagSenseMaps [NUMPARTS + 1]tagSenseMap
}

type Indexer interface {
	Lookup([]byte, int) ([]int64, error)
}

type indexFiles []io.Reader
type dataFiles []io.Reader

type WordNetDb struct {
	Index Indexer
	Data  dataFiles
}

func errMsg(n int) string {
	switch n {
	case 1 :
		return "UNKNOWN WORD"
	default :
		return "UNKNOWN ERROR MSG"
	}
	return "" // will not reach here
}

func (i *InMemIndex) Lookup(b []byte, pos int) ([]int64, error) {
	m := i.indexMaps[pos]
	offsets, ok := m[string(b)]
	if !ok {
		return nil, ERR_MSG(UNKNOWN_WORD)
	}
	return offsets, nil
}

func New() (*WordNetDb, error) {
	searchdir := os.Getenv("WNSEARCHDIR")
	var err error
	wndb := WordNetDb{}

	wndb.Index, err = loadIndex(searchdir)
	if err != nil {
		return nil, err
	}

	wndb.Data, err = dataFh(searchdir)
	if err != nil {
		return nil, err
	}
	return &wndb, nil
}

// Loads an Indexer
// For now, the indexes are loaded in a map (in memory)
// TODO: Check for WordNet version. It is stated more or less at the beginning of the file
func loadIndex(searchdir string) (Indexer, error) {
	fmt.Fprintf(os.Stderr, "Reading Index (in memory)...\n") // TODO: Time this
	index := InMemIndex{}

	for i := 1; i <= NUMPARTS; i++ {
		indexMap := make(indexMap, NINDEXRECS)
		tagSenseMap := make(tagSenseMap, NINDEXRECS)
		indexpath := fmt.Sprintf("%s/index.%s", searchdir, partnames[i]) // TODO: Make this portable
		fmt.Fprintf(os.Stderr, "Processing index file %s\n", indexpath)
		indexfh, err := os.Open(indexpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open indexfile (%s)\n", indexpath)
			return nil, err
		}

		defer func() {
			e := indexfh.Close()
			if e != nil {
				fmt.Fprintf(os.Stderr, "Problem closing the file %s: %s\n", indexpath, e)
			}
		}()

		bufindexfh := bufio.NewReader(io.Reader(indexfh))
		nlines := 0
		for {
			nlines++
			fmt.Fprintf(os.Stderr, "\r%d lines                 ", nlines)
			line, isPrefix, err := bufindexfh.ReadLine()
			if isPrefix {
				fmt.Fprintf(os.Stderr, "Line too long reading file (%s)\n", indexpath)
				os.Exit(1) // TODO: Exit?
			}
			if err == io.EOF {
				index.indexMaps[i] = indexMap
				index.tagSenseMaps[i] = tagSenseMap
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred while reading line from (%s)\n%s\n", indexpath, line)
			}
			if line[0] == ' ' { // header line
				continue
			}
			newIndexInfo := parseIndexLine(line)
			key := string(newIndexInfo.lemma)
			indexMap[key] = newIndexInfo.offsets
			tagSenseMap[key] = newIndexInfo.tagsense_cnt
		}
	}
	fmt.Fprintf(os.Stderr, "Done\n")

	return &index, nil
}

// Loads an array of io.Reader's containing handlers for the datafiles
func dataFh(searchdir string) (dataFiles, error) {
	var err error
	datafps := make([]io.Reader, NUMPARTS+1)
	for i := 1; i <= NUMPARTS; i++ {
		datapath := fmt.Sprintf("%s/data.%s", searchdir, partnames[i]) // TODO: Make this portable
		fmt.Fprintf(os.Stderr, "Opening data file: %s in slot %d\n", datapath, i)
		datafps[i], err = os.Open(datapath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open datafile (%s)\n", datapath)
			return nil, err
		}
	}
	return datafps, nil
}

func parseIndexLine(l []byte) *indexInfo {
	//	fmt.Fprintf (os.Stderr, "Processing line:\n%s\n", l)
	newIndexInfo := indexInfo{}
	var err error
	fields := bytes.Fields(l)
	ptr_cnt, err := strconv.Atoi(string(fields[PTR_CNT]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "I had a problem trying to convert '%s' to int\n", fields[PTR_CNT])
		os.Exit(1)
	}
	newIndexInfo.lemma = fields[LEMMA]
	//	newIndexInfo.pos, err = strconv.Atoui64(string(fields[POS]))
	if len(fields[POS]) > 1 {
		fmt.Fprintf(os.Stderr, "POS has to be 1 letter code ('n', 'v', 'a' or 'r') and I have %s\n", fields[POS])
		os.Exit(1)
	}
	newIndexInfo.pos = fields[POS][0]

	newIndexInfo.tagsense_cnt, err = strconv.Atoi(string(fields[TAGSENSE_CNT+ptr_cnt]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "I had a problem trying to convert %s to int32\n", fields[TAGSENSE_CNT+ptr_cnt])
		os.Exit(1)
	}
	offsets_strs := fields[(SYNSET_OFFSET + ptr_cnt):]
	offsets := make([]int64, len(offsets_strs))
	for i, offset := range offsets_strs {
		offsets[i], err = strconv.Atoi64(string(offset))
		if err != nil {
			fmt.Fprintf(os.Stderr, "I had a problem trying to convert the offset %s to int63\n", offset)
			os.Exit(1) // log.Fatal?
		}
	}
	newIndexInfo.offsets = offsets
	return &newIndexInfo
}

func main() {
	gownDb, err := New()
	if err != nil {
		log.Fatal(err)
	}
	offsets, err := gownDb.Index.Lookup([]byte("action"), 1)
	if err != nil {
		log.Fatal(err)
	}
	for _, offset := range offsets {
		fmt.Printf("[OFFSET] %d\n", offset)
	}
}
