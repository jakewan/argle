package argle

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type ArgumentProcessor interface {
	AddSubcommand(subcmd string, parser Parser)
	Execute() error
	ExecuteWithArgs(a []string) error
}

type argProc struct {
	router map[string]Parser
}

type commandSpec struct {
	subcommandList []string
}

func (c commandSpec) subcommandKey() string {
	return strings.Join(c.subcommandList, " ")
}

// ExecuteWithArgs implements ArgumentProcessor.
func (proc argProc) ExecuteWithArgs(a []string) error {
	log.Printf("args: %+v", a)
	log.Printf("router: %+v", proc.router)

	commandSpec := commandSpecFromArgs(a[1:])
	parser, ok := proc.router[commandSpec.subcommandKey()]
	if !ok {
		return fmt.Errorf("subcommand key %s not found", commandSpec.subcommandKey())
	}
	log.Printf("Subcommand parser: %+v", parser)
	return nil
}

// Execute implements ArgumentProcessor.
func (proc argProc) Execute() error {
	return proc.ExecuteWithArgs(os.Args)
}

// AddSubcommand implements ArgumentProcessor.
func (proc argProc) AddSubcommand(subcmd string, parser Parser) {
	proc.router[subcmd] = parser
}

func NewArgumentProcessor() ArgumentProcessor {
	return argProc{
		router: map[string]Parser{},
	}
}

type Parser interface {
	SetBoolOption(opt string, v bool) error
}

func commandSpecFromArgs(a []string) commandSpec {
	result := commandSpec{
		subcommandList: []string{},
	}
	for _, currentArg := range a {
		if strings.HasPrefix(currentArg, "-") {
			break
		}
		result.subcommandList = append(result.subcommandList, currentArg)
	}
	return result
}
