package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func handleSplit(input string) []string {
	if input == "" {
		return []string{}
	}
	var args []string

	word := ""
	quotedWord := ""
	activeQuote := ""

	for _, runeValue := range input {
		currentChar := fmt.Sprintf("%c", runeValue)

		// handle quote
		charIsQuote := currentChar == "'" || currentChar == `"`
		if charIsQuote {
			if activeQuote == "" {
				quotedWord = ""
				activeQuote = currentChar
				continue
			}

			charIsActiveQuote := currentChar == activeQuote
			if !charIsActiveQuote {
				quotedWord += currentChar
				continue
			}

			word += quotedWord
			quotedWord = ""
			activeQuote = ""
			continue
		}

		if activeQuote != "" {
			quotedWord += currentChar
			continue
		}
		
		// handle whitespace
		isCurrentCharWhiteSpace := unicode.IsSpace(runeValue)
		if isCurrentCharWhiteSpace {
			wordCount := utf8.RuneCountInString(word)
			if wordCount == 0 {
				continue
			}
			args = append(args, word)
			word = ""
			continue
		}
		
		word += currentChar
	}

	quotedWordCount := utf8.RuneCountInString(word)
	if quotedWordCount > 0 {
		word += quotedWord
	}

	wordCount := utf8.RuneCountInString(word)
	if wordCount > 0 {
		args = append(args, word)
		word = ""
	}
	return args
}