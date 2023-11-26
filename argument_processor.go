package argle

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type optionProcessor interface {
	processArgs(o commandOption) error
}

type stringOption struct{}

func (stringOption) processArgs(o commandOption) error {
	log.Printf("stringOption.ProcessOption")
	return nil
}

type fileOption struct{}

func (fileOption) processArgs(o commandOption) error {
	log.Printf("fileOption.ProcessOption %+v", o)
	return nil
}

type boolOption struct{}

func (boolOption) processArgs(o commandOption) error {
	log.Printf("boolOption.ProcessOption %+v", o)
	return nil
}

type subcommandOption func(string, *argProc)

func WithBoolOption(name string) subcommandOption {
	return func(subcmd string, processor *argProc) {
		processor.options[subcmd] = boolOption{}
	}
}

func WithFileOption(name string) subcommandOption {
	return func(subcmd string, processor *argProc) {
		processor.options[subcmd] = fileOption{}
	}
}

func WithStringOption(name string) subcommandOption {
	return func(subcmd string, processor *argProc) {
		processor.options[subcmd] = stringOption{}
	}
}

type ArgumentProcessor interface {
	AddSubcommand(subcmd string, opts ...subcommandOption)
	Execute() error
	ExecuteWithArgs(a []string) error
}

type argProc struct {
	options map[string]optionProcessor
}

func (processor *argProc) AddSubcommand(subcmd string, opts ...subcommandOption) {
	for _, o := range opts {
		o(subcmd, processor)
	}
}

type commandOption struct {
	name  string
	value string
}

type commandSpec struct {
	subcommandList []string
	options        []commandOption
}

func (s commandSpec) subcommandListAsString() string {
	return strings.Join(s.subcommandList, " ")
}

func (proc argProc) ExecuteWithArgs(a []string) error {
	log.Printf("args: %+v", a)
	userArgs := parseUserArgs(a[1:])
	log.Printf("User args: %+v", userArgs)
	log.Printf("Processor options %+v", proc.options)
	subcommandOptions := proc.options[userArgs.subcommandListAsString()]
	if subcommandOptions == nil {
		return fmt.Errorf("invalid subcommand key %s", userArgs.subcommandListAsString())
	}
	log.Printf("Subcommand options: %+v", subcommandOptions)
	return nil
}

func (proc argProc) Execute() error {
	return proc.ExecuteWithArgs(os.Args)
}

func NewArgumentProcessor() ArgumentProcessor {
	return &argProc{
		options: map[string]optionProcessor{},
	}
}

func parseUserArgs(a []string) commandSpec {
	result := commandSpec{
		subcommandList: []string{},
	}
	nextArgIdx := 0
	for _, currentArg := range a {
		if strings.HasPrefix(currentArg, "-") {
			break
		}
		result.subcommandList = append(result.subcommandList, currentArg)
		nextArgIdx += 1
	}
	a = a[nextArgIdx:]
	log.Printf("After processing subcommand args: %+v", a)
	for currentArgIdx, currentArg := range a {
		if !strings.HasPrefix(currentArg, "-") {
			break
		}
		log.Printf("Processing flag: %s", currentArg)
		option := strings.TrimLeft(currentArg, "-")
		log.Printf("Bare option: %s", option)

		// Handle -opt=value form
		optionParts := strings.SplitN(option, "=", 2)
		log.Printf("Option parts: %s", optionParts)
		optionName := optionParts[0]
		log.Printf("Option name: %s", optionName)
		if len(optionParts) == 2 {
			result.options = append(
				result.options,
				commandOption{name: optionParts[0], value: optionParts[1]},
			)
		} else {
			nextArgIdx := currentArgIdx + 1
			if len(a) > nextArgIdx {
				nextArg := a[nextArgIdx]
				result.options = append(
					result.options,
					commandOption{name: optionParts[0], value: nextArg},
				)
			} else {
				log.Printf("No next arg")
			}
		}
		nextArgIdx += 1
	}
	return result
}
