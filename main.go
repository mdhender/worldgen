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
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	Height = 320

	SQUARE           = 0
	MERCATOR         = 1
	SPHERICAL        = 2
	ORTHOGRAPHIC_NP  = 3
	ORTHOGRAPHIC_SP  = 4
	STEREOGRAPHIC_NP = 5
	STEREOGRAPHIC_SP = 6
	GNOMIC_NP        = 7
	GNOMIC_SP        = 8
	LAMBERT_AREAP_NP = 9
	LAMBERT_AREAP_SP = 10
	KACHUNK          = 11

	HEIGHT = 0
	RADIUS = 1
)

var (
	showGlobe = false

	colorMap     [][]int
	heightMap    [][]int
	Histogram    [256]int
	FilledPixels int

	SinIterPhi []float64

	ProjectionType = SQUARE
	ScrollDegrees  = 0
	ScrollDistance int

	XRange = 6400 // 320
	YRange = 3200 // 160

	YRangeDiv2  float64
	YRangeDivPI float64
)

func main() {
	started := time.Now()

	// save color card
	ColorCard(false)

	//argc, argv := len(os.Args), os.Args

	TwoColorMode := false
	var Color int
	var i, j, row int
	var Threshold int

	//Mode := HEIGHT

	switch ProjectionType {
	case KACHUNK:
		YRange = Height
		XRange = 2 * YRange
	case SQUARE:
		YRange = Height
		XRange = 2 * YRange
	case MERCATOR:
		YRange = Height
		XRange = int(float64(YRange) * math.Pi / 2)
		if 2*(XRange/2)-XRange != 0 {
			XRange++
		}
	case SPHERICAL:
		YRange = int(math.Round(Height * math.Pi / 2))
		XRange = int(math.Round(Height * math.Pi))
	case ORTHOGRAPHIC_NP, ORTHOGRAPHIC_SP:
		YRange = int(math.Round(Height * math.Pi / 2))
		XRange = int(math.Round(Height * math.Pi))
	case STEREOGRAPHIC_NP, STEREOGRAPHIC_SP:
		YRange = Height
		XRange = int(math.Round(Height * math.Pi))
	case GNOMIC_NP, GNOMIC_SP:
		YRange = int(math.Round(Height * math.Pi / 2))
		XRange = int(math.Round(Height * math.Pi))
	case LAMBERT_AREAP_NP, LAMBERT_AREAP_SP:
		YRange = Height
		XRange = int(math.Round(Height * math.Pi))
	}
	// cache some frequently used values based on the size of the world
	YRangeDiv2 = float64(YRange) / 2
	YRangeDivPI = float64(YRange) / math.Pi

	switch ProjectionType {
	case SQUARE, MERCATOR:
		log.Printf("WIDTH=%d HEIGHT=%d\n", XRange, YRange)
	default:
		log.Printf("WIDTH=%d HEIGHT=%d\n", Height, Height)
	}

	SinIterPhi = make([]float64, 2*XRange, 2*XRange)
	for i = 0; i < XRange; i++ {
		SinIterPhi[i] = math.Sin(float64(i) * 2 * math.Pi / float64(XRange))
		SinIterPhi[i+XRange] = SinIterPhi[i]
	}

	Seed := 0x638bb317ac47a6ba
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(int64(Seed))

	NumberOfFaults := 100
	PercentWater := 10
	PercentIce := 10
	log.Printf("Seed: %d\n", Seed)
	log.Printf("Number of faults: %d\n", NumberOfFaults)
	log.Printf("Percent water: %d\n", PercentWater)
	log.Printf("Percent ice: %d\n", PercentIce)

	/* Threshold now holds how many pixels PercentWater means */
	Threshold = PercentWater * XRange * YRange / 100
	log.Printf("threshold is %d: %v\n", Threshold, time.Now().Sub(started))

	heightMap = twoDimensionalArray(XRange, YRange)
	for row = 0; row < len(heightMap); row++ {
		heightMap[row][0] = 0 // why store zero here?
		for col := 1; col < len(heightMap[row]); col++ {
			heightMap[row][col] = math.MinInt
		}
	}
	log.Printf("filled %d x %d world map: %v\n", XRange, YRange, time.Now().Sub(started))

	/* Generate the map! */
	hasSymmetry := false
	switch ProjectionType {
	case KACHUNK:
		for a := 0; a < NumberOfFaults; a++ {
			FractureKachunk(XRange, YRange)
		}
	case MERCATOR:
		hasSymmetry = true
		for a := 0; a < NumberOfFaults; a++ {
			GenerateMercatorWorldMap(rand.Intn(2) == 0)
		}
	default:
		hasSymmetry = true
		for a := 0; a < NumberOfFaults; a++ {
			GenerateSquareWorldMap(rand.Intn(2) == 0)
		}
	}
	log.Printf("generated %d faults: %v\n", NumberOfFaults, time.Now().Sub(started))

	if hasSymmetry {
		// copy data. the generator calculated faults for 1/2 the image.
		for row = 0; row < len(heightMap); row++ {
			for col := 1; col < XRange/2; col++ {
				heightMap[row][XRange-col] = heightMap[row][col]
			}
		}
		log.Printf("flipped the image: %v\n", time.Now().Sub(started))
	}

	/* Reconstruct the real WorldMap from the WorldMapArray and FaultArray */
	for row = 0; row < len(heightMap); row++ {
		/* We have to start somewhere, and the top row was initialized to 0,
		 * but it might have changed during the iterations... */
		firstColValue := heightMap[row][0]
		for col := 1; col < len(heightMap[row]); col++ {
			/* We "fill" all positions with values != INT_MIN with firstColValue */
			cur := heightMap[row][col]
			if cur != math.MinInt {
				firstColValue += cur
			}
			heightMap[row][col] = firstColValue
		}
	}
	log.Printf("rebuilt the world map: %v\n", time.Now().Sub(started))

	/* Compute MAX and MIN values in WorldMapArray */
	MinZ, MaxZ := -1, 1
	for row = 0; row < len(heightMap); row++ {
		for col := 0; col < len(heightMap[row]); col++ {
			z := heightMap[row][col]
			if z < MinZ {
				MinZ = z
			}
			if z > MaxZ {
				MaxZ = z
			}
		}
	}
	log.Printf("computed minz %d and maxz %d: %v\n", MinZ, MaxZ, time.Now().Sub(started))

	/* Compute color-histogram of WorldMapArray. */
	rangeHeight := float64(MaxZ - MinZ + 1)
	for row = 0; row < len(heightMap); row++ {
		for col := 0; col < len(heightMap[row]); col++ {
			normalizedHeight := float64(heightMap[row][col] - MinZ + 1)
			Histogram[int(((normalizedHeight/rangeHeight)*30)+1)]++
		}
	}
	log.Printf("computed histogram: %v\n", time.Now().Sub(started))
	for k := 0; k < len(Histogram); k += 8 {
		log.Printf(" %7d %7d %7d %7d %7d %7d %7d %7d\n",
			Histogram[k+0], Histogram[k+1], Histogram[k+2], Histogram[k+3], Histogram[k+4], Histogram[k+5], Histogram[k+6], Histogram[k+7])
	}

	/* "Integrate" the histogram to decide where to put sea-level */
	Count := 0
	for j = 0; j < 256; j++ {
		Count += Histogram[j]
		if Count > Threshold {
			break
		}
	}
	log.Printf("integrated histogram %d %d / %d: %v\n", j, Count, Threshold, time.Now().Sub(started))

	/* Threshold now holds where sea-level is */
	Threshold = j*(MaxZ-MinZ+1)/30 + MinZ
	log.Printf("threshold is %d * (%d - %d + 1) / 30 + %d: %d: %v\n", j, MaxZ, MinZ, MinZ, Threshold, time.Now().Sub(started))

	colorMap = twoDimensionalArray(XRange, YRange)
	if TwoColorMode {
		for row = 0; row < len(heightMap); row++ {
			for col := 0; col < len(heightMap[row]); col++ {
				Color = heightMap[row][col]
				if Color < Threshold {
					heightMap[row][col] = 3
				} else {
					heightMap[row][col] = 20
				}
				colorMap[row][col] = heightMap[row][col]
			}
		}
		log.Printf("filled two color mode: %v\n", time.Now().Sub(started))
	} else {
		/* Scale WorldMapArray to color range in a way that gives you a certain Ocean/Land ratio */
		for row = 0; row < len(heightMap); row++ {
			for col := 0; col < len(heightMap[row]); col++ {
				if heightMap[row][col] < Threshold {
					Color = int(((float64(Color-MinZ) / float64(Threshold-MinZ)) * 15) + 1)
				} else {
					Color = int(((float64(Color-Threshold) / float64(MaxZ-Threshold)) * 15) + 16)
				}

				/* Just in case... I DON't want the GIF-saver to flip out! :) */
				if Color < 1 {
					Color = 1
				} else if Color > 255 {
					Color = 255
				}
				heightMap[row][col] = Color
				colorMap[row][col] = heightMap[row][col]
			}
		}
		log.Printf("scaled color range: %v\n", time.Now().Sub(started))

		/* "Recycle" Threshold variable, and, eh, the variable still has something
		 * like the same meaning... :) */
		Threshold = PercentIce * XRange * YRange / 100
		if 0 < Threshold && Threshold <= XRange*YRange {
			// fill from the "north"?
			filledPixels := 0
			for row = 0; row < len(colorMap) && filledPixels < Threshold; row++ {
				for col := 0; col < len(colorMap[row]) && filledPixels < Threshold; col++ {
					if colorMap[row][col] < 32 {
						filledPixels += FloodFill4(col, row, colorMap[row][col])
					}
				}
			}
			log.Printf("filled %8d/%8d north pixels: %v\n", filledPixels, Threshold, time.Now().Sub(started))

			// fill from the "south"?
			filledPixels = 0
			/* i==y, j==x */
			for row = 0; row < len(colorMap) && filledPixels < Threshold; row++ {
				for col := len(colorMap[row]) - 1; col >= 0 && filledPixels < Threshold; col-- {
					if colorMap[row][col] < 32 {
						filledPixels += FloodFill4(col, row, colorMap[row][col])
					}
				}
			}
			log.Printf("filled %8d/%8d south pixels: %v\n", filledPixels, Threshold, time.Now().Sub(started))
		}
	}
	log.Printf("finished map generation: %v\n", time.Now().Sub(started))

	/* Somehow, this seems to be the easy way of patching the problem of scrolling the wrong direction... ;) */
	ScrollDegrees = 33
	ScrollDistance = -1 * int((float64(ScrollDegrees%360))*(float64(XRange)/360))

	// create map
	var m *image.RGBA
	var saveFile string
	switch ProjectionType {
	case ORTHOGRAPHIC_NP, ORTHOGRAPHIC_SP, STEREOGRAPHIC_NP, STEREOGRAPHIC_SP, GNOMIC_NP, GNOMIC_SP, LAMBERT_AREAP_NP, LAMBERT_AREAP_SP:
		/*
		 * If it's a spherical projection, it will be a square map we output.
		 */
		Diameter := Height
		if 2*(Height/2)-Height != 0 {
			Diameter++
		}
		saveFile = "other"
		m = Project(ProjectionType, Diameter, Diameter, ScrollDegrees)
	case KACHUNK:
		saveFile = "kachunk"
		m = Project(ProjectionType, XRange, YRange, 0)
	case MERCATOR:
		saveFile = "mercator"
		m = Project(ProjectionType, XRange, YRange, 0)
	case SPHERICAL:
		// spherical projection, but still a square map on output
		Diameter := Height
		if 2*(Height/2)-Height != 0 {
			Diameter++
		}
		saveFile = "spherical"
		if showGlobe {
			for degrees := 0; degrees < 360; degrees = degrees + 15 {
				m = Project(ProjectionType, Diameter, Diameter, degrees)
				outFile, err := os.Create(fmt.Sprintf("%s-%03d.png", saveFile, degrees))
				if err != nil {
					log.Fatal(err)
				}
				_ = png.Encode(outFile, m)
				_ = outFile.Close()
			}
		}
		m = Project(ProjectionType, Diameter, Diameter, ScrollDegrees)
	case SQUARE:
		saveFile = "square"
		m = Project(ProjectionType, XRange, YRange, 0)
	default:
		saveFile = "rectangle"
		m = Project(ProjectionType, XRange, YRange, 0)
	}
	saveFile = fmt.Sprintf("%x-%s.png", Seed, saveFile)
	outFile, err := os.Create(saveFile)
	if err != nil {
		log.Fatal(err)
	}
	_ = png.Encode(outFile, m)
	_ = outFile.Close()

	log.Printf("created image: %s: %v\n", saveFile, time.Now().Sub(started))
}

