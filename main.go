package main

import (
	"fmt"

	BPTree "github.com/WatchJani/BPlustTree/btree"
)

func main() {
	tree := BPTree.New[int, int](10)

	store := make([]int, 100)

	for index := range store {
		store[index] = index
	}

	for _, key := range store {
		tree.Insert(key, 52)
	}

	fmt.Println(tree.Range(11, 11))

}
