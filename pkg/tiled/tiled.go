/*
 * worldgen - fractured terrain generator
 * Copyright (C) 1999  John Olsson
 * Copyright (C) 2023 Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package tiled

import (
	"log"
	"math"
	"math/rand"
	"time"
)

func Run(height, width, iterations int, saveFile string) error {
	started := time.Now()

	world := twoDimensionalArray(height, width)
	for iterations > 0 {
		// decide the amount that we're going to raise or lower
		switch rand.Intn(2) {
		case 0:
			fracture(rand.Intn(2) == 0, 1, world)
		case 1:
			fracture(rand.Intn(2) == 0, -1, world)
		}
		iterations--
	}

	normalizeMap(world)

	img := generateImage(world)

	if err := savePNG(saveFile, img); err != nil {
		return err
	}
	log.Printf("tile: created %s: %v\n", saveFile, time.Now().Sub(started))

	return nil
}

func fracture(inside bool, bump float64, world [][]float64) {
	height, width := len(world), len(world[0])
	diagonal := math.Sqrt(float64(height*height + width*width))
	radius := 0
	for n := rand.Float64(); radius < 1; n = rand.Float64() {
		radius = int(n * n * diagonal / 2)
	}
	//log.Printf("fracture: height %3d width %3d diagonal %6.3f radius %3d\n", height, width, diagonal, radius)

	cx, cy := rand.Intn(width), rand.Intn(height)
	//log.Printf("fracture: cx %3d cy %3d radius %3d\n", cx, cy, radius)

	// bump all points within in the radius
	rSquared := radius * radius
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx, dy := x-cx, y-cy
			isInside := dx*dx+dy*dy < rSquared
			if !isInside {
				world[y][x] += bump
			}
		}
	}
}

// normalizeMap normalizes the values in the map to the range of 0..1
func normalizeMap(world [][]float64) {
	height, width := len(world), len(world[0])
	// determine the range of values
	minValue, maxValue := world[0][0], world[0][0]
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if world[y][x] < minValue {
				minValue = world[y][x]
			}
			if maxValue < world[y][x] {
				maxValue = world[y][x]
			}
		}
	}
	deltaValue := maxValue - minValue
	if deltaValue+maxValue == 0 {
		maxValue++
	}
	// now normalize
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			world[y][x] = (world[y][x] + minValue) / (maxValue + deltaValue)

			if world[y][x] < 0 {
				world[y][x] = 0
			} else if world[y][x] < 0 || world[y][x] > 1 {
				world[y][x] = 1
			}
		}
	}
	log.Println(minValue, maxValue, deltaValue)
}

func twoDimensionalArray(height, width int) [][]float64 {
	v, rows := make([][]float64, height), make([]float64, height*width, height*width)
	for row := 0; row < height; row++ {
		v[row] = rows[row*width : (row+1)*width]
	}
	return v
}
