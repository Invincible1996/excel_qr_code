package model

import "image/draw"

type DrawText struct {
	JPG   draw.Image
	Title string
	X0    int
	Y0    int
	Size0 float64

	Subtitle string
	X1       int
	Y1       int
	Size1    float64
}
