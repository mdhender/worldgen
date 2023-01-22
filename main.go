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
	"github.com/mdhender/worldgen/pkg/fractal"
	"github.com/mdhender/worldgen/pkg/tiled"
	"github.com/mdhender/worldgen/sliced"
	"log"
	"math/rand"
	"os"
)

func main() {
	Seed := 1812 // 0x638bb317ac47a6ba
	rand.Seed(int64(Seed))
	//rand.Seed(time.Now().UnixNano())

	doFractal, doSliced, doTiled := false, false, false
	height, width, iterations := 600, 1_200, 1_000
	var saveFile string
	for _, arg := range os.Args[1:] {
		if arg == "--fractal" {
			doFractal = true
			if saveFile == "" {
				saveFile = "fractal.png"
			}
		} else if arg == "--sliced" {
			doSliced = true
			if saveFile == "" {
				saveFile = "sliced.png"
			}
		} else if arg == "--tiled" {
			doTiled = true
			if saveFile == "" {
				saveFile = "tiled.png"
			}
		}
	}

	if doFractal {
		if err := fractal.Run(); err != nil {
			log.Fatal(err)
		}
	}
	if doSliced {
		if err := sliced.Run(height, width, iterations, saveFile); err != nil {
			log.Fatal(err)
		}
	}
	if doTiled {
		if err := tiled.Run(height, width, iterations, saveFile); err != nil {
			log.Fatal(err)
		}
	}
}
