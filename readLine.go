package main

import (
	"errors"
	"fmt"
	"bytes"
	"io"
	"os"
	"log"
	"strconv"
	"encoding/hex" // only for the error types
)

type Data os.File

const BUFFSIZE = 10

type lemma struct {
	word   []byte
	lex_id int // 1-digit hexadecimal integer
}

type synsetPtr struct {
	symbol    []byte
	offset    int64
	pos       byte
	source    int
	target    int
}

type dataData struct {
	synset_offset int64 // Current byte offset in the file represented as an 8-digit dec integer
	lex_filenum   int   // 2-digit integer
	ss_type       byte  // n => NOUN, v => VERB, a => ADJECTIVE, s => ADJECTIVE SATELLITE, r => ADVERB
	w_cnt         int   // 2-digit hexadecimal integer
	lemmas        []*lemma
	p_cnt         int
	ptrs          []*synsetPtr
	gloss         []byte
}

func dataLookup(fh *os.File, offset int64) []byte {
	_, err := fh.Seek(offset, os.SEEK_SET)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, BUFFSIZE) // initial size of the buffer is 3kb
	line := make([]byte, 0, BUFFSIZE)
	prevLen := 0
	for {
		prevLen = len(line)
		n, err := fh.Read(buffer) // we read the next 3kb (or less)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		line = append(line, buffer[:]...)
		until := bytes.IndexByte(buffer, '\n')
		if until > 0 { // We have a full line
			return line[:prevLen+until]
		}
		if err == io.EOF || n < BUFFSIZE {
			return line
		}
	}
	log.Fatal("We can't be here")
	return nil
}

func GetRelation(dataLine []byte, symbol []byte) ([]synsetPtr, error) {
	ptrs := make([]synsetPtr, 0, 2) // larger cap?
	for {
		posInLine := bytes.Index(dataLine, symbol)
		if posInLine < 0 { // no occurrence
			return ptrs, nil
		}
		
		ptr, newPos, err := nextPtr(dataLine, posInLine)
		if err != nil {
			return nil, err
		}
		dataLine = dataLine[newPos:]
		ptrs = append(ptrs, *ptr)
	}
	log.Fatal("We can't be here")
	return nil, nil // never reached!!
}

func parseDataLine(dataLine []byte) (*dataData, error) {
	data := &dataData{}
	var err error

	// gloss
	lastIndex := bytes.LastIndex(dataLine, []byte{'|', ' '})
	if lastIndex == -1 {
		return nil, errors.New(`No gloss delimiter found "| " in line`)
	}
	data.gloss = dataLine[lastIndex+2:]

	// synset_offset  --- not used
	synsetOffsetBytes := dataLine[:8]
	fmt.Printf("synsetOffsetBytes: [%s]\n", synsetOffsetBytes)
	synset_offset, err := strconv.Atoi64(string(synsetOffsetBytes))
	if err != nil {
		return nil, err
	}
	data.synset_offset = synset_offset

	// lex_filenum  --- not used
	lexFilenumBytes := dataLine[9:11]
	fmt.Printf("lexFilenumBytes: [%s]\n", lexFilenumBytes)
	lexFilenum, err := strconv.Atoi(string(lexFilenumBytes))
	if err != nil {
		return nil, err
	}
	data.lex_filenum = lexFilenum

	// ss_type
	switch ss_type := dataLine[12]; {
	case ss_type == 'n' ||
		ss_type == 'v' ||
		ss_type == 'a' ||
		ss_type == 's' ||
		ss_type == 'r':
		data.ss_type = ss_type
	default:
		return nil, errors.New(fmt.Sprintf("Invalid ss_type: %c\n", ss_type))
	}

	// w_cnt
	w_cntBytes := dataLine[14:16]
	fmt.Printf("w_cntBytes: [%s]\n", w_cntBytes)
	w_cnt, err := x2i(w_cntBytes)
	if err != nil {
		return nil, err
	}
	data.w_cnt = w_cnt

	// lemmas
	lemmas := make([]*lemma, w_cnt)
	fromPos := 17
	for i := 0; i < w_cnt; i++ {
		nextLemma, posInLine, err := nextSense(dataLine, fromPos)
		if err != nil {
			return nil, err
		}
		fromPos = posInLine
		lemmas[i] = nextLemma
	}
	data.lemmas = lemmas

	// p_cnt
	p_cntBytes := dataLine[fromPos:fromPos+3]
	fmt.Printf("p_cnt: [%s]\n", p_cntBytes)
	p_cnt, err := strconv.Atoi(string(p_cntBytes))
	if err != nil {
		return nil, err
	}
	data.p_cnt = p_cnt

	// ptrs
	ptrs := make([]*synsetPtr, p_cnt)
	fromPos = fromPos + 4
	for i := 0; i < p_cnt; i++ {
		nextPtr, posInLine, err := nextPtr(dataLine, fromPos)
		if err != nil {
			return nil, err
		}
		fromPos = posInLine
		ptrs[i] = nextPtr
	}
	data.ptrs = ptrs

	//frames: In data.verb only, a list of numbers corresponding to the generic verb sentence frames for word s in the synset. frames is of the form:

	return data, nil
}

func nextSense(line []byte, pos int) (*lemma, int, error) {
	lemma := &lemma{}
	acc := make([]byte, 30)
	from := pos
	for i, ch := range line[pos:] {
		if ch == ' ' {
			lemma.word = acc
			from += i + 1
			break
		}
		acc[i] = ch
	}
	xval := line[from]
	ival, ok := fromHexChar(xval)
	if !ok {
		return nil, 0, errors.New(fmt.Sprintf("Invalid hex byte (%c)", xval))
	}
	lemma.lex_id = int(ival)

	return lemma, (from + 2), nil
}

func nextPtr(line []byte, pos int) (*synsetPtr, int, error) {
	ptr := &synsetPtr{}
	acc := make([]byte, 0, 2)
	from := pos
	for i, ch := range line[pos:] {
		if ch == ' ' {
			ptr.symbol = acc
			from += i + 1
			break
		}
		acc = append(acc, ch) // at most 2 digits in symbol
	}
	offsetBytes := line[from:from+8]
	offset, err := strconv.Atoi64(string(offsetBytes))
	if err != nil {
		return nil, 0, err
	}
	ptr.offset = offset
	from = from + 9

	ptrpos := line[from]
	ptr.pos = ptrpos

	from = from+2
	sourceBytes := line[from:from+2]
	source, err := x2i(sourceBytes)
	if err != nil {
		return nil, 0, err
	}
	ptr.source = source

	from = from+2
	targetBytes := line[from:from+2]
	target, err := x2i(targetBytes)
	if err != nil {
		return nil, 0, err
	}
	ptr.target = target

	return ptr, from + 3, nil
}

// Utility function to convert a 2-digit hexadecimal number to int
func x2i(src []byte) (int, error) {
	if len(src) != 2 {
		return 0, errors.New(fmt.Sprintf("Only 2-bytes hexadecimal integers allowed (passed [%s])\n", src))
	}
	d1, ok := fromHexChar(src[0])
	if !ok {
		return 0, hex.InvalidHexCharError(src[0])
	}
	d2, ok := fromHexChar(src[1])
	if !ok {
		return 0, hex.InvalidHexCharError(src[1])
	}
	val := int(d1)*16 + int(d2)
	return val, nil
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
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

