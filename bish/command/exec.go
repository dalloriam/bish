package command

import (
	"errors"
	"fmt"
	"os"
)

type ExecutionPlanner struct {
	Args []string

	idx        int
	currentTok *string
	nextTok    *string
	done       bool
}

func NewExecutionPlanner(args []string) *ExecutionPlanner {
	e := &ExecutionPlanner{Args: args}

	e.advance()

	return e
}

func (p *ExecutionPlanner) advance() {
	p.currentTok = p.nextTok
	if p.idx < len(p.Args) {
		p.nextTok = &p.Args[p.idx]
	} else {
		p.done = true
	}
	p.idx++
}

func (p *ExecutionPlanner) accept(char string) bool {
	if *p.nextTok == char {
		p.advance()
		return true
	}
	return false
}

func (p *ExecutionPlanner) Argument() (Argument, error) {
	if p.accept("(") {
		subcmd, err := p.Command(false)
		if err != nil {
			return nil, err
		}
		fmt.Println("Subcommand: ", subcmd)
		if !p.accept(")") {
			return nil, errors.New("no closing paren")
		}
		return subcmd, nil

	} else if p.accept("<") {
		p.advance()
		subcmd := &PythonArgument{*p.currentTok}
		if !p.accept(">") {
			return nil, errors.New("no closing angle bracket")
		}
		return subcmd, nil

	} else {
		p.advance()
		arg, err := ProcessArg(*p.currentTok)
		if err != nil {
			return nil, err
		}
		return &StringArgument{arg}, nil
	}
}

func (p *ExecutionPlanner) Command(topLevel bool) (Command, error) {
	var argumentBuffer []Argument
	for {

		arg, err := p.Argument()
		if err != nil {
			return nil, err
		}
		argumentBuffer = append(argumentBuffer, arg)

		if p.done || *p.nextTok == ")" || *p.nextTok == ">" {
			baseCmd := &CommandTree{
				Args: argumentBuffer,
			}
			if topLevel {
				// TODO: Plug this to shell settings.
				baseCmd.Shell = true
				baseCmd.StdOut = os.Stdout
				baseCmd.StdErr = os.Stderr
				baseCmd.StdIn = os.Stdin
			}
			return baseCmd, nil
		} else if p.accept("|") {
			fmt.Println("Executing Pipe")
			aCmd := &CommandTree{
				Args: argumentBuffer,
				StdOut: os.Stdout,
				StdErr: os.Stderr,
				StdIn: os.Stdin,
				Shell: true,
			}
			bCmd, err := p.Command(true)
			if err != nil {
				return nil, err
			}
			return &PipeCommand{
				SrcCommand:aCmd,
				DstCommand:bCmd,
				StdOut: os.Stdout,
				StdErr: os.Stderr,
				StdIn: os.Stdin,
				Shell:true,
			}, nil
		}
	}
}
