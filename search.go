package gown

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
	overflag int = 0              // set when output buffer overflows - PROBABLY NOT NEEDED
	searchbuffer []char
	lastholomero int              // keep track of last holo/meronym printed
        tmpbuf []char                 // general purpose printing buffer - PROBABLY NOT NEEDED
	wdbuf []char                  // general purpose word buffer - PROBABLY NOT NEEDED
	msgbuf []char                 // buffer for constructing error messages - NOT NEEDED?
)

// Find word in index file and return parsed entry in data structure
// Input word must be exact match of string in database

func index_lookup (word []char, dbase int) {
	
}
