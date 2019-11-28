package main

import "strings"

type Complete struct {
	candidates []string
}

func newComplete() *Complete {
	return &Complete{
		candidates: []string{"go run", "github", "git checkout"},
	}
}

func (c *Complete) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var candidates [][]rune

	for _, c := range c.candidates {
		if strings.HasPrefix(c, string(line[:pos])) {
			candidates = append(candidates, []rune(c[pos:]))
		}
	}

	return candidates, len(line)
}
