package hash_table

import "math"

const initialBucketsGroupSize = 4

type HashTable struct {
	bucketsGroups []bucketsGroup
}

func (ht *HashTable) Store(key Key, value interface{}) {
	if key == nil {
		panic("cannot store nil key")
	}

	if len(ht.bucketsGroups) == 0 {
		ht.grow()
	}

	// Store key in the last bucketsGroup
	lastGroupIdx := len(ht.bucketsGroups) - 1
	evacuateBucketIdx := ht.bucketsGroups[lastGroupIdx].store(key, value, false)

	// Evacuate overflowed bucket
	if evacuateBucketIdx >= 0 {
		ht.grow()
		newLastGroupIdx := len(ht.bucketsGroups) - 1
		ht.bucketsGroups[lastGroupIdx].evacuate(evacuateBucketIdx, &ht.bucketsGroups[newLastGroupIdx], nil)
	}

	ht.evacuateElementsFromOldBuckets(key)
}

func (ht *HashTable) Load(key Key) (interface{}, bool) {
	for groupIdx := len(ht.bucketsGroups) - 1; groupIdx >= 0; groupIdx-- {

		bucketIdx := ht.bucketsGroups[groupIdx].getBucketIdx(key)
		elementIdx := ht.bucketsGroups[groupIdx].buckets[bucketIdx].search(key)

		if elementIdx >= 0 {
			return ht.bucketsGroups[groupIdx].buckets[bucketIdx].elements[elementIdx].value, true
		}
	}
	return nil, false
}

func (ht *HashTable) Range(f func(key Key, value interface{}) (continueRange bool)) {
	for _, group := range ht.bucketsGroups {
		for _, theBucket := range group.buckets {
			if theBucket.evacuate {
				continue
			}
			for _, elem := range theBucket.elements {
				if !f(elem.key, elem.value) {
					return
				}
			}
		}
	}
}

func (ht *HashTable) evacuateElementsFromOldBuckets(key Key) {
	// Iterate over old bucketsGroups
	for groupIdx := 0; groupIdx < len(ht.bucketsGroups)-1; groupIdx++ {

		group := ht.bucketsGroups[groupIdx]
		bucketIdx := group.getBucketIdx(key)

		if !group.buckets[bucketIdx].evacuate && group.buckets[bucketIdx].search(key) >= 0 {
			// Evacuate bucket for a specified key to the last bucketsGroup
			group.buckets[bucketIdx].evacuate = true
			group.evacuate(bucketIdx, &ht.bucketsGroups[len(ht.bucketsGroups)-1], key)
		}
	}
}

// grow adds a new bucketsGroup into a ht.
// The size of first bucketsGroup will be initialBucketsGroupSize.
// Each next bucketsGroup's size will be doubled.
func (ht *HashTable) grow() {
	size := initialBucketsGroupSize * power2(len(ht.bucketsGroups))
	ht.bucketsGroups = append(ht.bucketsGroups, bucketsGroup{
		buckets: make([]bucket, size),
	})
}

func power2(n int) int {
	return int(math.Pow(2, float64(n)))
}
