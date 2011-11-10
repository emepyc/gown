package main

import (
	"fmt"
	"log"
	"./gown"
)

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
		symbol := []byte{'@'}
		hypernyms, err := wn
	}
}

func main() {
//	file := "/Users/pignatelli/Downloads/WordNet-3.0/dict/data.noun"
	file := "/home/mp/Downloads/WordNet-3.0/dict/data.noun"
	fh, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	pos := 8199025  //952963 //  576451 // 37396 //34777
	line := dataLookup(fh, int64(pos))
	fmt.Printf("[LINE] %s\n", line)

	symbol := []byte{';','c'}
	hypernyms, err := GetRelation(line, symbol)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("========\n%s\n==========\n[PTRS]\n", symbol)
	for _, l := range hypernyms {
		fmt.Printf(" [SYMBOL] %s\n", l.symbol)
		fmt.Printf(" [OFFSET] %d\n", l.offset)
		fmt.Printf(" [POS] %c\n", l.pos)
		fmt.Printf(" [SOURCE] %d\n", l.source)
		fmt.Printf(" [TARGET] %d\n", l.target)
	}
	fmt.Printf("==================\n");
//	os.Exit(0)
	data, err := parseDataLine(line)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[SYNSET_OFFSET] %d\n", data.synset_offset)
	fmt.Printf("[LEX_FILENUM] %d\n", data.lex_filenum)
	fmt.Printf("[SS_TYPE] %c\n", data.ss_type)
	fmt.Printf("[W_CNT] %d\n", data.w_cnt)
	fmt.Printf("[LEMMAS]\n")
	for _, l := range data.lemmas {
		fmt.Printf(" [LEMMA] %s %d\n", l.word, l.lex_id)
	}
	fmt.Printf("[P_CNT] %d\n", data.p_cnt)
	fmt.Printf("[PTRS]\n")
	for _, l := range data.ptrs {
		fmt.Printf(" [SYMBOL] %s\n", l.symbol)
		fmt.Printf(" [OFFSET] %d\n", l.offset)
		fmt.Printf(" [POS] %c\n", l.pos)
		fmt.Printf(" [SOURCE] %d\n", l.source)
		fmt.Printf(" [TARGET] %d\n", l.target)
	}
	fmt.Printf("[GLOSS] %s\n", data.gloss)
}
