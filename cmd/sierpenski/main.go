package main

import (
	"math"
	"os"
	"strconv"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const INITIAL_SIZE float64 = 50

func triangle(length float64) (sdf.SDF2, error) {
	halfL := length / 2
	poly, err := sdf.Polygon2D(sdf.V2Set{
		{0, math.Sqrt(length*length - (halfL)*(halfL))},
		{halfL, 0},
		{-halfL, 0},
	})
	return poly, err
}

const ITER_MAX int = 2

func sierp(triangle sdf.SDF2, i int) sdf.SDF2 {
	if i > ITER_MAX {
		return nil
	}
	height := triangle.BoundingBox().Size().Y
	width := triangle.BoundingBox().Size().X
	scaled := sdf.ScaleUniform2D(triangle, 0.53)
	// Make 3 copies that line up at the vertices.
	copies := []sdf.SDF2{
		sdf.Transform2D(scaled,
			sdf.Translate2d(triangle.BoundingBox().BottomLeft().Sub(scaled.BoundingBox().BottomLeft()))),
		sdf.Transform2D(scaled,
			sdf.Translate2d(triangle.BoundingBox().BottomLeft().Sub(scaled.BoundingBox().BottomLeft()).Add(sdf.V2{X: width / 2, Y: 0}))),
		sdf.Transform2D(scaled,
			sdf.Translate2d(triangle.BoundingBox().BottomLeft().Sub(scaled.BoundingBox().BottomLeft()).Add(sdf.V2{X: width / 4, Y: height / 2}))),
	}
	for j := range copies {
		inner := sierp(copies[j], i+1)
		if inner != nil {
			copies[j] = sierp(copies[j], i+1)
		}
	}
	sip := sdf.Union2D(copies...)
	return sip
}

func sierpenski_ex() (sdf.SDF3, error) {
	tri, err := triangle(INITIAL_SIZE)
	if err != nil {
		return nil, err
	}
	// Scale width and height by 1/2
	sip := sierp(tri, 0)
	return sdf.Extrude3D(sip, 5), nil
}

var NUMCELLS = 0

func main() {
	numCells, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	NUMCELLS = numCells
	box, err := sierpenski_ex()
	if err != nil {
		panic(err)
	}
	render.RenderSTL(box, NUMCELLS, "./sierp.stl")

}
