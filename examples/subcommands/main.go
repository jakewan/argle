package main

import (
	"log"

	"github.com/jakewan/argle"
)

type Shape int

const (
	shapeCircle Shape = iota
	shapeRectangle
	shapeTriangle
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
			argle.WithStringOptionsArg(
				"shape",
				argle.WithStringOption("circle", shapeCircle),
				argle.WithStringOption("rectangle", shapeRectangle),
				argle.WithStringOption("triangle", shapeTriangle),
			),
			argle.WithHandler(
				func(a argle.ArgumentHolder) error {
					return drawShapes(DrawShapesArgs{})
				},
			),
		),
		argle.WithSubcommand(
			"lines",
			argle.WithIntArg("count"),
			argle.WithFloat32Arg("line-length"),
			argle.WithHandler(
				func(a argle.ArgumentHolder) error {
					c, err := a.GetIntArg("count")
					if err != nil {
						return err
					}
					var l float32 = 0.0
					return drawLines(DrawLinesArgs{
						count:      c,
						lineLength: l,
					})
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
