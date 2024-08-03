package main

import (
	"math/rand"
	"root/BPTree"
)

func main() {
	// list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// DeleteElement(list, 0, 5)

	// fmt.Println(list)

	// list := []int{1, 2, 3, 4, 0, 0, 0, 0, 0, 0}
	// migrate := []int{6, 6, 7, 8, 9}
	// MigrateElement(list, migrate, 4)

	// fmt.Println(list)

	// list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Insert(list, []int{1, 2, 5, 5}, 5)

	// fmt.Println(list)

	tree := BPTree.New(5)
	for range 5 {
		tree.Insert(rand.Intn(100000), 52)
	}

	tree.TestFunc()
}

type SearchFn func([]int, int) int

func Insert(list []int, insert, position int) int {
	copy(list[position+1:], list[position:])
	return copy(list[position:], []int{insert})
}

func DeleteElement(list []int, position, deletion int) int {
	return copy(list[position:], list[position+deletion:])
}

func MigrateElement(list, migrateElement []int, position int) int {
	return copy(list[position:], migrateElement)
}

func LinerSearchLess(arr []int, key int) int {
	for index, value := range arr {
		if key <= value {
			return index
		}
	}

	return len(arr)
}

func BinarySearchLess(arr []int, target int) int {
	low, high := 0, len(arr)-1

	if target > arr[high] {
		return high + 1
	}

	for low <= high {
		mid := (low + high) / 2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low
}
