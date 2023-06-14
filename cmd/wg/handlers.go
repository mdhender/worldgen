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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mdhender/worldgen/pkg/cmap"
	"github.com/mdhender/worldgen/pkg/gen"
	"github.com/mdhender/worldgen/pkg/way"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("..", "templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func cartoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var seed uint64

		started := time.Now()
		defer func() {
			log.Printf("%s %q elapsed %v\n", r.Method, r.URL, time.Now().Sub(started))
		}()

		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		var err error
		seed, err = strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		// get shiftX and shiftY from query parameters
		var shiftX, shiftY int
		if qShiftX := r.URL.Query()["shiftX"]; len(qShiftX) == 0 {
			shiftX = 0
		} else if len(qShiftX) > 1 {
			http.Error(w, "shiftX repeated", http.StatusBadRequest)
			return
		} else if shiftX, err = strconv.Atoi(qShiftX[0]); err != nil {
			http.Error(w, fmt.Sprintf("shiftX %v", err), http.StatusBadRequest)
			return
		} else if shiftX < 0 || shiftX > 99 {
			http.Error(w, "shiftX out of range", http.StatusBadRequest)
			return
		}
		if qShiftY := r.URL.Query()["shiftY"]; len(qShiftY) == 0 {
			shiftY = 0
		} else if len(qShiftY) > 1 {
			http.Error(w, "shiftY repeated", http.StatusBadRequest)
			return
		} else if shiftY, err = strconv.Atoi(qShiftY[0]); err != nil {
			http.Error(w, fmt.Sprintf("shiftY %v", err), http.StatusBadRequest)
			return
		} else if shiftY < 0 || shiftY > 99 {
			http.Error(w, "shiftY out of range", http.StatusBadRequest)
			return
		}

		// get pctWater from query parameters
		var pctWater int
		if qPctWater := r.URL.Query()["pctWater"]; len(qPctWater) == 0 {
			http.Error(w, "pctWater missing", http.StatusBadRequest)
			return
		} else if len(qPctWater) > 1 {
			http.Error(w, "pctWater repeated", http.StatusBadRequest)
			return
		} else if pctWater, err = strconv.Atoi(qPctWater[0]); err != nil {
			http.Error(w, fmt.Sprintf("pctWater %v", err), http.StatusBadRequest)
			return
		} else if pctWater < 1 || pctWater > 255 {
			http.Error(w, "pctWater out of range", http.StatusBadRequest)
			return
		}

		// get pctIce from query parameters
		var pctIce int
		if qPctIce := r.URL.Query()["pctIce"]; len(qPctIce) == 0 {
			http.Error(w, "pctIce missing", http.StatusBadRequest)
			return
		} else if len(qPctIce) > 1 {
			http.Error(w, "pctIce repeated", http.StatusBadRequest)
			return
		} else if pctIce, err = strconv.Atoi(qPctIce[0]); err != nil {
			http.Error(w, fmt.Sprintf("pctIce %v", err), http.StatusBadRequest)
			return
		} else if pctIce < 1 || pctIce > 255 {
			http.Error(w, "pctIce out of range", http.StatusBadRequest)
			return
		}

		var m gen.Map
		if data, err := os.ReadFile(fmt.Sprintf("%x-asteroids.json", seed)); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if err = json.Unmarshal(data, &m); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}
		if shiftX != 0 && shiftY != 0 {
			m.Shift(m.Width()*shiftX/100, m.Height()*shiftY/100)
		} else if shiftX != 0 {
			m.Shift(m.Width()*shiftX/100, 0)
		} else if shiftY != 0 {
			m.Shift(0, m.Height()*shiftY/100)
		}

		// generate color map
		cm := cmap.FromHistogram(m.Histogram(), pctWater, pctIce, cmap.Water, cmap.Terrain, cmap.Ice)

		png, err := m.AsPNG(m.AsCarto(cm))
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(png)
	}
}

