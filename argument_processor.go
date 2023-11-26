package argle

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type optionProcessor interface {
	processArgs(o runtimeOption) error
}

type boolOption struct {
	name string
}

func (boolOption) processArgs(o runtimeOption) error {
	log.Printf("boolOption.ProcessOption %+v", o)
	return nil
}

type subcommandOption func(string, *argProc)

func WithBoolOption(name string) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Print("Inside WithBoolOption func")
		config := ap.getSubcommandConfig(subcmd)
		config.options = append(config.options, boolOption{name})
		ap.subcommandConfigs[subcmd] = config
	}
}

func WithHandler(handler func(int)) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Print("Inside WithHandler func")
		config := ap.getSubcommandConfig(subcmd)
		config.handler = handler
		ap.subcommandConfigs[subcmd] = config
	}
}

type ArgumentProcessor interface {
	AddSubcommand(subcmd string, opts ...subcommandOption)
	Execute() error
	ExecuteWithArgs(a []string) error
}

type subcommandConfig struct {
	handler func(int)
	options []optionProcessor
}

func (sc subcommandConfig) translateRuntimeArgs(ra runtimeArgs) int {
	return 123
}

type argProc struct {
	subcommandConfigs map[string]subcommandConfig
}

func (processor *argProc) AddSubcommand(subcmd string, opts ...subcommandOption) {
	for _, o := range opts {
		o(subcmd, processor)
	}
}

func (ap *argProc) getSubcommandConfig(subcmd string) subcommandConfig {
	config, ok := ap.subcommandConfigs[subcmd]
	if !ok {
		config = subcommandConfig{
			options: []optionProcessor{},
		}
		ap.subcommandConfigs[subcmd] = config
	}
	return config
}

type runtimeOption struct {
	name  string
	value string
}

type runtimeArgs struct {
	subcommandList []string
	options        []runtimeOption
}

func (s runtimeArgs) subcommandListAsString() string {
	return strings.Join(s.subcommandList, " ")
}

func (ap argProc) ExecuteWithArgs(args []string) error {
	log.Printf("ExecuteWithArgs arg processor: %+v", ap)
	log.Printf("Args: %+v", args)
	runtimeArgs := parseRuntimeArgs(args[1:])
	log.Printf("Runtime args: %+v", runtimeArgs)
	subcommandConfig, ok := ap.subcommandConfigs[runtimeArgs.subcommandListAsString()]
	log.Printf("Subcommand config: %+v", subcommandConfig)
	if !ok {
		return fmt.Errorf("invalid subcommand key %s", runtimeArgs.subcommandListAsString())
	}
	typeSafeArgs := subcommandConfig.translateRuntimeArgs(runtimeArgs)
	log.Printf("Type-safe runtime args: %+v", typeSafeArgs)
	subcommandConfig.handler(typeSafeArgs)
	return nil
}

func (proc argProc) Execute() error {
	return proc.ExecuteWithArgs(os.Args)
}

func NewArgumentProcessor() ArgumentProcessor {
	return &argProc{
		subcommandConfigs: map[string]subcommandConfig{},
	}
}

func parseRuntimeArgs(a []string) runtimeArgs {
	result := runtimeArgs{
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
				runtimeOption{name: optionParts[0], value: optionParts[1]},
			)
		} else {
			nextArgIdx := currentArgIdx + 1
			var value string
			if len(a) > nextArgIdx {
				value = a[nextArgIdx]
			}
			result.options = append(
				result.options,
				runtimeOption{name: optionParts[0], value: value},
			)
		}
		nextArgIdx += 1
	}
	return result
}
