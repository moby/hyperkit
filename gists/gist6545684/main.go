// Package gist6545684 parses GPC format files.
package gist6545684

import (
	"fmt"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl64"
)

type Contour struct {
	Vertices []mgl64.Vec2
}

type Polygon struct {
	Contours []Contour
}

func ReadGpcFromReader(r io.Reader) (Polygon, error) {
	p := Polygon{}

	var numContours uint64
	fmt.Fscan(r, &numContours)
	p.Contours = make([]Contour, numContours)

	for contourIndex := range p.Contours {
		var numVertices uint64
		fmt.Fscan(r, &numVertices)
		p.Contours[contourIndex].Vertices = make([]mgl64.Vec2, numVertices)

		for vertexIndex := range p.Contours[contourIndex].Vertices {
			fmt.Fscan(r, &p.Contours[contourIndex].Vertices[vertexIndex][0], &p.Contours[contourIndex].Vertices[vertexIndex][1])
		}
	}

	return p, nil
}

func ReadGpcFile(path string) (Polygon, error) {
	f, err := os.Open(path)
	if err != nil {
		return Polygon{}, err
	}
	defer f.Close()

	return ReadGpcFromReader(f)
}
