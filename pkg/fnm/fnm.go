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

// Package fnm returns a new filename.
package fnm

import (
	"errors"
	"fmt"
	"os"
)

func UniqueName(kind string, seed int) string {
	name := fmt.Sprintf("%x-%s.png", seed, kind)
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return name
	}
	for i := 1; i < 1024; i++ {
		name = fmt.Sprintf("%x-%s-%04d.png", seed, kind, i)
		if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
			return name
		}
	}
	panic(fmt.Sprintf("too many attemps %q", name))
}
