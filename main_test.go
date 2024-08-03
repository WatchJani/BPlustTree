package main

import (
	"math/rand"
	"root/BPTree"
	"slices"
	"testing"
)

func TestDeleteElements(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	numberDeletion := DeleteElement(list, 0, 5)
	expected := []int{6, 7, 8, 9, 10}
	get := list[:len(list)-numberDeletion]

	if slices.Compare(get, expected) != 0 {
		t.Errorf("expected slice: %v | get slice %v", expected, get)
	}
}

func TestMigrateElement(t *testing.T) {
	list := []int{1, 2, 3, 4, 0, 0, 0, 0, 0, 0}
	migrate := []int{6, 6, 7, 8, 9}

	MigrateElement(list, migrate, 4)
	expected := []int{1, 2, 3, 4, 6, 6, 7, 8, 9, 0}

	if slices.Compare(list, expected) != 0 {
		t.Errorf("expected slice: %v | get slice %v", expected, list)
	}
}

func SearchTestData() ([]int, []struct {
	key              int
	expectedPosition int
}) {
	list := []int{5, 15, 150, 154, 199, 222, 315, 451, 458, 500}

	test := []struct {
		key              int
		expectedPosition int
	}{{
		key:              3,
		expectedPosition: 0,
	}, {
		key:              1511,
		expectedPosition: 10,
	}, {
		key:              256,
		expectedPosition: 6,
	}, {
		key:              15,
		expectedPosition: 1,
	}}

	return list, test
}

func TestLinearSearchLess(t *testing.T) {
	list, test := SearchTestData()

	for index, test := range test {
		if get := LinerSearchLess(list, test.key); get != test.expectedPosition {
			t.Errorf("%d | expected: %d | get: %d", index, test.expectedPosition, get)
		}
	}
}

func TestBinarySearchLess(t *testing.T) {
	list, test := SearchTestData()
	for index, test := range test {
		if get := BinarySearchLess(list, test.key); get != test.expectedPosition {
			t.Errorf("%d | expected: %d | get: %d", index, test.expectedPosition, get)
		}
	}
}

func TestInsert(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	Insert(list, 1, 5)

	expected := []int{1, 2, 3, 4, 5, 1, 2, 5, 5, 6}

	if slices.Compare(list, expected) != 0 {
		t.Errorf("expected: %v | get: %v", expected, list)
	}
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	forInsert := 1
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Insert(list, forInsert, 5)
	}
}

func BenchmarkInsertBPTree(b *testing.B) {
	b.StopTimer()

	tree := BPTree.New(1000)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(rand.Intn(150), 5)
	}
}
