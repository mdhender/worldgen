/*
 * worldgen - fractured terrain generator
 * Copyright (C) 1999  John Olsson
 * Copyright (C) 2022 Michael D Henderson
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

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

var (
	Red = [255]uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 34, 68, 102, 119, 136, 153, 170, 187,
		0, 34, 34, 119, 187, 255, 238, 221, 204, 187, 170, 153,
		136, 119, 85, 68,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Green = [255]uint8{
		0, 0, 17, 51, 85, 119, 153, 204, 221, 238, 255, 255, 255,
		255, 255, 255, 68, 102, 136, 170, 221, 187, 170, 136,
		136, 102, 85, 85, 68, 51, 51, 34,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Blue = [255]uint8{
		0, 68, 102, 136, 170, 187, 221, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 0, 0, 0, 0, 0, 34, 34, 34, 34, 34, 34,
		34, 34, 34, 17, 0,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func ColorCard(borg bool) {
	width, height := 1280, 640
	cardsPerRow := width / 50
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		c := x / cardsPerRow
		pc := color.RGBA{R: Red[c], G: Green[c], B: Blue[c], A: 255}
		for y := 0; y < height; y++ {
			m.Set(x, y, pc)
		}
	}
	outFile, err := os.Create("color-card.png")
	if err != nil {
		log.Fatal(err)
	}
	_ = png.Encode(outFile, m)
	_ = outFile.Close()
	if borg {
		os.Exit(0)
	}
}

func Project(projectionType int, xRange, yRange int, scrollDegrees int) *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, xRange, yRange))
	switch projectionType {
	case SPHERICAL:
		projectSpherical(m, Height, scrollDegrees)
	default:
		projectSquare(m, xRange, yRange)
	}
	return m
}

func projectSpherical(m *image.RGBA, height int, scrollDegrees int) {
	diameter := height
	if 2*(height/2)-height != 0 {
		diameter++
	}
	radius := diameter / 2
	rSquared := radius * radius
	xRange, yRange := diameter, diameter
	xRange64, yRange64 := float64(xRange), float64(yRange)

	scrollDistance := -1 * int((float64(scrollDegrees%360))*(float64(XRange)/360))

	for x := 0; x < xRange; x++ {
		for y := 0; y < yRange; y++ {
			newX, newY := x-radius, y-radius

			temp := newX*newX + newY*newY
			if temp > rSquared {
				// point is not within the circle
				continue
			}
			sphereX := math.Sqrt(float64(rSquared - temp))

			newX = (int)((math.Atan((float64(newX))/sphereX)*xRange64/math.Pi+xRange64)/2) + scrollDistance
			if newX < 0 {
				newX = XRange - 1 - ((-1 * newX) % XRange)
			}
			if newX >= XRange {
				newX = newX % XRange
			}

			theta := math.Acos((float64(newY)) / (float64(radius)))
			newY = YRange - int(theta*yRange64/math.Pi)

			height := WorldMapArray[newX*YRange+newY]
			pc := color.RGBA{R: Red[height], G: Green[height], B: Blue[height], A: 255}
			m.Set(x, y, pc)
		}
	}
}

func projectSquare(m *image.RGBA, xRange, yRange int) {
	for x := 0; x < xRange; x++ {
		for y := 0; y < yRange; y++ {
			height := WorldMapArray[x*yRange+y]
			pc := color.RGBA{R: Red[height], G: Green[height], B: Blue[height], A: 255}
			m.Set(x, y, pc)
		}
	}
}
