package geometry

import (
	"math"

	"fmt"
	"image/color"
)

type Point struct {
	X, Y float64
}

type Path []Point //named slice type

type ColoredPoint struct {
	Point
	Color color.RGBA
}

func (p *Point) PointDistance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor

}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}
	for i := range path {
		path[i].Add(offset)
		path[i] = op(path[i], offset)
	}
}

func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].PointDistance(path[i])
		}
	}
	return sum
}

func GeometryDemo(x1 float64, y1 float64, x2 float64, y2 float64) {
	p := Point{x1, y1}
	q := Point{x2, y2}

	fmt.Printf("Distance from x1,y1 to x2,y2 %g\n", p.PointDistance(q))
	l := Point{1, 1}
	s := Point{2, 2}
	perim := Path{p, q, l, s}

	fmt.Printf("Permiter for p,q and 1,1 and 2,2 points %g\n", perim.Distance())

	(&p).ScaleBy(6)
	fmt.Printf("Scale p point  %g\n", p)

	blue := color.RGBA{255, 0, 0, 255}
	cp := ColoredPoint{p, blue}
	fmt.Printf("Colored point  %g\n", cp.Point.X)

	p1 := p.Add(q)
	fmt.Printf("Add p to q  %g\n", p1)

	p2 := p.Sub(q)
	fmt.Printf("Substract p from q  %g\n", p2)

	perim.TranslateBy(q, true)
	fmt.Printf("Offset perim to q  %g\n", perim)

}
