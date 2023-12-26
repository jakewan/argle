package argle

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type ArgumentHolder interface{}

type tempArgumentHolder struct{}

type SubcommandHandler = func(a ArgumentHolder) error

type Subcommand interface {
	SetHandler(SubcommandHandler) Subcommand
}

type tempSubcommand struct {
	handler SubcommandHandler
}

func (sc *tempSubcommand) SetHandler(handler SubcommandHandler) Subcommand {
	sc.handler = handler
	return sc
}

type Executor interface {
	Exec() error
}

type tempExecutor struct {
	args    []string
	handler SubcommandHandler
}

// Exec implements Executor.
func (ex *tempExecutor) Exec() error {
	fmt.Printf("Exec args: %s\n", ex.args)
	return ex.handler(tempArgumentHolder{})
}

type Config interface {
	AddSubcommand(sc []string) Subcommand
	Parse() (Executor, error)
	ParseWithArgs(a []string) (Executor, error)
}

type tempConfig struct {
	subcommands map[string]*tempSubcommand
}

func (c *tempConfig) AddSubcommand(sc []string) Subcommand {
	key := strings.Join(sc, " ")
	_, ok := c.subcommands[key]
	if ok {
		panic(fmt.Sprintf("Subcommand for %s exists", sc))
	}
	result := &tempSubcommand{}
	c.subcommands[key] = result
	return result
}

func (c *tempConfig) Parse() (Executor, error) {
	return c.ParseWithArgs(os.Args)
}

func (c *tempConfig) ParseWithArgs(a []string) (Executor, error) {
	fmt.Printf("ParseWithArgs given %s\n", a)
	prog := a[0]
	fmt.Printf("Program: %s\n", prog)
	incomingArgs := a[1:]
	fmt.Printf("Incoming args: %s\n", incomingArgs)
	for k, v := range c.subcommands {
		scParts := strings.Split(k, " ")
		if len(incomingArgs) >= len(scParts) {
			fmt.Printf("Comparing %v to %v\n", scParts, incomingArgs[:len(scParts)])
			if slices.Equal(scParts, incomingArgs[:len(scParts)]) {
				fmt.Printf("Found subcommand for %s: %v\n", k, v)
				return &tempExecutor{
					args:    incomingArgs[len(scParts):],
					handler: v.handler,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("subcommand not found: %s", strings.Join(incomingArgs, " "))
}

func NewConfig() Config {
	return &tempConfig{
		subcommands: map[string]*tempSubcommand{},
	}
}
