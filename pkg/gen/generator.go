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

// Package gen implements a few map generators.
package gen

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
)

type Map struct {
	height, width int
	diagonal      float64
	rnd           *rand.Rand
	points        []int
	yx            [][]int // points indexed by y, x
}

func New(height, width int, rnd *rand.Rand) *Map {
	m := &Map{
		height:   height,
		width:    width,
		diagonal: math.Sqrt(float64(height*height + width*width)),
		rnd:      rnd,
		points:   make([]int, height*width, height*width),
		yx:       make([][]int, height),
	}
	for row := 0; row < height; row++ {
		m.yx[row] = m.points[row*width : (row+1)*width]
	}
	return m
}

func (m *Map) AsPNG(img *image.RGBA) ([]byte, error) {
	bb := &bytes.Buffer{}
	err := png.Encode(bb, img)
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

func (m *Map) FractureCircle(bump int) {
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
	minx, maxx := cx-radius-1, cx+radius+1
	//log.Printf("fractureCircle: cx %3d/%4d/%3d/%3d cy %3d/%4d/%3d/%3d radius %3d\n", cx, width, minx, maxx, cy, height, miny, maxy, radius)

	// bump all points within the radius
	rSquared := radius * radius
	for y := miny; y < maxy; y++ {
		for x := minx; x < maxx; x++ {
			dx, dy := x-cx, y-cy
			isInside := dx*dx+dy*dy < rSquared
			if isInside {
				px, py := x, y
				for px < 0 {
					px += width
				}
				for px >= width {
					px -= width
				}
				for py < 0 {
					py += height
				}
				for py >= height {
					py -= height
				}
				m.yx[py][px] += bump
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
	// fetch the minimum value in the set of points
	minValue, maxValue := m.points[0], m.points[0]
	for _, val := range m.points {
		if val < minValue {
			minValue = val
		}
		if maxValue < val {
			maxValue = val
		}
	}
	log.Printf("normalize: min %8d max %8d\n", minValue, maxValue)
	// update the values in the set so that zero is now the minimum value
	for n, val := range m.points {
		m.points[n] = val - minValue
	}
	minValue = 0
	// fetch the maximum value in the set of points
	maxValue = m.points[0]
	for _, val := range m.points {
		if maxValue < val {
			maxValue = val
		}
	}
	log.Printf("normalize: min %8d max %8d\n", minValue, maxValue)
	deltaValue := maxValue - minValue
	if deltaValue == 0 {
		// all the points are zero so there's no work to be done
		return
	}
	// now normalize to range of 0..255
	for n, val := range m.points {
		m.points[n] = val * 255 / maxValue
	}
	minValue, maxValue = m.points[0], m.points[0]
	for _, val := range m.points {
		if val < minValue {
			minValue = val
		}
		if maxValue < val {
			maxValue = val
		}
	}
	log.Printf("normalize: min %8d max %8d\n", minValue, maxValue)
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

func (m *Map) Shift(dx, dy int) error {
	height, width := m.Height(), m.Width()

	if dx < 0 || dx > width {
		return fmt.Errorf("invalid dx shift")
	} else if dx == width {
		dx = 0
	}
	if dx = dx * -1; dx != 0 {
		for y := 0; y < height; y++ {
			shiftX(m.yx[y], dx)
		}
		log.Printf("shifted m x %d\n", dx)
	}

	if dy < 0 || dy > height {
		return fmt.Errorf("invalid dy shift")
	} else if dy == height {
		dy = 0
	}
	if dy != 0 {
		shiftY(m.yx, dy)
	}

	return nil
}

func shiftX(s []int, n int) {
	for n < 0 {
		n += len(s)
	}
	for n > len(s) {
		n -= len(s)
	}
	if n == 0 {
		return
	}
	tmp := make([]int, n)
	copy(tmp, s[len(s)-n:])
	copy(s[n:], s)
	copy(s, tmp)
}

func shiftY(s [][]int, n int) {
	for n < 0 {
		n += len(s)
	}
	for n > len(s) {
		n -= len(s)
	}
	if n == 0 {
		return
	}
	tmp := make([][]int, n)
	copy(tmp, s[len(s)-n:])
	copy(s[n:], s)
	copy(s, tmp)
}
