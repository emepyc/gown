package gown

import (
	"io"
)

const (
	wnrelease = "3.0"
)

// buffer sizes *** should these constants be in another source file? ***
const (
	SEARCHBUF = (200*1024) // *** long in wn.h ***
	LINEBUF = (15*1024) /* 15K buffer to read index & data files */
	SMLINEBUF = (3*1024) /* small buffer for output lines */
	WORDBUF = 256 /* buffer for one word or collocation */
)

const (
	ALLSENSES = 0 // pass to findtheinfo() if want all senses
	MAXID = 15 // maximum id number in lexicographer file
	MAXDEPTH = 20 // maximum tree depth - used to find cycles
	MAXSENSE = 75 // maximum number of senses in database
	MAX_FORMS = 5 // max # of different 'forms' word can have
	MAXFNUM = 44 // maximum number of lexicographer files
)

/* Pointer type and search type counts */
/* Pointers */
const (
	_ = iota
	ANTPTR /* ! */
	HYPERPTR /* @ */
	HYPOPTR /* ~ */
	ENTAILPTR /* * */
	SIMPTR /* & */

	ISMEMBERPTR /* #m */
	ISSTUFFPTR /* #s */
	ISPARTPTR /* #p */
	HASMEMBERPTR /* %m */
	HASSTUFFPTR /* %s */
	HASPARTPTR /* %p */

	MERONYM /* % (not valid in lexicographer file) */
	HOLONYM /* # (not valid in lexicographer file) */
	CAUSETO /* > */
	PPLPTR /* < */
	SEEALSOPTR /* ^ */
	PERTPTR /* \ */
	ATRIBUTE /* = */
	VERBGROUP /* $ */
	DERIVATION /* + */
	CLASSIFICATION /* ; */
	CLASS /* - */
)

/* Misc searches */
const (
	LASTTYPE = CLASS
	SYNX = LASTTYPE + 1
	FREQ = LASTTYPE + 2
	FRAMES = LASTTYPE + 3
	COORDS = LASTTYPE + 4
	RELATIVES = LASTTYPE + 5
	HMERONYM = LASTTYPE + 6
	HHOLONYM = LASTTYPE + 7
	WNGREP = LASTTYPE + 8
	OVERVIEW = LASTTYPE + 9

	MAXSEARCH = OVERVIEW

	CLASSIF_START = MAXSEARCH + 1

	CLASSIF_CATEGORY = CLASSIF_START // ;c
	CLASSIF_USAGE = CLASSIF_START + 1 // ;u
	CLASSIF_REGIONAL = CLASSIF_START + 2 // ;r

	CLASSIF_END = CLASSIF_REGIONAL

	CLASS_START = CLASSIF_END + 1

	CLASS_CATEGORY = CLASS_START // -c
	CLASS_USAGE = CLASS_START + 1 // -u
	CLASS_REGIONAL = CLASS_START + 2 // -r

	CLASS_END = CLASS_REGIONAL

	INSTANCE = CLASS_END + 1 // @i
	INSTANCES = CLASS_END + 2 // ~i

	MAXPTR = INSTANCES
)

const (
	NUMPARTS = 4
	NUMFRAMES = 35
)

const (
	_ = iota
	NOUN
	VERB
	ADJ
	ADV
	SATELLITE
	ADJSAT = SATELLITE
)

const (
	ALL_POS = iota
	PADJ // (p)
	NPADJ // (a)
	IPADJ // (ip)
)

var (
//	datafps [NUMPARTS + 1]io.Reader             // *** with or without pointer?? ***
//	indexfps [NUMPARTS + 1]io.Reader            // *** with or without pointer?? ***
	sensefp, cntlistfp, keyindexfp, revkeyindexfp, vidxfilefp, vsentfilefp io.Reader
)


