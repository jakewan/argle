package argle

type runtimeArgs struct{}

// Load implements RuntimeArguments.
func (runtimeArgs) Load(target any) error {
	panic("unimplemented")
}
