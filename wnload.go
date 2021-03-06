package gown

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"bufio"
	"strconv"
)

const NINDEXRECS = 363000 // Current number of lines in index.* (wc -l index.*)
const BUFFSIZE = 3072 // for reading lines in data

type indexMap map[string][]int64 // TODO: Profile in-memory indexing
type tagSenseMap map[string]int
type InMemIndex struct {
	indexMaps [NUMPARTS+1]indexMap
	tagSenseMaps [NUMPARTS+1]tagSenseMap
}

// I think this is not used
type indexInfo struct {
	lemma []byte
	pos byte // 'n', 'v', 'a', 'r' ... should we have these in a map? probably yes
	offsets []int64
	tagsense_cnt int
}

type Indexer interface {
	Lookup([]byte, int) ([]int64, os.Error)
}

type indexFiles []io.Reader
type dataFiles  []io.Reader

type WordNetDb struct {
	Index Indexer
	Data dataFiles
}

func (m indexMap) Lookup(word []byte, db int) ([]int64, os.Error) {
	offsets, ok := m[string(word)]
	if !ok {
		newErr := os.NewError("Word <" + word + "> not present in db")
		return nil, newErr
	}
	return offsets, nil
}

func (d Data) dataLookup(pos int, offset int64) []byte { // ret os.Error
	fh := d[pos]
	os.File(fh).Seek(offset, os.SEEK_SET)
	
	buffer := make([]byte, BUFFSIZE) // initial size of the buffer is 3kb
	line := make([]byte, 0, BUFFSIZE)
	prevLen := 0
	for {
		prevLen = len(line)
		n, err := os.File(fh).Read(buffer) // we read the next 3kb (or less)
		if err != nil && err != os.EOF {
			fmt.Fprintf(os.Stderr, "Error while reading file: %s\n", err)
			os.Exit(1) // return err?
		}
		line = append(line, buffer)
		until := IndexByte(buffer, '\n')
		if until > 0 { // We have a full line
			return line[:prevLen+until]
		}
		if err == os.EOF || n < BUFFSIZE {
			return line
		}
	}
	return nil // We cant' be here
}

func New() *WordNetDb {
	searchdir := os.Getenv("WNSEARCHDIR")
	var err os.Error
	wndb := WordNetDb{}

	wndb.Index, err = loadIndex(searchdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I can't load the indexes files from %s\n", searchdir)
		os.Exit(1) // Return error? Without indexes we can't do anything
	}

	wndb.Data, err = dataFh(searchdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I can't open the data files from %s\n", searchdir)
		os.Exit(1)
	}
	return &wndb
}

// Loads an Indexer
// For now, the indexes are loaded in a map (in memory)
// TODO: Check for WordNet version. It is stated more or less at the beginning of the file
func loadIndex(searchdir string) (Indexer, os.Error) {
	fmt.Fprintf(os.Stderr, "Reading Index (in memory)...\n") // TODO: Time this
	index := InMemIndex{}

	for i:=1; i<=NUMPARTS; i++ {
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
			if e!=nil {
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
			if err == os.EOF {
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

	return index, nil
}

// Loads an array of io.Reader's containing handlers for the datafiles
func dataFh(searchdir string) (dataFiles, os.Error) {
	var err os.Error
	datafps := make([]io.Reader, NUMPARTS+1)
	for i:=1; i<=NUMPARTS; i++ {
		datapath := fmt.Sprintf("%s/data.%s", searchdir, partnames[i]) // TODO: Make this portable
		fmt.Fprintf(os.Stderr, "Opening data file: %s in slot %d\n", datapath, i)
		datafps[i], err = os.Open(datapath) 
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open datafile (%s)\n", datapath)
			return nil,err
		}
	}
	return datafps, nil
}
// 	/* This file isn't used by the library and doesn't have to
// 	 be present.  No error is reported if the open fails. */
// 	sensepath := fmt.Sprintf("%s/index.sense", searchdir) // *** NOT PORTABLE ***
// 	sensefp, _ = os.Open(sensepath)

// 	/* If this file isn't present, the runtime code will skip printint out
// 	 the number of times each sense was tagged. */
// 	cntlistpath := fmt.Sprintf("%s/cntlist.rev", searchdir) // *** NOT PORTABLE ***
// 	cntlistfp, _ = os.Open(cntlistpath)

// 	/* This file doesn't have to be present.  No error is reported if the
// 	 open fails. */
// 	keyidpath := fmt.Sprintf("%s/index.key", searchdir)
// 	keyindexfp, _ = os.Open(keyidpath)

// 	revkeyindexpath := fmt.Sprintf("%s/index.key.rev", searchdir) // *** NOT PORTABLE ***
// 	revkeyindexfp, _ = os.Open(revkeyindexpath)

// 	vsentfilepath := fmt.Sprintf("%s/sents.vrb", searchdir)  // *** NOT PORTABLE ***
// 	vsentfilefp, err = os.Open(vsentfilepath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "WordNet warning: Can't open verb example sentences file (%s) : %s\n", vsentfilepath, err)
// 	}

// 	vidxfilepath := fmt.Sprintf("%s/sentidx.vrb", searchdir) // *** NOT PORTABLE ***
// 	vidxfilefp, err = os.Open(vidxfilepath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "WordNet warning: Can't open verb example sentences index file (%s): %s", vidxfilepath, err)
// 	}

// 	return nil
// }

/* Count the number of underscore or space separated words in a string. */
func cntwords(s []byte, separator byte) (wdcnt int) {
	wdcnt = 0
	for i:=0; i<len(s); i++ {
		if s[i] == separator || s[i] == ' ' || s[i] == '_' {
			wdcnt++
			for ; i<len(s); i++ {
				if s[i] != separator && s[i] != ' ' && s[i] != '_' {
					break
				}
			}
		} 
	}
	wdcnt++
	return
}

/* Replace all occurences of 'from' with 'to' in 'str' */
func strsubst(src []byte, from, to byte) []byte {
	dest := make([]byte, len(src))
	for i, ch := range(src) {
		if ch == from {
			dest[i] = to
		} else {
			dest[i] = ch
		}
	}
	return dest
}

func parseIndexLine(l []byte) *indexInfo {
//	fmt.Fprintf (os.Stderr, "Processing line:\n%s\n", l)
	newIndexInfo := indexInfo{}
	var err os.Error
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

newIndexInfo.tagsense_cnt, err = strconv.Atoi(string(fields[TAGSENSE_CNT + ptr_cnt]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "I had a problem trying to convert %s to int32\n", fields[TAGSENSE_CNT + ptr_cnt])
		os.Exit(1)
	}
	offsets_strs := fields[(SYNSET_OFFSET + ptr_cnt) : ]
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
