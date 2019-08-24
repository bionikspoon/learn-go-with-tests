package structs

import "testing"

func TestPerimeter(t *testing.T) {
	tests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{"find perimeter of a rectangle", Rectangle{Length: 10.0, Width: 10.0}, 40.0},
		{"find perimeter of a circle", Circle{Radius: 10.0}, 62.83185307179586},
		{"find perimeter of a triangle", Triangle{Base: 12.0, Height: 6.0}, 24.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.shape.Perimeter(); got != tt.want {
				t.Errorf("%#v got %v want %v", tt.shape, got, tt.want)
			}
		})
	}
}

func TestArea(t *testing.T) {

	tests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{"find area of a rectangle", Rectangle{Length: 12.0, Width: 6.0}, 72.0},
		{"find area of a circle", Circle{Radius: 12.0}, 452.3893421169302},
		{"find area of a triangle", Triangle{Base: 12.0, Height: 6.0}, 36.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.shape.Area(); got != tt.want {
				t.Errorf("%#v got %v want %v", tt.shape, got, tt.want)
			}
		})
	}
}
