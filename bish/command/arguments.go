package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var keepChars = []rune{'(', ')', '[', ']', '!', ';', '|', '>', '<'}
var keepCharMap map[rune]struct{}

var nonSeparators = []rune{'$', '-', '_', '/', '.', '~'}
var nonSeparatorMap map[rune]struct{}

func init() {
	keepCharMap = make(map[rune]struct{})
	nonSeparatorMap = make(map[rune]struct{})

	for _, r := range keepChars {
		keepCharMap[r] = struct{}{}
	}

	for _, r := range nonSeparators {
		nonSeparatorMap[r] = struct{}{}
	}
}

func isSeparator(ch rune) bool {
	_, ok := nonSeparatorMap[ch]
	return !(unicode.IsLetter(ch) || unicode.IsDigit(ch) || ok)
}

func isKeepChar(ch rune) (ok bool) {
	_, ok = keepCharMap[ch]
	return
}

func ParseArguments(in string) ([]string, error) {
	var tokens []string

	tokenStart := 0
	charCount := len(in)
	lookingForTokenEnd := false

	inLiteralBlock := false
	var literalChar rune

	inEscape := false

	for i, currentChar := range in {
		fmt.Println(string(currentChar))

		if currentChar == '\\' {
			fmt.Println("Escaping")
			inEscape = true
		} else if inEscape {
			fmt.Println("Stopping escape")
			inEscape = false
			continue
		}

		if inLiteralBlock {
			if currentChar == literalChar {
				inLiteralBlock = false
				tokens = append(tokens, in[tokenStart:i])
				lookingForTokenEnd = true
			}
			continue
		}

		if currentChar == '"' || currentChar == '\'' {
			inLiteralBlock = true
			literalChar = currentChar
			tokenStart = i + 1
			continue
		}

		if isSeparator(currentChar) {
			if !lookingForTokenEnd {
				if i >= (tokenStart + 1) {
					tokens = append(tokens, in[tokenStart:i])
				}
				lookingForTokenEnd = true
			}
			if isKeepChar(currentChar) {
				tokens = append(tokens, string(currentChar))
			}
		} else if lookingForTokenEnd {
			if isKeepChar(currentChar) {
				tokens = append(tokens, string(currentChar))
			}
			tokenStart = i
			lookingForTokenEnd = false
		}

	}
	if !lookingForTokenEnd && charCount-1 >= (tokenStart+1) {
		tokens = append(tokens, in[tokenStart:charCount])
	}

	r, _ := json.Marshal(tokens)
	fmt.Println("parsed: ", string(r))

	return tokens, nil
}

// SubstituteArguments performs env substitution.
func SubstituteArguments(in []string) ([]string, error) {
	var vals []string

	for i, val := range in {

		if strings.HasPrefix(val, "$") {
			key := val[1:]
			in[i] = os.Getenv(key)
			val = os.Getenv(key)
			vals = append(vals, val)
		} else {
			vals = append(vals, val)
		}

	}
	r, _ := json.Marshal(vals)
	fmt.Println("subbed: ", string(r))
	return vals, nil
}
