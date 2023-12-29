package argle

import (
	"errors"
	"fmt"
	"os"
)

type RuntimeArguments interface {
	Load(target any) error
}

type runtimeArgs struct{}

// Load implements RuntimeArguments.
func (runtimeArgs) Load(target any) error {
	panic("unimplemented")
}

type SubcommandHandler = func(a RuntimeArguments) error

type subcommand struct {
	name        string
	handler     SubcommandHandler
	subcommands map[string]*subcommand
}

func newSubcommand(name string, opts ...subcommandOption) *subcommand {
	newSub := &subcommand{
		name:        name,
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

type tempExecutor struct {
	args    []string
	handler SubcommandHandler
}

func (ex *tempExecutor) Exec() error {
	fmt.Printf("Exec args: %s\n", ex.args)
	return ex.handler(runtimeArgs{})
}

type subcommandOption func(*subcommand)

func WithHandler[T any](f func(a T) error) subcommandOption {
	h := func(a RuntimeArguments) error {
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

func WithIntArg(name string) subcommandOption {
	return func(s *subcommand) {
		s.name = name
	}
}

func WithFloat32Arg(name string) subcommandOption {
	return func(s *subcommand) {
		s.name = name
	}
}

type stringOptionOption[T any] func() T

func WithStringOption[T any](name string, value T) stringOptionOption[T] {
	return func() T {
		fmt.Printf("String option handler %s\n", name)
		return value
	}
}

func WithStringOptionsArg[T any](name string, opts ...stringOptionOption[T]) subcommandOption {
	return func(s *subcommand) {
		s.name = name
	}
}

func WithSubcommand(name string, opts ...subcommandOption) subcommandOption {
	fmt.Printf("WithSubcommand (name=%s)\n", name)
	return func(s *subcommand) {
		_, ok := s.subcommands[name]
		if ok {
			panic(fmt.Sprintf("Subcommand %s exists", name))
		}
		s.subcommands[name] = newSubcommand(name, opts...)
		fmt.Printf("WithSubcommand inner (name=%s,subcommand=%+v)\n", name, s)
	}
}

type Config interface {
	AddSubcommand(name string, opts ...subcommandOption) Config
	Parse() (Executor, error)
	ParseWithArgs(a []string) (Executor, error)
	Run()
}

type tempConfig struct {
	subcommands map[string]*subcommand
}

// Run implements Config.
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
	s := newSubcommand(name, opts...)
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
	token := tokens[0]
	fmt.Printf("Current token: %s\n", token)
	sc, ok := c.subcommands[token]
	if !ok {
		return nil, fmt.Errorf("subcommand not found: %s", token)
	}
	fmt.Printf("Current subcommand: %v\n", sc)
	return sc.findSubcommand(tokens[1:])
	// for k, v := range c.subcommands {
	// 	// scParts := strings.Split(k, " ")
	// 	// if len(incomingArgs) >= len(scParts) {
	// 	// 	fmt.Printf("Comparing %v to %v\n", scParts, incomingArgs[:len(scParts)])
	// 	// 	if slices.Equal(scParts, incomingArgs[:len(scParts)]) {
	// 	// 		fmt.Printf("Found subcommand for %s: %v\n", k, v)
	// 	// 		return &tempExecutor{
	// 	// 			args:    incomingArgs[len(scParts):],
	// 	// 			handler: v.handler,
	// 	// 		}, nil
	// 	// 	}
	// 	// }
	// }
	// return nil, fmt.Errorf("subcommand not found: %s", strings.Join(tokens, " "))
}

func NewConfig() Config {
	return &tempConfig{
		subcommands: map[string]*subcommand{},
	}
}
