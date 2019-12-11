package builtins

import (
	"errors"

	"github.com/dalloriam/bish/bish/state"
)

// Attributes of the alias commands.
const (
	AliasName       = "alias"
	AliasContextKey = "alias"
)

func Alias(ctx *state.State, args []string) error {
	if len(args) < 2 {
		return errors.New("invalid alias syntax")
	}
	ctx.SetKey("alias", args[0], args[1:])

	return nil
}
