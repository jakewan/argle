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
			argle.WithHandler(drawShapes),
		),
		argle.WithSubcommand(
			"lines",
			argle.WithIntArg("count"),
			argle.WithFloat32Arg("line-length"),
			argle.WithHandler(drawLines),
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
	log.Printf("drawShapes (shape=%d,count=%d)", a.shape, a.count)
	return nil
}

func drawLines(a DrawLinesArgs) error {
	log.Printf("drawLines (line length=%f,count=%d)", a.lineLength, a.count)
	return nil
}
