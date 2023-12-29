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

type tempSubcommand struct {
	name    string
	handler SubcommandHandler
}

func (sc *tempSubcommand) findSubcommand(tokens []string) (Executor, error) {
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

type subcommandOption func(*tempSubcommand)

func WithHandler[T any](f func(a T) error) subcommandOption {
	h := func(a RuntimeArguments) error {
		var args T
		err := a.Load(&args)
		if err != nil {
			return err
		}
		return f(args)
	}
	return func(s *tempSubcommand) {
		s.handler = h
	}
}

func WithIntArg(name string) subcommandOption {
	return func(s *tempSubcommand) {
		s.name = name
	}
}

func WithFloat32Arg(name string) subcommandOption {
	return func(s *tempSubcommand) {
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
	return func(s *tempSubcommand) {
		s.name = name
	}
}

func WithSubcommand(name string, opts ...subcommandOption) subcommandOption {
	return func(s *tempSubcommand) {

	}
}

type Config interface {
	AddSubcommand(name string, opts ...subcommandOption)
	Parse() (Executor, error)
	ParseWithArgs(a []string) (Executor, error)
}

type tempConfig struct {
	subcommands map[string]*tempSubcommand
}

func (c *tempConfig) AddSubcommand(name string, opts ...subcommandOption) {
	_, ok := c.subcommands[name]
	if ok {
		panic(fmt.Sprintf("Subcommand %s exists", name))
	}
	newSubcommand := &tempSubcommand{}
	for _, o := range opts {
		o(newSubcommand)
	}
	c.subcommands[name] = newSubcommand
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
		subcommands: map[string]*tempSubcommand{},
	}
}
