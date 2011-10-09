package gown

import (
	"fmt"
	"os"
)

const (
	DONT_KNOW  =  iota
	DIRECT_ANT           // direct antonyms (cluster head)
	INDIRECT_ANT         // indirect antonyms (similar)
	PERTAINYM            // no antonyms or similars (pertainyms)
)

const (
	ALLWORDS = 0       // print all words
	SKIP_ANTS = 0      // skip printing antonyms in printsynset()
	PRINT_ANTS = 1     // print antonyms in printsynset()
	SKIP_MARKER = 0    // skip printing adjective marker
	PRINT_MARKER = 1   // print adjective marker
)

const (
	TRACEP = iota + 1  // traceptrs
	TRACEC             // tracecoords()
	TRACEI             // traceinherit()
)

const (
	DEFOFF = iota
	DEFON
)

const TMPBUFSIZE = 1024*10

var (
	prflag, sense, prlexid int
	overflag int = 0              // set when output buffer overflows *** PROBABLY NOT NEEDED ***
	searchbuffer []byte
	lastholomero int              // keep track of last holo/meronym printed
        tmpbuf []byte                 // general purpose printing buffer *** PROBABLY NOT NEEDED OR bufio ***
	wdbuf []byte                  // general purpose word buffer *** PROBABLY NOT NEEDED OR bufio ***
	msgbuf []byte                 // buffer for constructing error messages *** NOT NEEDED OR bufio ? ***
)

//type IndexDB struct {
	// ???
//}

type Index struct {
	idxoffset int64 // byte offset of entry in index file
	wd []byte // word string
	pos []byte // part of speech
	sense_cnt int // sense (collins) count
	off_cnt int // number of offsets
	tagget_cnt int // number of senses that are tagged
	offset uint64 // offsets of synsets containing word
	ptruse_cnt int // number of pointers used
	ptruse int // pointers used. ***IN wn.h THIS IS A POINTER TO A INT***
}

type Synset struct {
	hereiam int64 // current file position
	ssypte int // type of ADJ synset
	fnum int // file number that synset comes from
	pos []byte // part of speech
	wcount int // number of words in synset
	words [][]byte // words in synset
	lexid int //unique id in lexicographer file *** In wn.h this is a pointer to a int ***
	wnsns int // sense number in wordnet *** In wn.h this is a pointer to a int ***
	whichword int // which word in synsest we're looking for
	ptrcount int // number of pointers
	ptrtyp int // pointer types *** In wn.h this is a pointer to a int ***
	ptroff int64 // pointer offsets *** In wn.h this is a pointer to a long ***
	ppos int // pointer part of speech *** In wn.h this is a pointer to a int ***
	pto int // pointer 'to' fields *** In wn.h this is a pointer to a int ***
	pfrm int // pointer 'from' fields *** In wn.h this is a pointer to a int ***
	fcount int // number of verb frames
	frmid int // frame numbers *** In wn.h this is a pointer to a int ***
	frmto int // frame 'to' fields *** In wn.h this is a pointer to a int ***
	defn []byte // synset gloss (definition)
	key uint32 // unique synset key

/* these fields are used if a data structure is returned
 instead of a text buffer */

	nextss *Synset // ptr to next synset containing searchword
	nextform *Synset // ptr to list of synsets for alternate spelling of wordform

	searchtype int // type of search performed
	ptrlist []*Synset // ptr to synset list result of search
	headword []byte // if pos is "s", this is cluster head word
	headsense int16 // sense number of headword
}

type SnsIndex struct {
	sensekey []byte // sense key
	word []byte // word string
	loc int64 // synset offset
	wnsense int // WordNet sense number
	tag_cnt int // number of semantic tags to sense
	nextsi *SnsIndex // ptr to next sense index entry
}

type SearchResults struct {
	SenseCount [MAX_FORMS]byte // number of senses word form has
	OutSenseCount [MAX_FORMS]byte // number of senses printed for word form
	numforms int // number of word forms searchword has
	printcnt int // number of senses printed by search
	searchbuf []byte // buffer containing formatted results *** change to bufio.New ? ***
	searchds *Synset // data structure containing search results
}

// Find word in index file and return parsed entry in data structure
// Input word must be exact match of string in database
func (i Index) index_lookup (word []byte, dbase int) *Index {
	idx := &Index{}
	if dbase > len(i.dbs)-1 {
		return nil
	}
	if i.dbs[i] == nil {
		fmt.Fprintf(os.Stderr, "WordNet library error: %s indexfile not open\n");
		return nil
	}

	
}
