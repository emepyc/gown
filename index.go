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
	SYMBOL               // 4
	SENSE_CNT            // 5
	TAGSENSE_CNT         // 6
	SYNSET_OFFSET        // 7
)

// Error msgs -- may be refactored in another src file
const (
	_ = iota
	UNKNOWN_WORD
	SYNTACTIC_CATEGORY_TOO_LONG
	LINE_TOO_LONG
)

const (
	NINDEXRECS = 363000 // Current number of lines in index.* (wc -l index.*)
	BUFFSIZE = 3072 // Buffer for reading lines in data
	NUMPARTS = 4 // noun, adj, verb, adverb
	//	NUMFRAMES = 35
)

var partnames []string = []string{"", "noun", "verb", "adj", "adv"}

type indexInfo struct {
	lemma        []byte
	pos          byte // 'n', 'v', 'a', 'r' ... should we have these in a map? probably yes
	p_cnt        int
	ptr_symbols   [][]byte
	tagsense_cnt int
	offsets      []int64
}

// May be refactored in another source file
type ERR_MSG int

// May be refactored in another source file
func (t ERR_MSG) Error() string {
	return errMsg(int(t))
}

type indexMap map[string]indexInfo // TODO: Profile in-memory indexing (better *indexInfo?)
type indexMaps [NUMPARTS+1]indexMap

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
	case UNKNOWN_WORD :
		return "UNKNOWN WORD"
	case SYNTACTIC_CATEGORY_TOO_LONG:
		return "SYNTACTIC CATEGORY TOO LONG IN INDEX FILE (expected values are 'n', 'v', 'a' or 'r'"
	case LINE_TOO_LONG :
		return "LINE TOO LONG"
	default :
		return "UNKNOWN ERROR MSG"
	}
	return "" // will not reach here
}


func (i *indexMaps) Lookup(b []byte, pos int) ([]int64, error) {
	m := i[pos] // returns an indexMap
	lemma, ok := m[string(b)]
	if !ok {
		return nil, ERR_MSG(UNKNOWN_WORD)
	}
	return lemma.offsets, nil
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
	var index indexMaps

	for i := 1; i <= NUMPARTS; i++ {
		indexMap := make(indexMap, NINDEXRECS)
		indexpath := fmt.Sprintf("%s/index.%s", searchdir, partnames[i]) // TODO: Make this portable
		fmt.Fprintf(os.Stderr, "Processing index file %s\n", indexpath)
		indexfh, err := os.Open(indexpath)
		if err != nil {
			return nil, err
		}

		defer func() {
			e := indexfh.Close()
			if e != nil {
				log.Fatal("Problem closing the file")
			}
		}()

		bufindexfh := bufio.NewReader(io.Reader(indexfh))
		nlines := 0
		for {
			nlines++
			fmt.Fprintf(os.Stderr, "\r%d lines                 ", nlines)
			line, isPrefix, err := bufindexfh.ReadLine()
			if isPrefix {
				return nil, ERR_MSG(LINE_TOO_LONG)
			}
			if err == io.EOF {
				index[i] = indexMap
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred while reading line from (%s)\n%s\n", indexpath, line)
			}
			if line[0] == ' ' { // header line
				continue
			}
			newIndexInfo, err := parseIndexLine(line)
			if err != nil {
				return nil, err
			}
			key := string(newIndexInfo.lemma)
			indexMap[key] = *newIndexInfo
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

func parseIndexLine(l []byte) (*indexInfo, error) {
	//	fmt.Fprintf (os.Stderr, "Processing line:\n%s\n", l)
	newIndexInfo := indexInfo{}
	var err error
	fields := bytes.Fields(l)
	newIndexInfo.lemma = fields[LEMMA]
	//	newIndexInfo.pos, err = strconv.Atoui64(string(fields[POS]))
	if len(fields[POS]) > 1 {
		return nil, ERR_MSG(SYNTACTIC_CATEGORY_TOO_LONG)
	}
	newIndexInfo.pos = fields[POS][0]

	ptr_cnt, err := strconv.Atoi(string(fields[PTR_CNT]))
	if err != nil {
		return nil, err
	}

	ptr_symbols := fields[SYMBOL:SYMBOL+ptr_cnt]
	newIndexInfo.ptr_symbols = ptr_symbols

	newIndexInfo.tagsense_cnt, err = strconv.Atoi(string(fields[TAGSENSE_CNT+ptr_cnt]))
	if err != nil {
		return nil, err
	}
	offsets_strs := fields[(SYNSET_OFFSET + ptr_cnt - 1):]
	offsets := make([]int64, len(offsets_strs))
	for i, offset := range offsets_strs {
		offsets[i], err = strconv.Atoi64(string(offset))
		if err != nil {
			return nil, err
		}
	}
	newIndexInfo.offsets = offsets
	return &newIndexInfo, nil
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
