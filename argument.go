package argle

import "fmt"

type argument struct{}

func newArgument(opts ...argumentOption) *argument {
	a := &argument{}
	for _, o := range opts {
		o(a)
	}
	return a
}

type argumentOption func(*argument)

func WithArgOption[T any](o T) argumentOption {
	fmt.Printf("WithArgOption (option=%+v)\n", o)
	return func(a *argument) {
		fmt.Printf("WithArgOption inner after (option=%+v,argument=%+v)\n", o, a)
	}
}

func WithArg[T any](name string, opts ...argumentOption) subcommandOption {
	fmt.Printf("WithArg (name=%s)\n", name)
	return func(s *subcommand) {
		fmt.Printf("WithArg inner (name=%s,subcommand=%+v)\n", name, s)
		_, ok := s.arguments[name]
		if ok {
			panic(fmt.Sprintf("Argument %s exists", name))
		}
		s.arguments[name] = newArgument()
		fmt.Printf("WithArg inner after (name=%s,subcommand=%+v)\n", name, s)
	}
}
