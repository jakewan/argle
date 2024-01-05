package internal

import "fmt"

type DisplayHelp struct{}

// Exec implements argle.Executor.
func (*DisplayHelp) Exec() error {
	fmt.Print("Displaying help\n")
	return nil
}
