package command

import (
	"errors"
	"os"
	"os/user"
	"regexp"
	"strings"
	"unicode"

	"github.com/dalloriam/bish/bish/state"

	"github.com/dalloriam/bish/bish/builtins"
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

func applySupportedEscapes(in string) string {
	in = strings.ReplaceAll(in, "\\n", "\n")
	return strings.ReplaceAll(in, "\\t", "\t")
}

// ParseArguments tokenizes the command according to shell syntax.
func ParseArguments(in string) ([]string, error) {
	// Pre-process input to convert supported escapes.
	// Eg. []rune{'\', 'n'} => []rune{'\n'}
	in = applySupportedEscapes(in)

	var tokens []string
	tokenStart := 0
	charCount := len(in)
	lookingForTokenEnd := false

	inLiteralBlock := false
	var literalChar rune

	for i := 0; i < len(in); i++ {
		currentChar := rune(in[i])

		if currentChar == '\\' && i+1 < len(in) {
			in = string(append([]rune(in[:i]), []rune(in[i+1:])...))
			continue
		}

		if inLiteralBlock {
			if currentChar == literalChar {
				inLiteralBlock = false
				tokens = append(tokens, in[tokenStart:i])
				lookingForTokenEnd = true
				if currentChar == '>' {
					// We're closing an expression block -- not shellcode.
					tokens = append(tokens, ">")
				}
			}
			continue
		}

		if currentChar == '<' {
			inLiteralBlock = true
			literalChar = '>'
			tokenStart = i + 1
			tokens = append(tokens, string(in[i]))
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
	if !lookingForTokenEnd && charCount >= (tokenStart+1) {
		tokens = append(tokens, in[tokenStart:charCount])
	}

	// TODO: Hide under debug setting
	//r, _ := json.Marshal(tokens)
	//fmt.Println("parsed: ", string(r))

	return tokens, nil
}

func substituteEnvironmentVariables(argument string) string {
	pattern := regexp.MustCompile(`\$[a-zA-Z0-9]+`)
	matches := pattern.FindAllString(argument, -1)
	for _, match := range matches {
		key := os.Getenv(match[1:])
		argument = strings.ReplaceAll(argument, match, key)
	}
	return argument
}

func substituteHome(arg string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(arg, "~", usr.HomeDir), nil
}

func substituteAliases(ctx *state.State, arg string) ([]string, error) {
	var out []string
	if v, ok := ctx.GetKey(builtins.AliasContextKey, arg); ok {
		if s, sOk := v.(string); sOk {
			out = append(out, s)
		} else if ss, ssOk := v.([]string); ssOk {
			out = append(out, ss...)
		} else {
			return nil, errors.New("invalid alias")
		}
	} else {
		out = append(out, arg)
	}

	return out, nil
}

// ProcessArg performs env substitution.
func ProcessArg(arg string, ctx *state.State) ([]string, error) {
	args, err := substituteAliases(ctx, arg)

	for i, arg := range args {
		arg = substituteEnvironmentVariables(arg)
		arg, err = substituteHome(arg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, err
}
