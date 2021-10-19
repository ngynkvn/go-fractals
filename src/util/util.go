package util

import "github.com/deadsy/sdfx/sdf"

func XY(x, y float64) V2 {
	return sdf.V2{X: x, Y: y}
}

func XYZ(x, y, z float64) V3 {
	return sdf.V3{X: x, Y: y, Z: z}
}

type V2 = sdf.V2
type V2Set = sdf.V2Set
type V3 = sdf.V3
type V3Set = sdf.V3Set
