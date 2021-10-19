package main

import (
	"math"

	"github.com/alecthomas/kong"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/ngynkvn/go-fractals/cmd/cli"
	u "github.com/ngynkvn/go-fractals/src/util"
)

var CLI = &cli.CLISierp

func main() {
	ctx := kong.Parse(CLI)
	ctx.FatalIfErrorf(ctx.Error)

	box, err := sierpenski_ex()
	if err != nil {
		panic(err)
	}
	render.RenderSTL(box, CLI.MeshCells, CLI.Output)

}

func triangle(length float64) (sdf.SDF2, error) {
	halfL := length / 2
	poly, err := sdf.Polygon2D(u.V2Set{
		u.XY(0, math.Sqrt(length*length-(halfL)*(halfL))),
		u.XY(halfL, 0),
		u.XY(-halfL, 0),
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
			sdf.Translate2d(triangle.BoundingBox().BottomLeft().Sub(scaled.BoundingBox().BottomLeft()).Add(u.XY(width/2, 0)))),
		sdf.Transform2D(scaled,
			sdf.Translate2d(triangle.BoundingBox().BottomLeft().Sub(scaled.BoundingBox().BottomLeft()).Add(u.XY(width/4, height/2)))),
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

// Extrudes the 2D shape out to a desired width.
func sierpenski_ex() (sdf.SDF3, error) {
	tri, err := triangle(CLI.Size)
	if err != nil {
		return nil, err
	}
	// Scale width and height by 1/2
	sip := sierp(tri, 0)
	return sdf.Extrude3D(sip, CLI.Extrusion), nil
}
