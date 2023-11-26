package argle_test

import (
	"log"
	"testing"

	"github.com/jakewan/argle"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(0)
}

type processorSetupFunc = func(argle.ArgumentProcessor)

type testConfig struct {
	description         string
	executeArgs         []string
	processorSetupFunc  processorSetupFunc
	expectedErrorString string
}

func Test(t *testing.T) {
	testTable := []testConfig{
		{
			description: "bool flag with default",
			processorSetupFunc: func(ap argle.ArgumentProcessor) {
				ap.AddSubcommand("some-subcommand", argle.WithBoolOption("some-bool"), argle.WithHandler(func() {
					log.Print("Inside subcommand handler")
				}))
			},
			executeArgs: []string{"some-program", "some-subcommand", "-some-bool"},
		},
	}
	for _, cfg := range testTable {
		t.Run(cfg.description, func(t *testing.T) {
			p := argle.NewArgumentProcessor()
			cfg.processorSetupFunc(p)

			// Code under test
			err := p.ExecuteWithArgs(cfg.executeArgs)

			// Verify
			if err != nil {
				if cfg.expectedErrorString != "" {
					assert.EqualError(t, err, cfg.expectedErrorString)
				} else {
					t.Error(err)
				}
			} else if cfg.expectedErrorString != "" {
				t.Error("ExecuteWithArgs was expected to throw an error")
			}
		})
	}
}
