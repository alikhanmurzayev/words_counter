package hash_table

import (
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type keyType []byte

func (k keyType) Hash() int {
	h := fnv.New32()
	_, _ = h.Write(k)
	return int(h.Sum32())
}

func (k keyType) Interface() interface{} {
	return []byte(k)
}

func TestHashTable(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const n = 10000

	var (
		hashTable       HashTable
		mp              = make(map[string]int)
		mpFromHashTable = make(map[string]int)
	)

	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		mp[key] = i
		hashTable.Store(keyType(key), i)
	}

	for iter := 0; iter < 10; iter++ {
		for i := 0; i < n; i++ {
			randVal := rand.Intn(1000)
			key := strconv.Itoa(i)

			mp[key] += randVal

			oldValInterface, _ := hashTable.Load(keyType(key))
			oldVal := oldValInterface.(int)
			newVal := oldVal + randVal
			hashTable.Store(keyType(key), newVal)
		}
	}

	hashTable.Range(func(key Key, value interface{}) (continueRange bool) {
		mpFromHashTable[string(key.(keyType))] = value.(int)
		return true
	})

	assert.Equal(t, mp, mpFromHashTable)

	t.Logf("number of bucketsGruops: %d", len(hashTable.bucketsGroups))

}
