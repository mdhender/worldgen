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
	"bytes"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	XRange = 6400 // 320
	YRange = 3200 // 160

	YRangeDiv2  = float64(YRange) / 2
	YRangeDivPI = float64(YRange) / math.Pi
)

var (
	WorldMapArray []int
	Histogram     [256]int
	FilledPixels  int

	SinIterPhi []float64
)

func main() {
	started := time.Now()

	//argc, argv := len(os.Args), os.Args

	var Color int
	var i, j, row, TwoColorMode int
	var index2 int
	var Threshold int
	var Cur int
	var SaveName string /* 255 character filenames should be enough? */
	var SaveFile string /* SaveName + .gif */

	SinIterPhi = make([]float64, 2*XRange, 2*XRange)
	for i = 0; i < XRange; i++ {
		SinIterPhi[i] = math.Sin(float64(i) * 2 * math.Pi / float64(XRange))
		SinIterPhi[i+XRange] = SinIterPhi[i]
	}

	Seed := 0x638bb317ac47a6ba
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(int64(Seed))

	NumberOfFaults := 10_000
	PercentWater := 15
	PercentIce := 2
	SaveName = "olson"
	log.Printf("Seed: %d\n", Seed)
	log.Printf("Number of faults: %d\n", NumberOfFaults)
	log.Printf("Percent water: %d\n", PercentWater)
	log.Printf("Percent ice: %d\n", PercentIce)

	log.Printf("Save as (.GIF will be appended): %q\n", SaveName)

	//srand(Seed);

	WorldMapArray = make([]int, XRange*YRange, XRange*YRange)
	for j, row = 0, 0; j < XRange; j++ {
		WorldMapArray[row] = 0
		for i = 1; i < YRange; i++ {
			WorldMapArray[i+row] = math.MinInt
		}
		row += YRange
	}
	log.Printf("filled %d x %d world map: %v\n", XRange, YRange, time.Now().Sub(started))

	/* Generate the map! */
	for a := 0; a < NumberOfFaults; a++ {
		FractureWorldMap(rand.Intn(2) == 0)
	}
	log.Printf("generated %d faults: %v\n", NumberOfFaults, time.Now().Sub(started))

	/* Copy data (I have only calculated faults for 1/2 the image.
	 * I can do this due to symmetry... :) */
	index2 = (XRange / 2) * YRange
	for j, row = 0, 0; j < XRange/2; j++ {
		for i = 1; i < YRange; i++ { /* fix */
			WorldMapArray[row+index2+YRange-i] = WorldMapArray[row+i]
		}
		row += YRange
	}
	log.Printf("flipped the image: %v\n", time.Now().Sub(started))

	/* Reconstruct the real WorldMap from the WorldMapArray and FaultArray */
	for j, row = 0, 0; j < XRange; j++ {
		/* We have to start somewhere, and the top row was initialized to 0,
		 * but it might have changed during the iterations... */
		Color = WorldMapArray[row]
		for i := 1; i < YRange; i++ {
			/* We "fill" all positions with values != INT_MIN with Color */
			Cur = WorldMapArray[row+i]
			if Cur != math.MinInt {
				Color += Cur
			}
			WorldMapArray[row+i] = Color
		}
		row += YRange
	}
	log.Printf("rebuilt the world map: %v\n", time.Now().Sub(started))

	/* Compute MAX and MIN values in WorldMapArray */
	MinZ, MaxZ := -1, 1
	for j = 0; j < XRange*YRange; j++ {
		Color = WorldMapArray[j]
		if MinZ > Color {
			MinZ = Color
		}
		if MaxZ < Color {
			MaxZ = Color
		}
	}
	log.Printf("computed minz %d and maxz %d: %v\n", MinZ, MaxZ, time.Now().Sub(started))

	/* Compute color-histogram of WorldMapArray.
	 * This histogram is a very crude approximation, since all pixels are
	 * considered of the same size... I will try to change this in a
	 * later version of this program. */
	for j, row = 0, 0; j < XRange; j++ {
		for i = 0; i < YRange; i++ {
			Color = WorldMapArray[row+i]
			Color = int(((float64(Color-MinZ+1) / float64(MaxZ-MinZ+1)) * 30) + 1)
			Histogram[Color]++
		}
		row += YRange
	}
	log.Printf("computed histogram: %v\n", time.Now().Sub(started))

	/* Threshold now holds how many pixels PercentWater means */
	Threshold = PercentWater * XRange * YRange / 100
	log.Printf("threshold is %d: %v\n", Threshold, time.Now().Sub(started))

	/* "Integrate" the histogram to decide where to put sea-level */
	Count := 0
	for j = 0; j < 256 && Count <= Threshold; j++ {
		Count += Histogram[j]
	}
	log.Printf("integrated histogram %d %d / %d: %v\n", j, Count, Threshold, time.Now().Sub(started))

	/* Threshold now holds where sea-level is */
	Threshold = j*(MaxZ-MinZ+1)/30 + MinZ
	log.Printf("threshold is %d * (%d - %d + 1) / 30 + %d: %d: %v\n", j, MaxZ, MinZ, MinZ, Threshold, time.Now().Sub(started))

	if TwoColorMode != 0 {
		for j, row = 0, 0; j < XRange; j++ {
			for i = 0; i < YRange; i++ {
				Color = WorldMapArray[row+i]
				if Color < Threshold {
					WorldMapArray[row+i] = 3
				} else {
					WorldMapArray[row+i] = 20
				}
			}
			row += YRange
		}
		log.Printf("filled two color mode: %v\n", time.Now().Sub(started))
	} else {
		/* Scale WorldMapArray to color range in a way that gives you
		 * a certain Ocean/Land ratio */
		for j, row = 0, 0; j < XRange; j++ {
			for i = 0; i < YRange; i++ {
				Color = WorldMapArray[row+i]

				if Color < Threshold {
					Color = int(((float64(Color-MinZ) / float64(Threshold-MinZ)) * 15) + 1)
				} else {
					Color = int(((float64(Color-Threshold) / float64(MaxZ-Threshold)) * 15) + 16)
				}

				/* Just in case... I DON't want the GIF-saver to flip out! :) */
				if Color < 1 {
					Color = 1
				} else if Color > 32 { // was 255?
					Color = 31
				}
				WorldMapArray[row+i] = Color
			}
			row += YRange
		}
		log.Printf("scaled color range: %v\n", time.Now().Sub(started))

		/* "Recycle" Threshold variable, and, eh, the variable still has something
		 * like the same meaning... :) */
		Threshold = PercentIce * XRange * YRange / 100

		if Threshold <= 0 || Threshold > XRange*YRange {
			goto Finished
		}

		FilledPixels = 0
		/* i==y, j==x */
		for i = 0; i < YRange; i++ {
			for j, row = 0, 0; j < XRange; j++ {
				Color = WorldMapArray[row+i]
				if Color < 32 {
					FloodFill4(j, i, Color)
				}
				/* FilledPixels is a global variable which FloodFill4 modifies...
				 * I know it's ugly, but as it is now, this is a hack! :)
				 */
				if FilledPixels > Threshold {
					goto NorthPoleFinished
				}
				row += YRange
			}
		}
		log.Printf("filled pixels: %v\n", time.Now().Sub(started))

	NorthPoleFinished:
		FilledPixels = 0
		/* i==y, j==x */
		for i = YRange - 1; i > 0; i-- { /* fix */
			for j, row = 0, 0; j < XRange; j++ {
				Color = WorldMapArray[row+i]
				if Color < 32 {
					FloodFill4(j, i, Color)
				}
				/* FilledPixels is a global variable which FloodFill4 modifies...
				 * I know it's ugly, but as it is now, this is a hack! :)
				 */
				if FilledPixels > Threshold {
					goto Finished
				}
				row += YRange
			}
		}
		log.Printf("filled north pole: %v\n", time.Now().Sub(started))
	Finished:
	}
	log.Printf("finished map generation: %v\n", time.Now().Sub(started))

	/* append .gif to SaveFile */
	SaveFile = SaveName + ".gif"
	/* open binary SaveFile */
	Save := &bytes.Buffer{}

	/* Write GIF to savefile */
	//GIFEncode(Save, XRange, YRange, 1, 0, 8, Red, Green, Blue)
	CreateImage()
	log.Printf("created image: %v\n", time.Now().Sub(started))

	log.Printf("Map created, saved as %s (%d bytes).\n", SaveFile, len(Save.Bytes()))
}

