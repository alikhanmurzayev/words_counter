package main

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestCompareWithStdMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const wordLen = 4

	var (
		stdMap        = make(map[string]int)
		htMap         = make(map[string]int)
		wordsCountMap WordsCountMap
	)

	for i := 0; i < 100000; i++ {
		word := randomWord(wordLen)
		stdMap[string(word)]++
		wordsCountMap.Store(word, wordsCountMap.Load(word)+1)
	}

	wordsCountMap.Range(func(word []byte, count int) (continueRange bool) {
		htMap[string(word)] = count
		return true
	})

	assert.Equal(t, stdMap, htMap)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomWord(length int) (word []byte) {
	for i := 0; i < length; i++ {
		word = append(word, alphabet[rand.Intn(len(alphabet))])
	}
	return
}
