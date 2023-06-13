// worldgen - fractured terrain generator
// Copyright (c) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package generator implements a few map generators.
package generator

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
)

type Map struct {
	height, width int
	diagonal      float64
	rnd           *rand.Rand
	points        []float64
	yx            [][]float64 // points indexed by y, x
}

func New(height, width int, rnd *rand.Rand) *Map {
	m := &Map{
		height:   height,
		width:    width,
		diagonal: math.Sqrt(float64(height*height + width*width)),
		rnd:      rnd,
		points:   make([]float64, height*width, height*width),
		yx:       make([][]float64, height),
	}
	for row := 0; row < height; row++ {
		m.yx[row] = m.points[row*width : (row+1)*width]
	}
	return m
}

func (m *Map) AsImage() *image.RGBA {
	height, width := m.Height(), m.Width()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := uint8(m.yx[y][x] * 49)
			pc := color.RGBA{R: red[val], G: green[val], B: blue[val], A: 255}
			img.Set(x, y, pc)
		}
	}
	return img
}

func (m *Map) AsPNG() ([]byte, error) {
	bb := &bytes.Buffer{}
	err := png.Encode(bb, m.AsImage())
	return bb.Bytes(), err
}

func (m *Map) Diagonal() float64 {
	return m.diagonal
}

func (m *Map) Height() int {
	return m.height
}

func (m *Map) Width() int {
	return m.width
}

func (m *Map) FractureCircle(bump float64) {
	height, width, diagonal := m.Height(), m.Width(), m.Diagonal()

	// generate random radius for the circle
	radius := 0
	for n := m.rnd.Float64(); radius < 1; n = m.rnd.Float64() {
		radius = int(n * n * diagonal / 2)
	}
	//log.Printf("fractureCircle: height %3d width %3d diagonal %6.3f radius %3d\n", height, width, diagonal, radius)

	cx, cy := m.rnd.Intn(width), m.rnd.Intn(height)
	//log.Printf("fractureCircle: cx %3d cy %3d radius %3d\n", cx, cy, radius)

	// limit the x and y values that we look at
	miny, maxy := cy-radius-1, cy+radius+1
	if miny < 0 {
		miny = 0
	}
	if maxy > height {
		maxy = height
	}
	minx, maxx := cx-radius-1, cx+radius+1
	if minx < 0 {
		minx = 0
	}
	if maxx > width {
		maxx = width
	}
	//log.Printf("fractureCircle: cx %3d/%4d/%3d/%3d cy %3d/%4d/%3d/%3d radius %3d\n", cx, width, minx, maxx, cy, height, miny, maxy, radius)

	// bump all points within the radius
	rSquared := radius * radius
	for y := miny; y < maxy; y++ {
		for x := minx; x < maxx; x++ {
			dx, dy := x-cx, y-cy
			isInside := dx*dx+dy*dy < rSquared
			if isInside {
				m.yx[y][x] += bump
			}
		}
	}
}

func FractureSlice(bump float64, world [][]float64) {
	height, width := len(world), len(world[0])

	// create a random line on the world map
	var m, b float64
	for {
		x1, y1 := rand.Intn(width), rand.Intn(height)
		x2, y2 := rand.Intn(width), rand.Intn(height)
		if x1 == x2 && y1 == y2 { // want a line, not a single point
			continue
		} else if y1 == y2 { // can't have vertical lines
			continue
		}
		m = float64(x1-x2) / float64(y1-y2)
		b = rand.Float64() * float64(height)
		break
	}

	// y = (() / ()) x
	//log.Printf("kachunk: line y = m(%f)x + b(%f): bump %f\n", m, b, bump)

	// move all the points below the line up or down
	for x := 0; x < width; x++ {
		mxb := int(m*float64(x) + b)
		for y := 0; y < height; y++ {
			if y > mxb { // point is above the line
				world[y][x] += bump
			}
		}
	}
}

// Normalize the values in the map to the range of 0..1
func (m *Map) Normalize() {
	// determine the range of values
	minValue, maxValue := m.points[0], m.points[0]
	for _, val := range m.points {
		if val < minValue {
			minValue = val
		}
		if maxValue < val {
			maxValue = val
		}
	}
	deltaValue := maxValue - minValue
	if deltaValue+maxValue == 0 {
		maxValue++
	}
	// now normalize
	for n, val := range m.points {
		val = (val - minValue) / (maxValue + deltaValue)
		if val < 0 {
			val = 0
		} else if val > 1 {
			val = 1
		}
		m.points[n] = val
	}
	//log.Println(minValue, maxValue, deltaValue)
}

func (m *Map) RandomFractureCircle(n int) {
	for n > 0 {
		// decide the amount that we're going to raise or lower
		switch m.rnd.Intn(2) {
		case 0:
			m.FractureCircle(1)
		case 1:
			m.FractureCircle(-1)
		}
		n--
	}
}

var (
	red = []uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 34, 68, 102, 119, 136, 153, 170, 187,
		0, 34, 34, 119, 187, 255, 238, 221, 204, 187, 170, 153,
		136, 119, 85, 68,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
	green = []uint8{
		0, 0, 17, 51, 85, 119, 153, 204, 221, 238, 255, 255, 255,
		255, 255, 255, 68, 102, 136, 170, 221, 187, 170, 136,
		136, 102, 85, 85, 68, 51, 51, 34,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
	blue = []uint8{
		0, 68, 102, 136, 170, 187, 221, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 0, 0, 0, 0, 0, 34, 34, 34, 34, 34, 34,
		34, 34, 34, 17, 0,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
)
