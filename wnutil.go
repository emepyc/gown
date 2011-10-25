package gown

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"bufio"
)

const NINDEXRECS = 363000 // Current number of lines in index.* (wc -l index.*)

type indexMap map[string]uint64 // TODO: Profile in-memory indexing
type tagSenseMap map[string]int
type InMemIndex struct {
	indexMaps [4]indexMap
	tagSenseMaps [4]tagSenseMap
}

type indexInfo struct {
	lemma []byte
	pos uint64
	offsets []uint64
	tagsense_cnt int
}

type Indexer interface {
//	Lookup() indexOffset
}

type indexFiles []*os.File
type dataFiles  []*os.File

type WordNetDb struct {
	Index Indexer
	Data dataFiles
}

func New() *WordNetDb {
	searchdir := os.Getenv("WNSEARCHDIR")

	wndb, err := loadIndex(searchdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I can't load the indexes files from %s\n", searchdir)
		os.Exit(1) // Return error? Without indexes we can't do anything
	}

	
}

// Loads an Indexer
// For now, the indexes are loaded in a map (in memory)
// TODO: Check for WordNet version. It is stated more or less at the beginning of the file
func loadIndex(searchdir string) (Indexer, os.Error) {
	fmt.Fprintf(os.Stderr, "Reading Index (in memory)...") // TODO: Time this
	index := InMemIndex{}
	for i:=1; i<=NUMPARTS; i++ {
		indexMap := make(indexMap, NINDEXRECS)
		tagSenseMap := make(tagSenseMap, NINDEXRECS)
		indexpath := fmt.Sprintf("%s/index.%s", searchdir, partnames[i]) // TODO: Make this portable
		indexfh, err := os.Open(indexpath)
		bufindexfh := bufio.NewReader(io.Reader(indexfh))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open indexfile (%s)\n", indexpath)
			return nil, err
		}
		for {
			line, isPrefix, err := bufindexfh.ReadLine()
			if isPrefix {
				fmt.Fprintf(os.Stderr, "Line too long reading file (%s)\n", indexpath)
				os.Exit(1) // TODO: Exit?
			}
			if err == os.EOF {
				index.indexMaps[i] = indexMap
				index.tagSenseMaps[i] = tagSenseMap
				continue
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred while reading line from (%s)\n%s\n", indexpath, line)
			}
			newIndexInfo := parseIndexLine(line)
			key := string(newIndexInfo.lemma)
			indexMap[key] = newIndexInfo.offsets
			tagSenseMap[key] = newIndexInfo.tagsense_cnt
		}
	}
	fmt.Fprintf(os.Stderr, "Done\n")

	return Indexer(indexMap), nil
}

// Loads an array of io.Reader's containing handlers for the datafiles
func dataFh(searchdir []byte) (dataFiles, os.Error) {
	var err = os.Error
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

/* Convert to lowercase and remove trailing adjective marker if found */
func strtolower(str []byte) []byte {
	ret := make ([]byte, len(str))
	for i, ch := range str {
		if ch == '(' {
			return bytes.ToLower(ret)
		}
		ret[i] = ch
	}
	return bytes.ToLower(ret)
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
	newIndexInfo := indexInfo{}
	var error os.Error
	fields := bytes.Fields(l)
	ptr_cnt := fields[PTR_CNT]
	newIndexInfo.lemma = fields[LEMMA]
	newIndexInfo.pos, err = strconv.Atoi64(string(fields[POS]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "I had a problem trying to convert %s to uint64\n", fields[POS])
		os.Exit(1) // log.Fatal?
	}
	newIndexInfo.tagsense_cnt, err = strconv.Atoi32(string(fields[TAGSENSE_CNT + ptr_cnt]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "I had a problem trying to convert %s to int32\n", fields[TAGSENSE_CNT + ptr_cnt])
		os.Exit(1)
	}
	offsets_strs := fields[(SYNSET_OFFSET + ptr_cnt) : ]
	offsets := make([]uint64, len(offsets_str))
	for i, offset := range offsets_strs {
		offsets[i], err = strconv.Atoi64(offset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I had a problem trying to convert the offset %s to int63\n", offset)
			os.Exit(1) // log.Fatal?
		}
	}
	newIndexInfo.offsets = offsets
	return &newIndexInfo
}