var lexfiles []string = []string{
	"adj.all",		/* 0 */
	"adj.pert",		/* 1 */
	"adv.all",		/* 2 */
	"noun.Tops",		/* 3 */
	"noun.act",		/* 4 */
	"noun.animal",		/* 5 */
	"noun.artifact",	/* 6 */
	"noun.attribute",	/* 7 */
	"noun.body",		/* 8 */
	"noun.cognition",	/* 9 */
	"noun.communication",	/* 10 */
	"noun.event",		/* 11 */
	"noun.feeling",		/* 12 */
	"noun.food",		/* 13 */
	"noun.group",		/* 14 */
	"noun.location",	/* 15 */
	"noun.motive",		/* 16 */
	"noun.object",		/* 17 */
	"noun.person",		/* 18 */
	"noun.phenomenon",	/* 19 */
	"noun.plant",		/* 20 */
	"noun.possession",	/* 21 */
	"noun.process",		/* 22 */
	"noun.quantity",	/* 23 */
	"noun.relation",	/* 24 */
	"noun.shape",		/* 25 */
	"noun.state",		/* 26 */
	"noun.substance",	/* 27 */
	"noun.time",		/* 28 */
	"verb.body",		/* 29 */
	"verb.change",		/* 30 */
	"verb.cognition",	/* 31 */
	"verb.communication",	/* 32 */
	"verb.competition",	/* 33 */
	"verb.consumption",	/* 34 */
	"verb.contact",		/* 35 */
	"verb.creation",	/* 36 */
	"verb.emotion",		/* 37 */
	"verb.motion",		/* 38 */
	"verb.perception",	/* 39 */
	"verb.possession",	/* 40 */
	"verb.social",		/* 41 */
	"verb.stative",		/* 42 */
	"verb.weather",		/* 43 */
	"adj.ppl",		/* 44 */

}

var ptrtyp []string = []string{
	"",
	"!",
	"@",			/* 2 HYPERPTR */
	"~",			/* 3 HYPOPTR */
	"*",			/* 4 ENTAILPTR */
	"&",			/* 5 SIMPTR */
	"#m",			/* 6 ISMEMBERPTR */
	"#s",			/* 7 ISSTUFFPTR */
	"#p",			/* 8 ISPARTPTR */
	"%m",			/* 9 HASMEMBERPTR */
	"%s",			/* 10 HASSTUFFPTR */
	"%p",			/* 11 HASPARTPTR */
	"%",			/* 12 MERONYM */
	"#",			/* 13 HOLONYM */
	">",			/* 14 CAUSETO */
	"<",			/* 15 PPLPTR */
	"^",			/* 16 SEEALSO */
	"\\",			/* 17 PERTPTR */
	"=",			/* 18 ATTRIBUTE */
	"$",			/* 19 VERBGROUP */
	"+",		        /* 20 NOMINALIZATIONS */
	";",			/* 21 CLASSIFICATION */
	"-",			/* 22 CLASS */
	/* additional searches, but not pointers.  */
	"",			/* SYNS */
	"",			/* FREQ */
	"+",			/* FRAMES */
	"",			/* COORDS */
	"",			/* RELATIVES */
	"",			/* HMERONYM */
	"",			/* HHOLONYM */
	"",			/* WNGREP */
	"",			/* OVERVIEW */
	";c",			/* CLASSIF_CATEGORY */
	";u",			/* CLASSIF_USAGE */
	";r",			/* CLASSIF_REGIONAL */
	"-c",			/* CLASS_CATEGORY */
	"-u",			/* CLASS_USAGE */
	"-r",			/* CLASS_REGIONAL */
	"@i",			/* INSTANCE */
	"~i",			/* INSTANCES */
}

var partnames []string = []string{"", "noun", "verb", "adj", "adv"}
var partchars string = " nvara"  // Add char for satellites to end
var adjclass []string = []string{"", "(p)", "(a)", "(ip)"}

// Text of verb sentence frames

var frametext []string = []string{
   "",
    "Something ----s",
    "Somebody ----s",
    "It is ----ing",
    "Something is ----ing PP",
    "Something ----s something Adjective/Noun",
    "Something ----s Adjective/Noun",
    "Somebody ----s Adjective",
    "Somebody ----s something",
    "Somebody ----s somebody",
    "Something ----s somebody",
    "Something ----s something",
    "Something ----s to somebody",
    "Somebody ----s on something",
    "Somebody ----s somebody something",
    "Somebody ----s something to somebody",
    "Somebody ----s something from somebody",
    "Somebody ----s somebody with something",
    "Somebody ----s somebody of something",
    "Somebody ----s something on somebody",
    "Somebody ----s somebody PP",
    "Somebody ----s something PP",
    "Somebody ----s PP",
    "Somebody's (body part) ----s",
    "Somebody ----s somebody to INFINITIVE",
    "Somebody ----s somebody INFINITIVE",
    "Somebody ----s that CLAUSE",
    "Somebody ----s to somebody",
    "Somebody ----s to INFINITIVE",
    "Somebody ----s whether INFINITIVE",
    "Somebody ----s somebody into V-ing something",
    "Somebody ----s something with something",
    "Somebody ----s INFINITIVE",
    "Somebody ----s VERB-ing",
    "It ----s that CLAUSE",
    "Something ----s INFINITIVE",
    "",
}