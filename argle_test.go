package argle_test

import (
	"testing"

	"github.com/jakewan/argle"
)

func TestArgle(t *testing.T) {
	type testConfig struct {
		description string
		config      argle.Config
		args        []string
	}
	for _, cfg := range []testConfig{
		{
			description: "no args; display help",
			config:      argle.NewConfig(),
			args:        []string{"progname"},
		},
	} {
		t.Run(cfg.description, func(t *testing.T) {
			executor, err := cfg.config.ParseWithArgs(cfg.args)
			if err != nil {
				t.Fatal(err)
			}
			err = executor.Exec()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
