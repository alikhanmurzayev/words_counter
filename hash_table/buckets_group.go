package hash_table

import "reflect"

type bucketsGroup struct {
	buckets []bucket
}

// store stores key and value in one of buckets according to key.Hash.
// If used bucket is not overflowed, then evacuateBucketIdx is -1.
// Non-negative evacuateBucketIdx indicates that growth of hash table is required
// and a specific bucket must be evacuated first
func (bg *bucketsGroup) store(key Key, value interface{}, ignoreEvacuation bool) (evacuateBucketIdx int) {
	idx := bg.getBucketIdx(key)
	bg.buckets[idx].store(key, value, ignoreEvacuation)
	if !ignoreEvacuation && bg.buckets[idx].evacuate {
		return idx
	}
	return -1
}

func (bg *bucketsGroup) getBucketIdx(key Key) int {
	return key.Hash() % len(bg.buckets)
}

// evacuate moves elements from one bucket to another bucketsGroup
// with a bigger size and redistributes elements among buckets.
func (bg *bucketsGroup) evacuate(fromBucketIdx int, toBucketsGroup *bucketsGroup, omitKey Key) {
	for i := range bg.buckets[fromBucketIdx].elements {

		key := bg.buckets[fromBucketIdx].elements[i].key
		if omitKey != nil && reflect.DeepEqual(key.Interface(), omitKey.Interface()) {
			continue
		}

		value := bg.buckets[fromBucketIdx].elements[i].value

		_ = toBucketsGroup.store(key, value, true)
	}

	// Delete evacuated elements from old bucket
	bg.buckets[fromBucketIdx].elements = nil
}
