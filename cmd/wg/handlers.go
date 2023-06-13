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

		var m gen.Map
		if data, err := os.ReadFile(fmt.Sprintf("%x-asteroids.json", seed)); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		} else if err = json.Unmarshal(data, &m); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		data := struct {
			Height, Width int
			Seed          string
		}{
			Height: m.Height(),
			Width:  m.Width(),
			Seed:   fmt.Sprintf("%x", seed),
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
