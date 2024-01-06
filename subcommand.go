package argle

import (
	"errors"
	"fmt"
)

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
