package wtv

import (
	"image"
	"math"
)

type Shape interface {
	In(x, y int) bool
	Move(x, y int) error
	Point() (float64, float64)
}

type Rectangle struct {
	x int
	y int
	w int
	h int
}

func NewRectangle(x, y, w, h int) *Rectangle {
	var r Rectangle
	r.x = x
	r.y = y
	r.w = w
	r.h = h
	return &r
}

func (r Rectangle) In(x, y int) bool {
	ir := image.Rect(r.x, r.y, r.w+r.x, r.h+r.y)
	return image.Pt(x, y).In(ir)
}

func (r *Rectangle) Move(x, y int) error {
	r.x = x
	r.y = y
	return nil
}

func (r *Rectangle) Point() (float64, float64) {
	return float64(r.x), float64(r.y)
}

type Circle struct {
	r int
	x int
	y int
}

func NewCircle(x, y, r int) *Circle {
	var c Circle
	c.x = x
	c.y = y
	c.r = r
	return &c
}

func (c Circle) In(x, y int) bool {
	dx := x - c.x
	dy := y - c.y
	dr := math.Sqrt(float64(dx*dx + dy*dy))
	return c.r > int(dr)
}

func (c *Circle) Move(x, y int) error {
	c.x = x
	c.y = y
	return nil
}

func (c *Circle) Point() (float64, float64) {
	return float64(c.x - c.r), float64(c.y - c.r)
}

/*
//TODO 実装時は基準点を別に実装した方がいいかも

type Polygon struct {
	points []*image.Point
}

func (p Polygon) NewPolygon(pts ...image.Point) *Polygon {
	var p Polygon
	p.points = append(pts...)
	return &p
}

const FloatTolerance = 0.00001

func (p Polygon) In(x, y int) bool {
	angle := 0.0
	for idx := 0; idx < len(p.points)-1; idx++ {
		angle += s.calcTan(x, y, idx)
	}
	if math.Abs(2.0*math.Pi-math.Abs(angle)) < FloatTolerance {
		return true
	}
	return false
}

func (p *Polygon) calcTan(x, y int, idx int) float64 {

	a := p.points[idx]
	b := p.points[idx+1]

	ax := a.X - x
	ay := a.Y - y

	bx := b.X - x
	by := b.Y - y

	avb := ax*bx + ay*by
	axb := ax*bx - ay*by

	return math.Atan2(float64(axb), float64(avb))
}

*/
