package argle

import (
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
		log.Printf("Subcommand config: %+v", config)
	}
}

func WithHandler(handler func()) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Print("Inside WithHandler func")
		config := ap.getSubcommandConfig(subcmd)
		config.handler = handler
		ap.subcommandConfigs[subcmd] = config
		log.Printf("Subcommand config: %+v", config)
	}
}

type ArgumentProcessor interface {
	AddSubcommand(subcmd string, opts ...subcommandOption)
	Execute() error
	ExecuteWithArgs(a []string) error
}

type subcommandConfig struct {
	handler func()
	options []optionProcessor
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

func (proc argProc) ExecuteWithArgs(args []string) error {
	log.Printf("Arg processor: %+v", proc)
	log.Printf("args: %+v", args)
	runtimeArgs := parseRuntimeArgs(args[1:])
	log.Printf("Runtime args: %+v", runtimeArgs)
	// log.Printf("Processor options %+v", proc.options)
	// subcommandOptions := proc.options[runtimeArgs.subcommandListAsString()]
	// if subcommandOptions == nil {
	// 	return fmt.Errorf("invalid subcommand key %s", runtimeArgs.subcommandListAsString())
	// }
	// log.Printf("Subcommand options: %+v", subcommandOptions)
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
			if len(a) > nextArgIdx {
				nextArg := a[nextArgIdx]
				result.options = append(
					result.options,
					runtimeOption{name: optionParts[0], value: nextArg},
				)
			} else {
				log.Printf("No next arg")
			}
		}
		nextArgIdx += 1
	}
	return result
}
