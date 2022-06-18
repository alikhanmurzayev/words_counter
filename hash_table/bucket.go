package hash_table

import "reflect"

const maxBucketSize = 7

type bucket struct {
	evacuate bool
	elements []bucketElement
}

type bucketElement struct {
	key   Key
	value interface{}
}

func (b *bucket) store(key Key, value interface{}, ignoreEvacuation bool) {
	idx := b.search(key)
	if idx >= 0 {
		b.elements[idx].value = value
		return
	}

	// Insert new key
	b.elements = append(b.elements, bucketElement{
		key:   key,
		value: value,
	})

	if !ignoreEvacuation && len(b.elements) > maxBucketSize {
		// Bucket if overflowed, need to evacuate
		b.evacuate = true
	}
}

// search key in a bucket. If key is present, then returns its index.
// If key is not present, then return -1
func (b *bucket) search(key Key) int {
	for i := range b.elements {
		if reflect.DeepEqual(b.elements[i].key.Interface(), key.Interface()) {
			return i
		}
	}
	return -1
}
