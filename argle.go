package argle

import (
	"errors"
	"fmt"
	"os"
)

type runtimeArgs struct{}

// Load implements RuntimeArguments.
func (runtimeArgs) Load(target any) error {
	panic("unimplemented")
}

type SubcommandHandler = func(a runtimeArgs) error

type subcommand struct {
	arguments   map[string]*argument
	handler     SubcommandHandler
	subcommands map[string]*subcommand
}

func newSubcommand(opts ...subcommandOption) *subcommand {
	newSub := &subcommand{
		arguments:   map[string]*argument{},
		subcommands: map[string]*subcommand{},
	}
	for _, o := range opts {
		o(newSub)
	}
	return newSub
}

func (sc *subcommand) findSubcommand(tokens []string) (Executor, error) {
	fmt.Printf("findSubcommand tokens: %s\n", tokens)
	return nil, errors.New("not implemented")
}

type Executor interface {
	Exec() error
}

// type tempExecutor struct {
// 	args    []string
// 	handler SubcommandHandler
// }

// func (ex *tempExecutor) Exec() error {
// 	fmt.Printf("Exec args: %s\n", ex.args)
// 	return ex.handler(runtimeArgs{})
// }

type subcommandOption func(*subcommand)

func WithHandler[T any](f func(a T) error) subcommandOption {
	h := func(a runtimeArgs) error {
		var args T
		err := a.Load(&args)
		if err != nil {
			return err
		}
		return f(args)
	}
	return func(s *subcommand) {
		s.handler = h
	}
}

type argument struct{}

func newArgument(opts ...argumentOption) *argument {
	a := &argument{}
	for _, o := range opts {
		o(a)
	}
	return a
}

type argumentOption func(*argument)

func WithArgOption[T any](o T) argumentOption {
	fmt.Printf("WithArgOption (option=%+v)\n", o)
	return func(a *argument) {
		fmt.Printf("WithArgOption inner after (option=%+v,argument=%+v)\n", o, a)
	}
}

func WithArg[T any](name string, opts ...argumentOption) subcommandOption {
	fmt.Printf("WithArg (name=%s)\n", name)
	return func(s *subcommand) {
		fmt.Printf("WithArg inner (name=%s,subcommand=%+v)\n", name, s)
		_, ok := s.arguments[name]
		if ok {
			panic(fmt.Sprintf("Argument %s exists", name))
		}
		s.arguments[name] = newArgument()
		fmt.Printf("WithArg inner after (name=%s,subcommand=%+v)\n", name, s)
	}
}

func WithSubcommand(name string, opts ...subcommandOption) subcommandOption {
	fmt.Printf("WithSubcommand (name=%s)\n", name)
	return func(s *subcommand) {
		_, ok := s.subcommands[name]
		if ok {
			panic(fmt.Sprintf("Subcommand %s exists", name))
		}
		s.subcommands[name] = newSubcommand(opts...)
		fmt.Printf("WithSubcommand inner (name=%s,subcommand=%+v)\n", name, s)
	}
}

type Config interface {
	AddSubcommand(name string, opts ...subcommandOption) Config
	Parse() (Executor, error)
	ParseWithArgs(a []string) (Executor, error)
	Run()
}

type invalidSubcommandBehavior int

const (
	invalidSubcommandBehaviorDisplayHelp invalidSubcommandBehavior = iota
)

type tempConfig struct {
	invalidSubcommandBehavior invalidSubcommandBehavior
	subcommands               map[string]*subcommand
}

func (c *tempConfig) Run() {
	fmt.Printf("Argle config: %v\n", c)
	if exec, err := c.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		exec.Exec()
	}
}

func (c *tempConfig) AddSubcommand(name string, opts ...subcommandOption) Config {
	_, ok := c.subcommands[name]
	if ok {
		panic(fmt.Sprintf("Subcommand %s exists", name))
	}
	s := newSubcommand(opts...)
	fmt.Printf("%+v\n", *s)
	c.subcommands[name] = s
	return c
}

func (c *tempConfig) Parse() (Executor, error) {
	return c.ParseWithArgs(os.Args)
}

func (c *tempConfig) ParseWithArgs(a []string) (Executor, error) {
	fmt.Printf("ParseWithArgs given %s\n", a)
	prog := a[0]
	fmt.Printf("Program: %s\n", prog)
	tokens := a[1:]
	fmt.Printf("Tokens: %s\n", tokens)
	if len(tokens) == 0 {
		if c.invalidSubcommandBehavior == invalidSubcommandBehaviorDisplayHelp {
			return newDisplayHelp(c), nil
		}
		return nil, noSubcommandGiven{}
	}
	token := tokens[0]
	fmt.Printf("Current token: %s\n", token)
	sc, ok := c.subcommands[token]
	if !ok {
		return nil, fmt.Errorf("subcommand not found: %s", token)
	}
	fmt.Printf("Current subcommand: %v\n", sc)
	return sc.findSubcommand(tokens[1:])
}

func NewConfig() Config {
	return &tempConfig{
		invalidSubcommandBehavior: invalidSubcommandBehaviorDisplayHelp,
		subcommands:               map[string]*subcommand{},
	}
}
