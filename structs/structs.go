package structs

import "math"

type Shape interface {
	Perimeter() float64
	Area() float64
}

type Rectangle struct {
	Length float64
	Width  float64
}

func (rectangle Rectangle) Perimeter() float64 {
	return 2 * (rectangle.Length + rectangle.Width)
}

func (rectangle Rectangle) Area() float64 {
	return rectangle.Length * rectangle.Width
}

type Circle struct {
	Radius float64
}

func (circle Circle) Perimeter() float64 {
	return 2 * math.Pi * circle.Radius
}

func (circle Circle) Area() float64 {
	return math.Pi * circle.Radius * circle.Radius
}

type Triangle struct {
	Base   float64
	Height float64
}

func (triangle Triangle) Perimeter() float64 {
	return triangle.Base + 2*triangle.Height
}

func (triangle Triangle) Area() float64 {
	return triangle.Base * triangle.Height * .5
}
