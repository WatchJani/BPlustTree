package main

import (
	"fmt"
	"log"
	"root/BPTree"
)

func main() {
	tree := BPTree.New(5)

	myMap := map[int]struct{}{}

	key := []int{12, 26, 10, 21, 24, 12, 13, 22, 29, 27, 25, 9, 20, 5, 6, 8, 23, 26, 10, 7, 1, 28, 28, 10, 4, 23, 0, 20, 26, 23}
	fmt.Println(len(key))

	for _, num := range key {
		tree.Insert(num, 52)
		myMap[num] = struct{}{}
	}

	tree.TestFunc()
	fmt.Println(len(myMap))

	fmt.Println("==========================================================")

	if err := tree.Delete(26); err != nil {
		log.Println(err)
		return
	}

	if err := tree.Delete(29); err != nil {
		log.Println(err)
		return
	}

	if err := tree.Delete(25); err != nil {
		log.Println(err)
		return
	}

	//
	if err := tree.Delete(21); err != nil {
		log.Println(err)
		return
	}

	// if err := tree.Delete(0); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	//key not founded, check what is problem
	// if err := tree.Delete(23); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if err := tree.Delete(21); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	tree.TestFunc()

}
