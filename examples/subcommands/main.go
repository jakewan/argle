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
	argle.NewConfig().
		AddSubcommand(
			"draw",
			argle.WithSubcommand(
				"shapes",
				argle.WithArg[int]("count"),
				argle.WithArg[Shape](
					"shape",
					argle.WithArgOption(shapeCircle),
					argle.WithArgOption(shapeRectangle),
					argle.WithArgOption(shapeTriangle),
				),
				argle.WithHandler(drawShapes),
			),
			argle.WithSubcommand(
				"lines",
				argle.WithArg[int]("count"),
				argle.WithArg[float32]("line-length"),
				argle.WithHandler(drawLines),
			),
		).Run()
}

func drawShapes(a DrawShapesArgs) error {
	log.Printf("drawShapes (shape=%d,count=%d)", a.shape, a.count)
	return nil
}

func drawLines(a DrawLinesArgs) error {
	log.Printf("drawLines (line length=%f,count=%d)", a.lineLength, a.count)
	return nil
}
