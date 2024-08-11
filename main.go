package main

import (
	"fmt"
	"log"
	"math/rand"
	"root/BPTree"
)

func main() {
	tree := BPTree.New(5)

	myMap := map[int]struct{}{}

	//2. infinite loop

	key := make([]int, 30)

	for index := range key {
		num := rand.Intn(30)
		fmt.Println(num)
		key[index] = num
	}

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
	fmt.Println(root.Children[4])
}
