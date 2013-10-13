package gist6545684

import (
	"fmt"
	. "gist.github.com/5286084.git"
	"github.com/Jragonmiris/mathgl"
	"os"

	"github.com/shurcooL/go-goon"
)

type Contour struct {
	Vertices []mathgl.Vec2d
}

type Polygon struct {
	Contours []Contour
}

func ReadGpcFile(path string) Polygon {
	p := Polygon{}

	file, err := os.Open(path)
	CheckError(err)
	defer file.Close()

	var numContours uint64
	fmt.Fscan(file, &numContours)
	p.Contours = make([]Contour, numContours)

	for contourIndex := range p.Contours {
		var numVertices uint64
		fmt.Fscan(file, &numVertices)
		p.Contours[contourIndex].Vertices = make([]mathgl.Vec2d, numVertices)

		for vertexIndex := range p.Contours[contourIndex].Vertices {
			fmt.Fscan(file, &p.Contours[contourIndex].Vertices[vertexIndex][0], &p.Contours[contourIndex].Vertices[vertexIndex][1])
		}
	}

	return p
}

func main() {
	p := ReadGpcFile("/Users/Dmitri/Dmitri/^Work/^GitHub/eX0/eX0/levels/test_orientation.wwl")
	goon.Dump(p)
}
