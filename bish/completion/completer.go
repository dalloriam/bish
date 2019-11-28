package completion

import (
	"strings"
)

// Completer applies all autocompletion.
// It conforms to the Autocompleter interface of readline.
type Completer struct {
}

// New returns a new completion engine.
func New() *Completer {
	return &Completer{}
}

func (c *Completer) findCandidates(lineStr string) []string {
	// Find completion candidates
	if len(strings.Split(lineStr, " ")) == 1 {
		// Attempt executable completion
		return listExecutables()
	} else {
		// TODO: Path completion + context-sensitive completion.
		return nil
	}

}

func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line)

	catalog := c.findCandidates(lineStr[:pos])

	var candidates [][]rune
	commonGround := len(line)

	for _, trial := range catalog {
		if strings.HasPrefix(trial, lineStr[:pos]) {
			candidates = append(candidates, []rune(trial[pos:]))
		}
	}

	if len(candidates) == 0 {
		candidates = append(candidates, []rune{'\t'})
		commonGround = 0
	}

	return candidates, commonGround
}
