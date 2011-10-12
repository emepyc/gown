package binsearch

import (
//	"fmt"
//	"sort"
)

type Interface interface {
	LessOrEqual(i, j int) bool
}

func IsSorted (data Interface, n int) bool {
	for i:=0; i<n; i++ {
		if ! data.LessOrEqual(i,i+1) {
			return false
		}
	}
	return true
}



