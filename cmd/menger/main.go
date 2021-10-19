package main

import (
	"github.com/ngynkvn/go-fractals/cmd/cli"
	u "github.com/ngynkvn/go-fractals/src/util"

	"github.com/alecthomas/kong"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

var CLI = &cli.CLIMenger

func main() {
	ctx := kong.Parse(CLI)
	ctx.FatalIfErrorf(ctx.Error)

	box, err := cube(CLI.Size)
	if err != nil {
		panic(err)
	}
	box, err = menger(box, 0)
	if err != nil {
		panic(err)
	}
	render.RenderSTL(box, CLI.MeshCells, CLI.Output)
}

// Create a 3D box primitive with equal X,Y,Z lengths
func cube(length float64) (sdf.SDF3, error) {
	return sdf.Box3D(u.XYZ(length, length, length), 0)
}

// Offsets to "rotate" around a vector and create the sub holes.
var faceOffsets = u.V3Set{
	u.XYZ(-1, 0, 0),
	u.XYZ(1, 0, 0),
	u.XYZ(0, -1, 0),
	u.XYZ(0, 1, 0),
	u.XYZ(0, 0, -1),
	u.XYZ(0, 0, 1),
}

// Offsets for recursing the menger function.
var z_offs = []float64{-1, 0, 1}
var xy_offs = u.V2Set{u.XY(1, 0), u.XY(1, 1), u.XY(0, 1), u.XY(-1, 1), u.XY(-1, 0), u.XY(-1, -1), u.XY(0, -1), u.XY(1, -1)}

const ABOUT_ONE float64 = 1.0 + 0.01

func createOffsets() u.V3Set {
	set := make(u.V3Set, len(z_offs)*len(xy_offs))
	for _, z := range z_offs {
		for _, xy := range xy_offs {
			off := u.XYZ(xy.X, xy.Y, z)
			// Ignore vectors that have a magnitude of one,
			// they map to the offset for the face holes which does not have a cube to recurse onto.
			if off.Length() < ABOUT_ONE {
				continue
			}
			set = append(set, off)
		}
	}
	return set
}

var offsets = createOffsets()

func menger(box sdf.SDF3, i int) (sdf.SDF3, error) {
	if i >= CLI.Iterations {
		return nil, nil
	}
	size := box.BoundingBox().Size().X / 3
	centerHole, err := cube(size)
	if err != nil {
		return nil, err
	}
	// Punch the center hole.
	cut := sdf.Difference3D(box, sdf.Transform3D(
		centerHole,
		sdf.Translate3d(box.BoundingBox().Center()),
	))
	// Punch the face holes.
	for _, offset := range faceOffsets {
		offset = offset.MulScalar(size).Add(box.BoundingBox().Center())
		translation := sdf.Translate3d(offset)
		faceHole := sdf.Transform3D(centerHole, translation)
		cut = sdf.Difference3D(cut, faceHole)
	}
	// Along the z-axis, recurse the menger function.
	for _, next_off := range offsets {
		next_off = next_off.MulScalar(size).Add(box.BoundingBox().Center())
		next := sdf.Transform3D(
			centerHole,
			sdf.Translate3d(next_off),
		)
		innerMenger, err := menger(next, i+1)
		if err != nil {
			return nil, err
		} else if innerMenger != nil {
			cut = sdf.Difference3D(cut,
				sdf.Difference3D(
					next,
					innerMenger,
				))
		}
	}
	return cut, nil
}