/* 4-connective floodfill algorithm which I use for constructing the ice-caps.*/
func FloodFill4(x, y, oldColor int) int {
	filledPixels := 0
	if colorMap[y][x] == oldColor {
		if colorMap[y][x] < 16 {
			colorMap[y][x] = 32
		} else {
			colorMap[y][x] += 17
		}

		filledPixels++
		if y-1 > 0 {
			filledPixels += FloodFill4(x, y-1, oldColor)
		}
		if y+1 < YRange {
			filledPixels += FloodFill4(x, y+1, oldColor)
		}
		if x-1 < 0 {
			filledPixels += FloodFill4(XRange-1, y, oldColor) /* fix */
		} else {
			filledPixels += FloodFill4(x-1, y, oldColor)
		}
		if x+1 >= XRange { /* fix */
			filledPixels += FloodFill4(0, y, oldColor)
		} else {
			filledPixels += FloodFill4(x+1, y, oldColor)
		}
	}
	return filledPixels
}

/* Function that generates the worldmap */
func FractureKachunk(width, height int) {
	// decide the amount that we're going to raise or lower
	bump := 1
	if rand.Intn(2) == 0 {
		bump = -1
	}

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
	log.Printf("kachunk: line y = %f x + %f: bump %d\n", m, b, bump)

	// move all the points below the line up or down
	for x := 0; x < width; x++ {
		xl := float64(x)
		for y := 0; y < height; y++ {
			yl := m*xl + b
			if float64(y)-yl > 0 {
				// point is below the line
				if heightMap[y][x] == math.MinInt {
					heightMap[y][x] = 0
				}
				heightMap[y][x] += bump
			}
		}
	}
}

