package main

import (
	"fmt"
	"./binsearch"
)

type T []int

func (t T) LessOrEqual(i, j int) bool{
	return t[i] <= t[j]
}

func main() {
	data := []int{1,1,2,4}
	interf := binsearch.Interface(T(data))
	fmt.Println(interf)
	fmt.Println(binsearch.IsSorted(interf, 2))
}
