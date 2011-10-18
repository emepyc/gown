package gown

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"bufio"
)

type indexOffset uint64

type indexMap map[string]indexOffset // TODO: Profile the amount of memory for in-memory indexing

type indexInfo struct {
	lemma []byte,
	pos indexOffset,
	offsets []indexOffset,
	tagsense_cnt int,
}

type Indexer interface {
	Lookup() indexOffset
}

type indexFiles []*os.File
type dataFiles  []*os.File

func do_init() os.Error {
	searchdir := os.Getenv("WNSEARCHDIR") // *** TODO: allow other ways to get the path to these files ***
	                                      // Also check for errors from os.Getenv ??
	wndb, err := loadIndex(searchdir) //
	if err != nil {
		fmt.Fprintf(os.STDERR, "I can't load the indexes files from %s\n", searchdir)
		os.Exit(1) // Return error?
	}
}

// Loads an array of io.Reader's containing handlers for the datafiles
// For now, the indexes are loaded in a map (in memory)
// TODO: Check for WordNet version. It is stated more or less at the beginning of the file
func loadIndex(searchdir []byte) (Indexer, os.Error) {
	fmt.Fprintf(os.Stderr, "Reading Index (in memory)...") // TODO: Profile with timer
	indexMap := make(indexMap, 363000) // Current number of lines in index.*
	for i:=1; i<=NUMPARTS; i++ {
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
				return indexMap, nil
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred while reading line from (%s)\n%s\n", indexpath, line)
			}
			newIndexInfo := parseIndexLine(line)
		}
	}
	fmt.Fprintf(os.Stderr, "Done\n")

	return Indexer(indexMap), nil
}

// Loads an array of io.Reader's containing handlers for the datafiles
func loadData(searchdir []byte) (dataFiles, os.Error) {
	datafps := [NUMPARTS+1]io.Reader
	for i:=1; i<=NUMPARTS; i++ {
		datapath := fmt.Sprintf("%s/data.%s", searchdir, partnames[i]) // TODO: Make this portable
		fmt.Fprintf(os.Stderr, "Opening data file: %s in slot %d\n", datapath, i)
		datafps[i], err := os.Open(datapath) 
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open datafile (%s)\n", datapath)
			return nil,err
		}
	}
	return datafps, nil
}
	/* This file isn't used by the library and doesn't have to
	 be present.  No error is reported if the open fails. */
	sensepath := fmt.Sprintf("%s/index.sense", searchdir) // *** NOT PORTABLE ***
	sensefp, _ = os.Open(sensepath)

	/* If this file isn't present, the runtime code will skip printint out
	 the number of times each sense was tagged. */
	cntlistpath := fmt.Sprintf("%s/cntlist.rev", searchdir) // *** NOT PORTABLE ***
	cntlistfp, _ = os.Open(cntlistpath)

	/* This file doesn't have to be present.  No error is reported if the
	 open fails. */
	keyidpath := fmt.Sprintf("%s/index.key", searchdir)
	keyindexfp, _ = os.Open(keyidpath)

	revkeyindexpath := fmt.Sprintf("%s/index.key.rev", searchdir) // *** NOT PORTABLE ***
	revkeyindexfp, _ = os.Open(revkeyindexpath)

	vsentfilepath := fmt.Sprintf("%s/sents.vrb", searchdir)  // *** NOT PORTABLE ***
	vsentfilefp, err = os.Open(vsentfilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WordNet warning: Can't open verb example sentences file (%s) : %s\n", vsentfilepath, err)
	}

	vidxfilepath := fmt.Sprintf("%s/sentidx.vrb", searchdir) // *** NOT PORTABLE ***
	vidxfilefp, err = os.Open(vidxfilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WordNet warning: Can't open verb example sentences index file (%s): %s", vidxfilepath, err)
	}

	return nil
}

func loadIndex() {
	if verbose {
		fmt.Fprintf(os.Stderr, "Loading Indexes")
	}
	for i:=1; i <= 4; i++ {
				
	}
}

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
	fields := bytes.Fields(l)
	ptr_cnt := fields[PTR_CNT]
	newIndexInfo.lemma = fields[LEMMA]
	newIndexInfo.pos = fields[POS]
	newIndexInfo.tagsense_cnt = fields[TAGSENSE_CNT + ptr_cnt]
	newIndexInfo.offsets = fields[(SYNSET_OFFSET + ptr_cnt) : ]
	return &newIndexInfo
}
