package argle_test

import (
	"log"
	"testing"

	"github.com/jakewan/argle"
	"github.com/stretchr/testify/mock"
)

func init() {
	log.SetFlags(0)
}

type testParser struct {
	mock.Mock
}

// SetBoolOption implements argle.Parser.
func (p *testParser) SetBoolOption(opt string, v bool) error {
	args := p.Called(opt, v)
	return args.Error(0)
}

func TestSubcommandWithBoolOption(t *testing.T) {
	p := argle.NewArgumentProcessor()
	parser := testParser{}
	p.AddSubcommand("subcommand", &parser)
	p.ExecuteWithArgs([]string{"someprogram", "subcommand"})
}
