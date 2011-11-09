package main

import (
	"fmt"
	"./gown"
)

func main() {
	// db
	wndb := gown.New()

	// Get all the senses of a word
	synset := wndb.Get('dog')  // synset is of type []sense
	for _, sense := range synset {
		fmt.Println(sense)
	}	
}
