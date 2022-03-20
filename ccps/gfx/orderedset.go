package gfx

import (
	"errors"
)

type allocator struct {
	index  map[int]bool
	values []int
}

func makeAllocator(size int) *allocator {
	var alloc = allocator{}
	alloc.index = make(map[int]bool)

	for i := 0; i < size; i++ {
		alloc.index[i] = true
		alloc.values = append(alloc.values, i)
	}

	return &alloc
}

func (alloc *allocator) any() (int, error) {
	if len(alloc.values) == 0 {
		return 0, errors.New("Unable to allocate: Out of ROM")
	}

	allocated := alloc.values[0]
	alloc.values = alloc.values[1:]
	alloc.index[allocated] = false
	return allocated, nil
}

func (alloc *allocator) has(v int) bool {
	return alloc.index[v]
}

func (alloc *allocator) isEmpty() bool {
	return len(alloc.values) == 0
}
