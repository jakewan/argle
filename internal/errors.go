package internal

type NoSubcommandGiven struct{}

func (NoSubcommandGiven) Error() string {
	return "no subcommand given"
}
