package command

import (
	"fmt"
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
	v, err := cmd.Evaluate()
	fmt.Println("SYNC OUTPUT: ", v)
	return err
}
