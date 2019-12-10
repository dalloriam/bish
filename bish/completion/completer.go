package completion

import (
	"strings"

	"github.com/dalloriam/bish/bish/command"
)

// Completer applies all autocompletion.
// It conforms to the Autocompleter interface of readline.
type Completer struct {
	ctx ShellContxt
}

type ShellContxt interface {
	GetKey(domain string, key string) (interface{}, bool)
	SetKey(domain string, key string, value interface{})
}

// New returns a new completion engine.
func New(ctx ShellContxt) *Completer {
	return &Completer{ctx: ctx}
}

func (c *Completer) findCandidates(lineStr string, currentTok string) []string {
	// Find completion candidates
	if len(strings.Split(lineStr, " ")) == 1 {
		// Attempt executable completion
		return listExecutables()
	} else {
		// TODO: Path completion + context-sensitive completion.
		return listPathOptions(currentTok)
	}

}

func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line)

	var candidates [][]rune
	commonGround := len(line)

	parsedTrim, err := command.ParseArguments(string(line[:pos]))
	if err == nil {
		currentToks, err := command.ProcessArg(parsedTrim[len(parsedTrim)-1], c.ctx)
		if err != nil || len(currentToks) != 1 {
			return candidates, commonGround
		}
		currentTok := currentToks[0]
		catalog := c.findCandidates(lineStr[:pos], currentTok)

		for _, trial := range catalog {
			if strings.HasPrefix(trial, currentTok) {
				candidates = append(candidates, []rune(strings.TrimPrefix(trial, currentTok)))
			}
		}
	}

	if len(candidates) == 0 {
		candidates = append(candidates, []rune{'\t'})
		commonGround = 0
	}

	return candidates, commonGround
}
