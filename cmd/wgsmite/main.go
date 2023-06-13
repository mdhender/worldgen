// worldgen - fractured terrain generator
// Copyright (c) 2022-2023 Michael D Henderson
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

package main

import (
	"github.com/mdhender/worldgen/pkg/fnm"
	"github.com/mdhender/worldgen/pkg/smite"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	Seed := 0x638bb317ac47a6ba
	rand.Seed(int64(Seed))
	//rand.Seed(time.Now().UnixNano())

	started := time.Now()

	height, width, iterations := 600, 1_200, 10_000
	saveFile := fnm.UniqueName("smite", Seed)

	img, err := smite.Generate(height, width, iterations)
	if err != nil {
		log.Fatal(err)
	}

	err = savePNG(saveFile, img)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("smite: created %s: %v\n", saveFile, time.Now().Sub(started))
}

func savePNG(filename string, m *image.RGBA) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(outFile, m); err != nil {
		return err
	}
	return outFile.Close()
}
