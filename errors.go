package argle

type noSubcommandGiven struct{}

func (noSubcommandGiven) Error() string {
	return "no subcommand given"
}
