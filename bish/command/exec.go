package command

import (
	"errors"
	"fmt"
)

/*
Grammar:
	Command  :- Argument [Argument...]
			 |- Command | Command

	Argument :- String
			 |- ( Command )
			 |- < String >
*/

type ExecutionPlanner struct {
	Args []string

	ctx ShellContext

	idx        int
	currentTok *string
	nextTok    *string
	done       bool
}

func NewExecutionPlanner(ctx ShellContext, args []string) *ExecutionPlanner {
	e := &ExecutionPlanner{Args: args, ctx: ctx}

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
		return &StringArgument{*p.currentTok}, nil
	}
}

// Command reads a command from the current execution plan.
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
				Ctx:  p.ctx,
			}
			if topLevel {
				baseCmd.Shell = true
			}
			return baseCmd, nil
		} else if p.accept("|") {
			aCmd := &CommandTree{
				Args:  argumentBuffer,
				Shell: true,
				Ctx:   p.ctx,
			}
			bCmd, err := p.Command(true)
			if err != nil {
				return nil, err
			}
			return &PipeCommand{
				SrcCommand: aCmd,
				DstCommand: bCmd,
				Shell:      true,
			}, nil
		}
	}
}
