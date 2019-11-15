package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dalloriam/bish/bish"
)

func prompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("wduss@phoenix [%s]\n> ", cwd)
}

func main() {
	os.Setenv("TERM", "xterm-kitty")
	reader := bufio.NewReader(os.Stdin)

	for {
		// Display the prompt.
		fmt.Print(prompt())
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cmd, err := bish.ParseCommand(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err := cmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}
