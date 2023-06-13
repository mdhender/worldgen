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

var (
	red = []uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 34, 68, 102, 119, 136, 153, 170, 187,
		0, 34, 34, 119, 187, 255, 238, 221, 204, 187, 170, 153,
		136, 119, 85, 68,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
	green = []uint8{
		0, 0, 17, 51, 85, 119, 153, 204, 221, 238, 255, 255, 255,
		255, 255, 255, 68, 102, 136, 170, 221, 187, 170, 136,
		136, 102, 85, 85, 68, 51, 51, 34,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
	blue = []uint8{
		0, 68, 102, 136, 170, 187, 221, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 0, 0, 0, 0, 0, 34, 34, 34, 34, 34, 34,
		34, 34, 34, 17, 0,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200,
		195, 190, 185, 180, 175, 175}
)
