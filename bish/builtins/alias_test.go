package builtins_test

import (
	"github.com/dalloriam/bish/bish/constants"
	"github.com/dalloriam/bish/bish/state"
)

func validateAlias(ctx *state.State, args []string) bool {
	if len(args) < 2 {
		return false
	}
	aliasName := args[1]

	val, ok := ctx.GetKey(constants.AliasKey, aliasName)
	if !ok {
		return false
	}

	if lst, vOk := val.([]string); vOk {
		expectedLst := args[2:]
		if len(lst) != len(expectedLst) {
			return false
		}
		for i := 0; i < len(lst); i++ {
			if lst[i] != expectedLst[i] {
				return false
			}
		}
		return true
	}
	return false
}
