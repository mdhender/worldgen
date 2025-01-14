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

import (
	"github.com/mdhender/worldgen/pkg/cmap"
	"image"
	"image/color"
)

// AsCarto assumes the map has been normalized to 0..255
func (m *Map) AsCarto(cm cmap.ColorMap) *image.RGBA {
	height, width := m.Height(), m.Width()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, cm[m.yx[y][x]])
		}
	}
	return img
}

// AsGreyscale assumes the map has been normalized to 0..255
func (m *Map) AsGreyscale() *image.RGBA {
	height, width := m.Height(), m.Width()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := uint8(m.yx[y][x])
			pc := color.RGBA{R: val, G: val, B: val, A: 255}
			img.Set(x, y, pc)
		}
	}
	return img
}

// AsImage assumes the map has been normalized to 0..255
func (m *Map) AsImage() *image.RGBA {
	height, width := m.Height(), m.Width()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// scale point from 0..255 to 0..49
			val := uint8(m.yx[y][x] * 49 / 255)
			pc := color.RGBA{R: red[val], G: green[val], B: blue[val], A: 255}
			img.Set(x, y, pc)
		}
	}
	return img
}
