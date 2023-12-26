package argle

import (
	"fmt"
	"os"
	"strings"
)

type SubcommandHandler interface {
	Exec()
}

type Config interface {
	AddSubcommand(sc []string, h SubcommandHandler)
	Parse()
	ParseWithArgs(a []string)
}

type tempConfig struct {
	handlers map[string]SubcommandHandler
}

func (c *tempConfig) Parse() {
	c.ParseWithArgs(os.Args)
}

func (c *tempConfig) ParseWithArgs(a []string) {
	fmt.Printf("ParseWithArgs given %s\n", a)
	prog := a[0]
	fmt.Printf("Program: %s\n", prog)
	remaining := a[1:]
	fmt.Printf("Remaining: %s\n", remaining)
	key := strings.Join(remaining, " ")
	h, ok := c.handlers[key]
	if !ok {
		fmt.Println("Handler not found")
		return
	}
	h.Exec()
}

func (c *tempConfig) AddSubcommand(sc []string, h SubcommandHandler) {
	fmt.Printf("Inside AddSubcommand\n")
	key := strings.Join(sc, " ")
	c.handlers[key] = h
}

func NewConfig() Config {
	return &tempConfig{
		handlers: map[string]SubcommandHandler{},
	}
}
