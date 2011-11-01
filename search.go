package gown

import (
	"fmt"
	"os"
	"strings"
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

type Query struct {
	option []byte, // user's search request
	search int, // search to pass findtheinfo()
	pos int, // part-of-speech to pass findtheinfo()
	sense int, // Perform search on sense number # only
	rel []byte // Relation
	// helpmsgidx int, // index into help message table
	// label []char, // text for search header message
}

var relNameSyn = map[string][]string {
	"ants" : []string{"!"},
	"hype" : []string{"@"},
	"inst" : []string{"@i"},
	"hypes" : []string{"@", "@i"},
}

// type Index struct {
// 	idxoffset int64 // byte offset of entry in index file
// 	wd []byte // word string
// 	pos []byte // part of speech
// 	sense_cnt int // sense (collins) count
// 	off_cnt int // number of offsets
// 	tagget_cnt int // number of senses that are tagged
// 	offset uint64 // offsets of synsets containing word
// 	ptruse_cnt int // number of pointers used
// 	ptruse int // pointers used. ***IN wn.h THIS IS A POINTER TO A INT***
// }

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
// Input word must be exact match of string in database (case insensitive)
func (wndb *WordNetDb) indexOffsetLookup (query Query) []int64 {
	lword := strToLower(query.option)
	offsets := wndb.Index.Lookup(lword, query.pos)

	if query.sense != 0 {
		offsets = offsets[:query.sense]
	}

	return offsets
}


func (wndb *WordNetDb) QuerySense (query *Query) [][]byte {
	fmt.Fprintf(os.Stderr, "(QuerySense) STRING=%s#%d#%d#%d\n", query.option, query.search, query.pos, query.sense) // Only for dubugging

	lword := strToLower(query.option)

	if query.sense != 0 {
		if query.rel == nil {
			fmt.Fprintf(os.Stderr, "(QuerySense) Relation required\n")
			os.Exit(1) // return err?
		}
		if rel, ok := relSymName[string(query.rel)]; ok {
			query.rel = rel
		} else if _, ok := relNameSym[string(query.rel)]; ok || string(query.rel) == "glos" || string(query.rel) == "syns" {
		} else {
			fmt.Fprintf(os.Stderr, "(QuerySense) Bad relation: %s\n", query.rel)
			os.Exit(1) // return err?
		} 

		offsets = wndb.indexOffsetLookup(lword, query.pos, query.sense)
		offset = offset[0]
		dataLine := dataLookup(query.pos, offset)
		
	}

// 	if query.pos == ALL_POS {
// 		outsenses := make([][]byte, 4) // 4?
// 		for pos:=1; pos <= NUMPARTS; pos++ {
// 			sense := wndb.Lookup(query)
// 			outsenses[pos-1] = sense // pos-1??
// 		}
// 		return outsenses
// 	}
// 	outsenses := make([][]byte,1)
// [...]
	
}


/* Convert to lowercase and remove trailing adjective marker if found */
func strToLower(str []byte) []byte {
	ret := make ([]byte, len(str))
	for i, ch := range str {
		if ch == '(' {
			return bytes.ToLower(ret)
		}
		ret[i] = ch
	}
	return bytes.ToLower(ret)
}

