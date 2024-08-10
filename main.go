package main

import (
	"fmt"
	"log"
	"root/BPTree"
)

func main() {
	tree := BPTree.New(5)

	myMap := map[int]struct{}{}

	//1. index error
	//2. infinite loop

	key := []int{1, 8, 8, 20, 18, 23, 8, 28, 20, 10, 19, 10, 5, 16, 6, 23, 4, 11, 3, 17, 8, 14, 16, 21, 12, 29, 1, 25, 17, 3}

	//1, 8, 20, 18, 23, 28, 10, 19, 5, 16, 6, 4, 11, 3, 17, 14, 21, 12, 29, 25

	for _, num := range key {
		tree.Insert(num, 52)
		myMap[num] = struct{}{}
	}

	fmt.Println("==========================================================")

	if err := tree.Delete(1); err != nil {
		log.Println(err)
		return
	}

	if err := tree.Delete(8); err != nil {
		log.Println(err)
		return
	}

	if err := tree.Delete(20); err != nil {
		log.Println(err)
		return
	}

	// if err := tree.Delete(18); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if err := tree.Delete(23); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if err := tree.Delete(28); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if err := tree.Delete(10); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	tree.TestFunc()
	root := tree.GetRoot()
	fmt.Println(root.Children[1].Children[3])
}