/* Function that generates the worldmap */
func GenerateMercatorWorldMap(lower bool) {
	// decide the amount that we're going to raise or lower
	bump := 1
	if lower {
		bump = -1
	}

	/* Create a random greatcircle...
	 * Start with an equator and rotate it */
	alpha := (rand64() - 0.5) * math.Pi /* Rotate around x-axis */
	beta := (rand64() - 0.5) * math.Pi  /* Rotate around y-axis */

	tanB := math.Tan(math.Acos(math.Cos(alpha) * math.Cos(beta)))

	xsi := int((float64(XRange)/2 - float64(XRange)/math.Pi) * beta)

	for row, phi := 0, 0; phi < XRange/2; phi++ {
		theta := int((math.Tan(math.Atan(SinIterPhi[xsi-phi+XRange]*tanB)/2) * YRangeDiv2) + YRangeDiv2)
		if heightMap[theta][phi] == math.MinInt {
			heightMap[theta][phi] = 0
		} else {
			heightMap[theta][phi] += bump
		}
		row++
	}
}

/* Function that generates the worldmap */
func GenerateSquareWorldMap(lower bool) {
	// decide the amount that we're going to raise or lower
	bump := 1
	if lower {
		bump = -1
	}

	/* Create a random greatcircle...
	 * Start with an equator and rotate it */
	alpha := (rand64() - 0.5) * math.Pi /* Rotate around x-axis */
	beta := (rand64() - 0.5) * math.Pi  /* Rotate around y-axis */

	tanB := math.Tan(math.Acos(math.Cos(alpha) * math.Cos(beta)))

	xsi := int((float64(XRange)/2 - float64(XRange)/math.Pi) * beta)

	for row, phi := 0, 0; phi < XRange/2; phi++ {
		theta := int((YRangeDivPI * math.Atan(SinIterPhi[xsi-phi+XRange]*tanB)) + YRangeDiv2)
		if heightMap[row][theta] == math.MinInt {
			heightMap[row][theta] = 0
		} else {
			heightMap[row][theta] += bump
		}
		row++
	}
}

func rand64() float64 {
	return rand.Float64()
}

func twoDimensionalArray(width, height int) [][]int {
	a := make([][]int, height, height)
	for y := 0; y < height; y++ {
		a[y] = make([]int, width, width)
	}
	return a
}
