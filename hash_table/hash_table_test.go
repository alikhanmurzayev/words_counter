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

func TestCompareWithStdMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const wordLen = 3

	var (
		stdMap = make(map[string]int)
		htMap  = make(map[string]int)
		ht     HashTable
	)

	for i := 0; i < 100000; i++ {
		word := randomWord(wordLen)

		stdMap[string(word)]++

		countInterface, _ := ht.Load(keyType(word))
		count, _ := countInterface.(int)
		ht.Store(keyType(word), count+1)
	}

	ht.Range(func(key Key, value interface{}) (continueRange bool) {
		htMap[string(key.(keyType))] = value.(int)
		return true
	})

	assert.Equal(t, stdMap, htMap)

	t.Logf("number of bucketsGruops: %d", len(ht.bucketsGroups))

}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomWord(length int) (word []byte) {
	for i := 0; i < length; i++ {
		word = append(word, alphabet[rand.Intn(len(alphabet))])
	}
	return
}
