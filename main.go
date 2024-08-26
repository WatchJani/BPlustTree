package main

import (
	"fmt"

	BPTree "github.com/WatchJani/BPlustTree/btree"
)

func main() {
	tree := BPTree.New[int, int](99)

	store := make([]int, 100)

	for index := range store {
		store[index] = index
	}

	for _, key := range store {
		tree.Insert(key, 52)
	}

	val, _ := tree.Find(63)
	fmt.Println(val)

	fmt.Println(tree.RangeUp(63, 85, ">"))
}
