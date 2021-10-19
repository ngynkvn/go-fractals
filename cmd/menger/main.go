package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const INITIAL_SIZE float64 = 30

var ORIGIN = sdf.V3{
	X: INITIAL_SIZE,
	Y: INITIAL_SIZE,
	Z: INITIAL_SIZE,
}

func cube(length float64) (sdf.SDF3, error) {
	return sdf.Box3D(sdf.V3{length, length, length}, 0)
}

const ITER_MAX int = 2

var faceOffsets = sdf.V3Set{
	sdf.V3{-1, 0, 0},
	sdf.V3{1, 0, 0},
	sdf.V3{0, -1, 0},
	sdf.V3{0, 1, 0},
	sdf.V3{0, 0, -1},
	sdf.V3{0, 0, 1},
}

var z_offs = []float64{-1, 0, 1}
var xy_offs = []sdf.V2{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

func menger(box sdf.SDF3, i int) (sdf.SDF3, error) {
	if i >= ITER_MAX {
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
	for _, z := range z_offs {
		for _, xy := range xy_offs {
			next_off := sdf.V3{xy.X, xy.Y, z}
			if next_off.Length() < 1.0+0.01 {
				continue
			}
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
	}
	return cut, nil
}

func main() {
	box, err := cube(INITIAL_SIZE)
	if err != nil {
		panic(err)
	}
	box, err = menger(box, 0)
	if err != nil {
		panic(err)
	}
	render.RenderSTL(box, 300, "./tiny_menger.stl")

}
