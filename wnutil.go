package gown

import (
	"os"
	"fmt"
	"bytes"
)

type XX struct {
	
}

func do_init() os.Error { // *** Is having datafps and indexfps as globals, the best way to proceed?? ***
	// Find base directory for database, for now, we look only in the env var WNSEARCHDIR
	searchdir := os.Getenv("WNSEARCHDIR")
	var err os.Error

	for i:=1; i<NUMPARTS; i++ {
		datapath := fmt.Sprintf("%s/data.%s", searchdir, partnames[i]) // *** NOT PORTABLE ***
		fmt.Fprintf(os.Stderr, "Opening data file: %s in slot %d\n", datapath, i)
		datafps[i], err = os.Open(datapath) 
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open datafile (%s)\n", datapath)
			return err
		}
	}
	for i:=1; i<NUMPARTS; i++ {
		indexpath := fmt.Sprintf("%s/index.%s", searchdir, partnames[i]) // *** NOT PORTABLE ***
		fmt.Fprintf(os.Stderr, "Opening data file: %s in slot %d\n", indexpath, i)
		indexfps[i], err = os.Open(indexpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WordNet library error: Can't open indexfile (%s)\n", indexpath)
			return err
		}
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

