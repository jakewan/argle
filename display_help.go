package argle

import (
	"fmt"
)

type displayHelp struct {
	config *tempConfig
}

func newDisplayHelp(config *tempConfig) displayHelp {
	return displayHelp{
		config,
	}
}

// Exec implements argle.Executor.
func (d displayHelp) Exec() error {
	fmt.Printf("Displaying help for config %v\n", d.config)
	return nil
}
