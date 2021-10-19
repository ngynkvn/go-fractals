package cli

type CLISierp struct {
	Iterations int     `kong:"default=3,name='depth',help='The number of recursive calls to make when generating.'"`
	MeshCells  int     `kong:"default=300,help='Number of MeshCells to use in rendering.'"`
	Size       float64 `kong:"default=30,help='Size in mm of the requested geometry.'"`
	Extrusion  float64 `kong:"default=5,help='Size in mm of extrusion dimension.'"`
	Output     string  `kong:"default='sierpenski.stl',help='Output path to write rendered geometry to.'"`
}

type CLIMenger struct {
	Iterations int     `kong:"default=3,name:'depth',help='The number of recursive calls to make when generating.'"`
	MeshCells  int     `kong:"default=300,help='Number of MeshCells to use in rendering.'"`
	Size       float64 `kong:"default=30,help='Size in mm of the requested geometry.'"`
	Output     string  `kong:"default='menger.stl',help='Output path to write rendered geometry to.'"`
}

var Sierp CLISierp
var Menger CLIMenger
