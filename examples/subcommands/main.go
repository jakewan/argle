package main

import (
	"log"

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

func init() {
	log.SetFlags(0)
}

func main() {
	config := argle.NewConfig()
	config.AddSubcommand(
		"draw",
		argle.WithSubcommand(
			"shapes",
			argle.WithIntArg("count"),
			argle.WithArg("count", 0),
			argle.WithHandler(
				func(a argle.ArgumentHolder) error {
					return drawShapes(DrawShapesArgs{})
				},
			),
		),
		argle.WithSubcommand(
			"lines",
			argle.WithIntArg("count"),
			argle.WithHandler(
				func(a argle.ArgumentHolder) error {
					return drawLines(DrawLinesArgs{})
				},
			),
		),
	)
	log.Printf("Argle config: %v", config)
	if exec, err := config.Parse(); err != nil {
		log.Fatal(err)
	} else {
		exec.Exec()
	}
}

func drawShapes(a DrawShapesArgs) error {
	log.Printf("drawShapes given %v", a)
	return nil
}

func drawLines(a DrawLinesArgs) error {
	log.Printf("drawLines given %v", a)
	return nil
}
