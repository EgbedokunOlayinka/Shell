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
	backSlashActive := false

	for _, runeValue := range input {
		currentChar := fmt.Sprintf("%c", runeValue)
		charIsQuote := currentChar == "'" || currentChar == `"`
		quoteIsOpen := activeQuote != ""
		charIsBackSlash := currentChar == `\`
		charIsWhiteSpace := unicode.IsSpace(runeValue)

		// the next character after a backslash
		if backSlashActive {
			backSlashActive = false
			if !quoteIsOpen {
				word += currentChar
				continue
			}
			isDoubleQuote := activeQuote == `"`
			charIsEscapable := currentChar == `"` || currentChar == `\`
			if !isDoubleQuote || !charIsEscapable {
				quotedWord += currentChar
				continue
			}
			lastQuotedRune, size := utf8.DecodeLastRuneInString(quotedWord)
			if lastQuotedRune == utf8.RuneError && size == 0 { // quotedWord is empty
				quotedWord += currentChar
				continue
			}
			quotedWord = quotedWord[:len(quotedWord)-size] + string(runeValue)  //replace the last character("\") with the escaped new character(.e.g. "\" or `"`)
			continue
		}

		// handle quote
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

		// we are inside an open quote here (.e.g. "example)
		if quoteIsOpen {
			isDoubleQuote := activeQuote == `"`
			if charIsBackSlash && isDoubleQuote {
				backSlashActive = true
			}
			quotedWord += currentChar
			continue
		}

		// handle backslash
		if charIsBackSlash {
			backSlashActive = true
			continue
		}
		
		// handle whitespace
		if charIsWhiteSpace {
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