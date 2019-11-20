package command

import (
	"strings"
)

func DoCommand(input string) error {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input to separate the command and the arguments.
	// TODO: Fancier argument parsing.
	args, err := ParseArguments(input)
	if err != nil {
		return err
	}

	planner := NewExecutionPlanner(args)
	cmd, err := planner.Command(true)
	if err != nil {
		return err
	}
	_, err = cmd.Evaluate()
	return err
}
