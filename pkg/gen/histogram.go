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

package gen

// Histogram assumes that the map has been normalized to 0..255
func (m *Map) Histogram() (hs [256]int) {
	for _, val := range m.points {
		hs[val] = hs[val] + 1
	}
	return hs
}

func (m *Map) IceLevel(pct int) int {
	return 200
}

func (m *Map) SeaLevel(pct int) int {
	threshold := pct * len(m.points) / 100
	if threshold <= 1 {
		return 1
	} else if threshold >= len(m.points) {
		return 254
	}

	// find the sea-level
	pixels := 0
	for n, val := range m.Histogram() {
		if pixels += val; pixels > threshold {
			return n
		}
	}

	return 254
}
