package argle

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type optionProcessor interface {
	normalizedName() string
	processRuntimeOption(runtimeOption) error
}

func normalizedOptionName(name string) string {
	log.Printf("Normalizing %s", name)
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, "-", "")
	log.Printf("Returning %s", result)
	return result
}

type boolOption struct {
	name         string
	value        bool
	defaultValue bool
}

// processRuntimeOption implements optionProcessor.
func (b *boolOption) processRuntimeOption(ro runtimeOption) error {
	log.Printf("processRuntimeOption %+v %+v", b, ro)
	if ro.value == "" {
		b.value = b.defaultValue
	} else if val, err := strconv.ParseBool(ro.value); err != nil {
		return err
	} else {
		b.value = val
	}
	return nil
}

// normalizedName implements optionProcessor.
func (o boolOption) normalizedName() string {
	return normalizedOptionName(o.name)
}

type subcommandOption func(string, *argProc)

func WithBoolOption(name string) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Print("Inside WithBoolOption func")
		config := ap.getSubcommandConfig(subcmd)
		config.options = append(config.options,
			&boolOption{name: name, defaultValue: true},
		)
		ap.subcommandConfigs[subcmd] = config
	}
}

func WithBoolOptionWithDefault(name string, d bool) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Print("Inside WithBoolOptionWithDefault func")
		config := ap.getSubcommandConfig(subcmd)
		config.options = append(config.options,
			&boolOption{name: name, defaultValue: d},
		)
		ap.subcommandConfigs[subcmd] = config
	}
}

func WithGenericHandler[T interface{}](handler SubcommandHandler[T]) subcommandOption {
	return func(subcmd string, ap *argProc) {
		log.Printf("Inside WithGenericHandler %+v", handler)
		config := ap.getSubcommandConfig(subcmd)
		config.handler = func() error {
			return handler.Func(handler.Args)
		}
		ap.subcommandConfigs[subcmd] = config
	}
}

type ArgumentProcessor interface {
	AddSubcommand(subcmd string, opts ...subcommandOption)
	Execute() error
	ExecuteWithArgs(a []string) error
}

type subcommandConfig struct {
	handler func() error
	options []optionProcessor
}

// translateRuntimeArgs takes the runtime arg given by the user and updates
// the subcommand's internal typesafe argument, which will be passed to the
// subcommand handler in a later step.
func (sc subcommandConfig) translateRuntimeArgs(ra runtimeArgs) error {
	log.Printf("translateRuntimeArgs subcommand config: %+v", sc)
	log.Printf("translateRuntimeArgs runtime args: %+v", ra)
	slices.SortFunc(sc.options, func(a, b optionProcessor) int {
		if a.normalizedName() < b.normalizedName() {
			return -1
		}
		if a.normalizedName() > b.normalizedName() {
			return 1
		}
		return 0
	})
	for _, runtimeOpt := range ra.options {
		log.Printf("Runtime arg: %+v", runtimeOpt)
		idx, found := slices.BinarySearchFunc(
			sc.options,
			runtimeOpt,
			func(op optionProcessor, ro runtimeOption) int {
				if op.normalizedName() < ro.normalizedName() {
					return -1
				}
				if op.normalizedName() > ro.normalizedName() {
					return 1
				}
				return 0
			},
		)
		if !found {
			return fmt.Errorf("no configuration found for runtime option %s", runtimeOpt.name)
		}
		configuredOpt := sc.options[idx]
		configuredOpt.processRuntimeOption(runtimeOpt)
	}
	return nil
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

func (o runtimeOption) normalizedName() string {
	return normalizedOptionName(o.name)
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
	if err := subcommandConfig.translateRuntimeArgs(runtimeArgs); err != nil {
		return err
	}
	return subcommandConfig.handler()
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

type SubcommandHandler[T interface{}] struct {
	Args T
	Func func(T) error
}
