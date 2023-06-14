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
		0, 34, 34, 119, 187, 255, 238, 221, 204, 187, 170, 153, 136, 119, 85, 68,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200, 195, 190, 185, 180, 175, 175}
	green = []uint8{
		0, 0, 17, 51, 85, 119, 153, 204, 221, 238, 255, 255, 255, 255, 255, 255,
		68, 102, 136, 170, 221, 187, 170, 136, 136, 102, 85, 85, 68, 51, 51, 34,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200, 195, 190, 185, 180, 175, 175}
	blue = []uint8{
		0, 68, 102, 136, 170, 187, 221, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		0, 0, 0, 0, 0, 34, 34, 34, 34, 34, 34, 34, 34, 34, 17, 0,
		255, 250, 245, 240, 235, 230, 225, 220, 215, 210, 205, 200, 195, 190, 185, 180, 175, 175}
)

type colormap struct {
	red, green, blue uint8
}

var (
	ice = []colormap{
		{red: 175, green: 175, blue: 175},
		{red: 180, green: 180, blue: 180},
		{red: 185, green: 185, blue: 185},
		{red: 190, green: 190, blue: 190},
		{red: 195, green: 195, blue: 195},
		{red: 200, green: 200, blue: 200},
		{red: 205, green: 205, blue: 205},
		{red: 210, green: 210, blue: 210},
		{red: 215, green: 215, blue: 215},
		{red: 220, green: 220, blue: 220},
		{red: 225, green: 225, blue: 225},
		{red: 230, green: 230, blue: 230},
		{red: 235, green: 235, blue: 235},
		{red: 240, green: 240, blue: 240},
		{red: 245, green: 245, blue: 245},
		{red: 250, green: 250, blue: 250},
		{red: 255, green: 255, blue: 255},
	}
	terrain = []colormap{
		{red: 0, green: 68, blue: 0},
		{red: 34, green: 102, blue: 0},
		{red: 34, green: 136, blue: 0},
		{red: 119, green: 170, blue: 0},
		{red: 187, green: 221, blue: 0},
		{red: 255, green: 187, blue: 34},
		{red: 238, green: 170, blue: 34},
		{red: 221, green: 136, blue: 34},
		{red: 204, green: 136, blue: 34},
		{red: 187, green: 102, blue: 34},
		{red: 170, green: 85, blue: 34},
		{red: 153, green: 85, blue: 34},
		{red: 136, green: 68, blue: 34},
		{red: 119, green: 51, blue: 34},
		{red: 85, green: 51, blue: 17},
		{red: 68, green: 34, blue: 0},
	}
	water = []colormap{
		{red: 0, green: 0, blue: 0},
		{red: 0, green: 0, blue: 68},
		{red: 0, green: 17, blue: 102},
		{red: 0, green: 51, blue: 136},
		{red: 0, green: 85, blue: 170},
		{red: 0, green: 119, blue: 187},
		{red: 0, green: 153, blue: 221},
		{red: 0, green: 204, blue: 255},
		{red: 34, green: 221, blue: 255},
		{red: 68, green: 238, blue: 255},
		{red: 102, green: 255, blue: 255},
		{red: 119, green: 255, blue: 255},
		{red: 136, green: 255, blue: 255},
		{red: 153, green: 255, blue: 255},
		{red: 170, green: 255, blue: 255},
		{red: 187, green: 255, blue: 255},
	}
)
