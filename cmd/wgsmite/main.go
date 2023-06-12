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
	"log"
	"math/rand"
)

func main() {
	Seed := 0x638bb317ac47a6ba
	rand.Seed(int64(Seed))
	//rand.Seed(time.Now().UnixNano())

	height, width, iterations := 600, 1_200, 10_000
	saveFile := fnm.UniqueName("smite", Seed)
	if err := smite.Run(height, width, iterations, saveFile); err != nil {
		log.Fatal(err)
	}
}