func customizeHandler(root, filename string) http.HandlerFunc {
	templ := template.Must(template.ParseFiles(filepath.Join(root, filename)))
	return func(w http.ResponseWriter, r *http.Request) {
		var seed uint64

		started := time.Now()
		defer func() {
			log.Printf("%s %q elapsed %v\n", r.Method, r.URL, time.Now().Sub(started))
		}()

		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		var err error
		seed, err = strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		// get shiftX and shiftY from query parameters
		var shiftX, shiftY int
		if qShiftX := r.URL.Query()["shiftX"]; len(qShiftX) == 0 {
			shiftX = 0
		} else if len(qShiftX) > 1 {
			http.Error(w, "shiftX repeated", http.StatusBadRequest)
			return
		} else if shiftX, err = strconv.Atoi(qShiftX[0]); err != nil {
			http.Error(w, fmt.Sprintf("shiftX %v", err), http.StatusBadRequest)
			return
		} else if shiftX < 0 || shiftX > 99 {
			http.Error(w, "shiftX out of range", http.StatusBadRequest)
			return
		} else {
			log.Printf("shiftX %v %d\n", qShiftX, shiftX)
		}
		if qShiftY := r.URL.Query()["shiftY"]; len(qShiftY) == 0 {
			shiftY = 0
		} else if len(qShiftY) > 1 {
			http.Error(w, "shiftY repeated", http.StatusBadRequest)
			return
		} else if shiftY, err = strconv.Atoi(qShiftY[0]); err != nil {
			http.Error(w, fmt.Sprintf("shiftY %v", err), http.StatusBadRequest)
			return
		} else if shiftY < 0 || shiftY > 99 {
			http.Error(w, "shiftY out of range", http.StatusBadRequest)
			return
		} else {
			log.Printf("shiftY %v %d\n", qShiftY, shiftY)
		}

		// get pctWater from query parameters
		var pctWater int
		if qPctWater := r.URL.Query()["pctWater"]; len(qPctWater) == 0 {
			pctWater = 50 // default value
		} else if len(qPctWater) > 1 {
			http.Error(w, "pctWater repeated", http.StatusBadRequest)
			return
		} else if pctWater, err = strconv.Atoi(qPctWater[0]); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if pctWater < 1 || pctWater > 99 {
			http.Error(w, "pctWater out of range", http.StatusBadRequest)
			return
		}

		// get pctIce from query parameters
		var pctIce int
		if qPctWater := r.URL.Query()["pctIce"]; len(qPctWater) == 0 {
			pctIce = 50 // default value
		} else if len(qPctWater) > 1 {
			http.Error(w, "pctIce repeated", http.StatusBadRequest)
			return
		} else if pctIce, err = strconv.Atoi(qPctWater[0]); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if pctIce < 1 || pctIce > 99 {
			http.Error(w, "pctIce out of range", http.StatusBadRequest)
			return
		}

		var m gen.Map
		if data, err := os.ReadFile(fmt.Sprintf("%x-asteroids.json", seed)); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if err = json.Unmarshal(data, &m); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		data := struct {
			Height, Width      int
			Seed               string
			Histogram          [256]string
			ShiftX, ShiftY     int
			PctWater, SeaLevel int
			PctIce, IceLevel   int
		}{
			Height:   m.Height(),
			Width:    m.Width(),
			Seed:     fmt.Sprintf("%x", seed),
			ShiftX:   shiftX,
			ShiftY:   shiftY,
			PctWater: pctWater,
			SeaLevel: m.SeaLevel(pctWater),
			PctIce:   pctIce,
			IceLevel: m.IceLevel(pctIce),
		}
		pixels, count := m.Height()*m.Width(), 0
		for n, val := range m.Histogram() {
			count += val
			data.Histogram[n] = fmt.Sprintf("%03d: %8d / %8d / %8d\n", n, val, count, pixels)
		}

		bb := &bytes.Buffer{}
		if err = templ.Execute(bb, data); err != nil {
			log.Printf("%s %q %v\n", r.Method, r.URL, err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bb.Bytes())
	}
}

func greyscaleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var seed uint64

		started := time.Now()
		defer func() {
			log.Printf("%s %q elapsed %v\n", r.Method, r.URL, time.Now().Sub(started))
		}()

		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		var err error
		seed, err = strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		var m gen.Map
		if data, err := os.ReadFile(fmt.Sprintf("%x-asteroids.json", seed)); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if err = json.Unmarshal(data, &m); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		png, err := m.AsPNG(m.AsGreyscale())
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(png)
	}
}

func indexHandler(root, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
