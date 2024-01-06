package argle

import (
	"fmt"
)

type displayHelp struct {
	c *config
}

func newDisplayHelp(c *config) displayHelp {
	return displayHelp{
		c,
	}
}

// Exec implements argle.Executor.
func (d displayHelp) Exec() error {
	fmt.Printf("Displaying help for config %v\n", d.c)
	return nil
}
