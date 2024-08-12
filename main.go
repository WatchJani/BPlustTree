package main

import (
	"fmt"
	"log"
	"root/BPTree"
)

func main() {

	//need to fix delete with 4
	tree := BPTree.New[int, int](4)

	myMap := map[int]struct{}{}

	key := []int{8, 5, 15, 6, 16, 17, 3, 11, 4, 19, 13, 16, 13, 10, 15, 18, 9, 1, 5, 4}

	// 8, 5, 15, 6, 16, 17, 3, 11, 4, 19, 13, 10, 18, 9, 1

	// for index := range key {
	// 	num := rand.Intn(20)
	// 	fmt.Println(num)
	// 	key[index] = num
	// }

	for _, num := range key {
		tree.Insert(num, 52)
		myMap[num] = struct{}{}
	}

	for _, key := range key {
		if err := tree.Delete(key); err != nil {
			log.Println(err)
		}
	}

	fmt.Println("==========================================================")
	tree.TestFunc()
	root := tree.GetRoot()
	fmt.Println(root.Children[1])
}
