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
	"github.com/mdhender/worldgen/pkg/sliced"
	"github.com/mdhender/worldgen/pkg/smite"
	"github.com/mdhender/worldgen/pkg/tiled"
	"log"
	"math/rand"
	"os"
)

func main() {
	Seed := 0x638bb317ac47a6ba
	rand.Seed(int64(Seed))
	//rand.Seed(time.Now().UnixNano())

	doFractal, doSliced, doSmite, doTiled := false, false, false, false
	height, width, iterations := 600, 1_200, 10_000
	var saveFile string
	for _, arg := range os.Args[1:] {
		if arg == "--fractal" {
			doFractal = true
		} else if arg == "--sliced" {
			doSliced = true
		} else if arg == "--smited" {
			doSmite = true
		} else if arg == "--tiled" {
			doTiled = true
		}
	}

	if doFractal {
		if err := fractal.Run(); err != nil {
			log.Fatal(err)
		}
	}
	if doSliced {
		saveFile = "sliced.png"
		if err := sliced.Run(height, width, iterations, saveFile); err != nil {
			log.Fatal(err)
		}
	}
	if doSmite {
		saveFile = "smited.png"
		if err := smite.Run(height, width, iterations, saveFile); err != nil {
			log.Fatal(err)
		}
	}
	if doTiled {
		saveFile = "tiled.png"
		if err := tiled.Run(height, width, iterations, saveFile); err != nil {
			log.Fatal(err)
		}
	}
}
