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
	shape Shape
	count int
}

type DrawLinesArgs struct {
	lineLength float32
	count      int
}

func main() {
	config := argle.NewConfig()
	config.AddSubcommand(
		[]string{"draw", "shapes"},
		argle.WithHandler(
			func(a argle.ArgumentHolder) error {
				return drawShapes(DrawShapesArgs{})
			},
		),
		argle.WithIntArg("count"),
	)
	config.AddSubcommand(
		[]string{"draw", "lines"},
		argle.WithHandler(
			func(a argle.ArgumentHolder) error {
				return drawLines(DrawLinesArgs{})
			},
		),
		argle.WithIntArg("count"),
	)
	fmt.Printf("Argle config: %v\n", config)
	if exec, err := config.Parse(); err != nil {
		fmt.Println(err)
	} else {
		exec.Exec()
	}
}

func drawShapes(a DrawShapesArgs) error {
	fmt.Printf("drawShapes given %v\n", a)
	return nil
}

func drawLines(a DrawLinesArgs) error {
	fmt.Printf("drawLines given %v\n", a)
	return nil
}
