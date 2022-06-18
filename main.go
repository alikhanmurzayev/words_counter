package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"unicode"
)

const wordsNum = 20

func main() {
	if len(os.Args) < 2 {
		log.Fatal("file path not provided")
	}

	mp, err := getWordsCountsMap(os.Args[1])
	if err != nil {
		log.Fatalf("getWordsCountsMap: %s", err)
	}

	counts := getWordsCountsSlice(mp)
	sortCounts(counts)

	for i := 0; i < min(wordsNum, len(counts)); i++ {
		fmt.Printf("%6d %s\n", counts[i].count, counts[i].word)
	}
}

func sortCounts(counts []wordCount) {
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].count == counts[j].count {
			// Sort lexicographically
			return bytes.Compare(counts[i].word, counts[j].word) > 0
		}
		return counts[i].count > counts[j].count
	})
}

func getWordsCountsSlice(mp WordsCountMap) (counts []wordCount) {
	mp.Range(func(word []byte, count int) (continueRange bool) {
		counts = append(counts, wordCount{word: word, count: count})
		return true
	})
	return
}

func getWordsCountsMap(filePath string) (mp WordsCountMap, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("could not open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := cleanWord(scanner.Bytes())
		if len(word) > 0 {
			mp.Store(word, mp.Load(word)+1)
		}
	}

	return
}

func cleanWord(word []byte) []byte {
	word = bytes.TrimFunc(word, func(r rune) bool { return !unicode.IsLetter(r) })
	word = bytes.ToLower(word)
	return word
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type wordCount struct {
	word  []byte
	count int
}
