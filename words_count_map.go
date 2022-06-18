package main

import (
	"github.com/alikhanmurzayev/words_counter/hash_table"
)

type WordsCountMap struct {
	ht hash_table.HashTable
}

func (m *WordsCountMap) Store(word []byte, count int) {
	m.ht.Store(ByteSlice(word), count)
}

func (m *WordsCountMap) Load(word []byte) int {
	count, ok := m.ht.Load(ByteSlice(word))
	if !ok {
		return 0
	}
	return count.(int)
}

func (m *WordsCountMap) Range(f func(word []byte, count int) (continueRange bool)) {
	m.ht.Range(func(key hash_table.Key, value interface{}) bool {
		return f(key.(ByteSlice), value.(int))
	})
}
