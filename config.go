package argle

import (
	"fmt"
	"os"
)

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

type config struct {
	invalidSubcommandBehavior invalidSubcommandBehavior
	subcommands               map[string]*subcommand
}

func (c *config) Run() {
	fmt.Printf("Argle config: %v\n", c)
	if exec, err := c.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		exec.Exec()
	}
}

func (c *config) AddSubcommand(name string, opts ...subcommandOption) Config {
	_, ok := c.subcommands[name]
	if ok {
		panic(fmt.Sprintf("Subcommand %s exists", name))
	}
	s := newSubcommand(opts...)
	fmt.Printf("%+v\n", *s)
	c.subcommands[name] = s
	return c
}

func (c *config) Parse() (Executor, error) {
	return c.ParseWithArgs(os.Args)
}

func (c *config) ParseWithArgs(a []string) (Executor, error) {
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
	return &config{
		invalidSubcommandBehavior: invalidSubcommandBehaviorDisplayHelp,
		subcommands:               map[string]*subcommand{},
	}
}
