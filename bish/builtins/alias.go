package builtins

import (
	"errors"

	"github.com/dalloriam/bish/bish/constants"

	"github.com/dalloriam/bish/bish/state"
)

// Attributes of the alias commands.
const (
	aliasName = "alias"
)

func init() {
	registry[aliasName] = alias
}

func alias(ctx *state.State, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid alias syntax")
	}
	ctx.SetKey(constants.AliasKey, args[0], args[1:])

	return "", nil
}
