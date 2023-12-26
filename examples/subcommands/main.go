package main

import (
	"fmt"

	"github.com/jakewan/argle"
)

type Shape int

const (
	circle Shape = iota
	rectangle
	triangle
)

type DrawShapesArgs struct {
	count Shape
}

type DrawShapesHandler struct{}

// Exec implements argle.SubcommandHandler.
func (DrawShapesHandler) Exec() {
	fmt.Println("Inside Exec")
}

func main() {
	config := argle.NewConfig()
	fmt.Printf("Argle config: %v\n", config)

	var drawShapesArgs DrawShapesArgs
	fmt.Printf("Draw shapes arg: %v\n", drawShapesArgs)

	handler := DrawShapesHandler{}
	config.AddSubcommand([]string{"draw", "shapes"}, handler)
	config.Parse()
}
