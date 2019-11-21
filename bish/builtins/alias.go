package builtins

import "errors"

const (
	AliasName       = "alias"
	AliasContextKey = "alias"
)

type ShellContext interface {
	GetKey(domain string, key string) (interface{}, bool)
	SetKey(domain string, key string, value interface{})
}

func Alias(ctx ShellContext, args []string) error {
	if len(args) < 2 {
		return errors.New("invalid alias syntax")
	}
	ctx.SetKey("alias", args[0], args[1:])

	return nil
}
