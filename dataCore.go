package main

import (
	"math/rand"
	"strings"
	"unicode"
)

func sanatizeString(input string) string {
	prev := string(input)
	prefiltered := strings.Map(func(r rune) rune {
		if unicode.IsNumber(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, prev)

	return prefiltered
}

func shuffleString(s string) string {
	runes := []rune(s)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

func createDada(input []string, poetLines int) []string {
	var asOneString string
	for _, a := range input {
		asOneString += " " + a
	}

	words := shuffleWords(getWords(sanatizeString(asOneString)))
	var lines = make([]string, poetLines)

	globalWordCount := 0
	wordCountForLines := getRandomWordCountForLines(poetLines)
	for currentLine, numWordsForThisLine := range wordCountForLines {
		for lineWords := 1; lineWords <= numWordsForThisLine; lineWords++ {
			lines[currentLine] = lines[currentLine] + words[globalWordCount%len(words)] + " "
			globalWordCount += 1
		}
	}
	return lines
}

func getRandomWordCountForLines(poetLines int) []int {
	var lineWordCounts = make([]int, poetLines)
	for i := 0; i < poetLines; i++ {
		lineWordCounts[i] = getRandomWordCountForLine()
	}

	return lineWordCounts
}

func getRandomWordCountForLine() int {
	minWordsPerLine := 1
	maxWordsPerLine := 10
	return rand.Intn(maxWordsPerLine-minWordsPerLine) + minWordsPerLine
}

func getWords(input string) []string {
	return strings.Fields(input)
}

func shuffleWords(input []string) []string {
	retVal := input

	rand.Shuffle(len(retVal), func(i, j int) {
		retVal[i], retVal[j] = retVal[j], retVal[i]
	})

	return retVal
}
