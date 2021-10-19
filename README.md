# Go Fractals!

Some renders I made based off fractal patterns.

<image src="docs/fractals.png"/>

## Usage

```
$ go run ./cmd/menger --help

Usage: menger.exe

Flags:
  -h, --help                   Show context-sensitive help.
      --iterations=3           The number of recursive calls to make when
                               generating.
      --mesh-cells=300         Number of MeshCells to use in rendering.
      --size=30                Size in mm of the requested geometry.
      --output="menger.stl"    Output path to write rendered geometry to.
```
