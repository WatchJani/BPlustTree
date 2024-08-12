package BPTree

import (
	"math/rand"
	"testing"
)

func TestInsert(t *testing.T) {
	for range 200 {
		tree := New(5)
		treeKey := map[int]struct{}{}

		for range 2000 {
			num := rand.Intn(200000)
			tree.Insert(num, 52)
			treeKey[num] = struct{}{}
		}

		leafKeyNumber := tree.TestFunc()

		if realNumber := len(treeKey); leafKeyNumber != realNumber {
			t.Errorf("real number of key: %d | tree number of key %d", realNumber, leafKeyNumber)
		}
	}
}

func TestDelete(t *testing.T) {
	tree := New(3)

	key := make([]int, 30)
	for index := range 30 {
		num := rand.Intn(200000)
		tree.Insert(num, 52)
		key[index] = num
	}

	for _, key := range key {
		tree.Delete(key)
	}

	if tree.root.pointer != 0 {
		t.Errorf("all elements is not deleted from tree")
	}
}

func BenchmarkInsertBPTree(b *testing.B) {
	b.StopTimer()

	tree := New(99)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(rand.Intn(100000), 5)
	}
}