/* 4-connective floodfill algorithm which I use for constructing the ice-caps.*/
func FloodFill4(x, y, OldColor int) {
	if WorldMapArray[x*YRange+y] == OldColor {
		if WorldMapArray[x*YRange+y] < 16 {
			WorldMapArray[x*YRange+y] = 32
		} else {
			WorldMapArray[x*YRange+y] += 17
		}

		FilledPixels++
		if y-1 > 0 {
			FloodFill4(x, y-1, OldColor)
		}
		if y+1 < YRange {
			FloodFill4(x, y+1, OldColor)
		}
		if x-1 < 0 {
			FloodFill4(XRange-1, y, OldColor) /* fix */
		} else {
			FloodFill4(x-1, y, OldColor)
		}

		if x+1 >= XRange { /* fix */
			FloodFill4(0, y, OldColor)
		} else {
			FloodFill4(x+1, y, OldColor)
		}
	}
}

/* Function that generates the worldmap */
func FractureWorldMap(lower bool) {
	/* Create a random greatcircle...
	 * Start with an equator and rotate it */
	Alpha := (rand64() - 0.5) * math.Pi /* Rotate around x-axis */
	Beta := (rand64() - 0.5) * math.Pi  /* Rotate around y-axis */

	TanB := math.Tan(math.Acos(math.Cos(Alpha) * math.Cos(Beta)))

	Xsi := int((float64(XRange)/2 - float64(XRange)/math.Pi) * Beta)

	var Phi, Theta, row int
	for Phi = 0; Phi < XRange/2; Phi++ {
		Theta = int((YRangeDivPI * math.Atan(SinIterPhi[Xsi-Phi+XRange]*TanB)) + YRangeDiv2)
		if lower {
			/* lower southern hemisphere */
			if WorldMapArray[row+Theta] != math.MinInt {
				WorldMapArray[row+Theta]--
			} else {
				WorldMapArray[row+Theta] = -1
			}
		} else {
			/* raise southern hemisphere */
			if WorldMapArray[row+Theta] != math.MinInt {
				WorldMapArray[row+Theta]++
			} else {
				WorldMapArray[row+Theta] = 1
			}
		}
		row += YRange
	}
}

func rand64() float64 {
	return rand.Float64()
}
