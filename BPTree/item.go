package BPTree

type item struct {
	key   int
	value int
}

func newItem(key, value int) item {
	return item{
		key:   key,
		value: value,
	}
}
