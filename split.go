package main

import (
	"fmt"
	"slices"
	"unicode"
	"unicode/utf8"
)

type SplitLoopState struct {
	args []string
	word, quotedWord, activeQuote, redirectChar string
	backSlashActive bool
	redirectCharIndex int
}

type SplitResult struct {
	fields []string
	redirectCharIndex int 
	redirectChar string
}

// handleSplit handles the splitting of the whole string the user types as their command.
// The string is split based on a set of rules.
// Backslashes escape the next character except if they are in a single quote (where they are taken literally and treated as just another character).
// Quotes remove the meaning from special characters, meaning they are parsed literally. Whitespaces, backslashes .e.t.c are treated as just another character.
// Consecutive whitespaces are reduced to just one (if they are not in quotes).
func handleSplit(input string) SplitResult {
	if input == "" {
		return SplitResult{
			fields: []string{},
			redirectCharIndex: -1,
			redirectChar: "",
		}
	}

	state := SplitLoopState{
		args: []string{},	
		word: "",
		quotedWord: "",
		activeQuote: "",
		backSlashActive: false,	
		redirectCharIndex: -1,
		redirectChar: "",
	}

	for _, runeValue := range input {
		currentChar := fmt.Sprintf("%c", runeValue)
		charIsQuote := currentChar == "'" || currentChar == `"`
		quoteIsOpen := state.activeQuote != ""
		charIsBackSlash := currentChar == `\`
		charIsWhiteSpace := unicode.IsSpace(runeValue)
		charIsRedirect := currentChar == ">"

		// the next character after a backslash
		if state.backSlashActive {
			state.handleBackSlashActive(currentChar, runeValue)
			continue
		}

		// handle quote
		if charIsQuote {
			state.handleCharIsQuote(currentChar)
			continue
		}

		// we are inside an open quote here (.e.g. "example)
		if quoteIsOpen {
			state.handleQuoteIsOpen(currentChar)
			continue
		}

		// handle backslash
		if charIsBackSlash {
			state.handleCharIsBackSlash()
			continue
		}

		// handle greater than
		if charIsRedirect {
			state.handleCharIsRedirect()
			continue
		}
		
		// handle whitespace
		if charIsWhiteSpace {
			state.handleCharIsWhiteSpace()
			continue
		}
		
		state.word += currentChar
	}

	state.finalizeSplit()
	// return state.args, state.redirectCharIndex
	return SplitResult{
		fields: state.args,
		redirectCharIndex: state.redirectCharIndex,
		redirectChar: state.redirectChar,
	}
}

func (state *SplitLoopState) handleBackSlashActive(currentChar string, runeValue rune) {
	quoteIsOpen := state.activeQuote != ""
	state.backSlashActive = false
	if !quoteIsOpen {
		state.word += currentChar
		return
	}
	isDoubleQuote := state.activeQuote == `"`
	charIsEscapable := currentChar == `"` || currentChar == `\`
	if !isDoubleQuote || !charIsEscapable {
		state.quotedWord += currentChar
		return
	}
	lastQuotedRune, size := utf8.DecodeLastRuneInString(state.quotedWord)
	if lastQuotedRune == utf8.RuneError && size == 0 { // quotedWord is empty
		state.quotedWord += currentChar
		return
	}
	state.quotedWord = state.quotedWord[:len(state.quotedWord)-size] + string(runeValue)  //replace the last character("\") with the escaped new character(.e.g. "\" or `"`)
}

func (state *SplitLoopState) handleCharIsQuote(currentChar string) {
	noQuoteOpen := state.activeQuote == ""
	if noQuoteOpen {
		state.quotedWord = ""
		state.activeQuote = currentChar
		return
	}

	charIsActiveQuote := currentChar == state.activeQuote
	if !charIsActiveQuote {
		state.quotedWord += currentChar
		return
	}

	state.word += state.quotedWord
	state.quotedWord = ""
	state.activeQuote = ""
}

func (state *SplitLoopState) handleQuoteIsOpen(currentChar string) {
	charIsBackSlash := currentChar == `\`
	isDoubleQuote := state.activeQuote == `"`
	if charIsBackSlash && isDoubleQuote {
		state.backSlashActive = true
	}
	state.quotedWord += currentChar
}

func (state *SplitLoopState) handleCharIsBackSlash() {
	state.backSlashActive = true
}

func (state *SplitLoopState) handleCharIsWhiteSpace() {
	wordCount := utf8.RuneCountInString(state.word)
	if wordCount == 0 {
		return
	}
	state.args = append(state.args, state.word)
	state.word = ""
}

func (state *SplitLoopState) finalizeSplit() {
	quotedWordCount := utf8.RuneCountInString(state.quotedWord)
	if quotedWordCount > 0 {
		state.word += state.quotedWord
		state.quotedWord = ""
	}
	wordCount := utf8.RuneCountInString(state.word)
	if wordCount > 0 {
		state.args = append(state.args, state.word)
		state.word = ""
	}
}

var supportedRedirectOperators = []string{"1", "2"}

func (state *SplitLoopState) handleCharIsRedirect() {
	lastWordRune, size := utf8.DecodeLastRuneInString(state.word)
	lastChar := ""
	lastCharIsSupported := false  // "1> or 2>"
	if size > 0 {
		lastChar = fmt.Sprintf("%c", lastWordRune)
		lastCharIsSupported = slices.Contains(supportedRedirectOperators, lastChar)
	}
	if lastCharIsSupported {
		state.redirectChar = lastChar
		state.word = state.word[:len(state.word)-size] //remove the last character, which is "1" or "2"
	}
	
	state.finalizeSplit()
	state.redirectCharIndex = len(state.args)
	state.args = append(state.args, ">")
}