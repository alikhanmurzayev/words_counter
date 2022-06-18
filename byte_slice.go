package main

import "hash/fnv"

type ByteSlice []byte

func (bs ByteSlice) Hash() int {
	h := fnv.New32()
	_, _ = h.Write(bs)
	return int(h.Sum32())
}

func (bs ByteSlice) Interface() interface{} {
	return []byte(bs)
}
