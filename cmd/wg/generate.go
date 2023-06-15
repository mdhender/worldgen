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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mdhender/worldgen/pkg/cmap"
	"github.com/mdhender/worldgen/pkg/gen"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func generateHandler(height, width, iterations int) http.HandlerFunc {
	var lock sync.Mutex
	secret := os.Getenv("WMG_SECRET")

	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		log.Printf("%s %s: entering\n", r.Method, r.URL)
		lock.Lock()
		defer func() {
			lock.Unlock()
			log.Printf("%s %s: elapsed %vn", r.Method, r.URL, time.Now().Sub(started))
		}()

		if err := r.ParseForm(); err != nil {
			log.Printf("%s %s: %v\n", r.Method, r.URL, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		log.Printf("%s %s: %v\n", r.Method, r.URL, r.PostForm)

		// get form values
		var err error
		var input struct {
			fname            string
			seed             uint64
			height, width    int
			iterations       int
			pctWater, pctIce int
			shiftX, shiftY   int
			secret           string
		}
		input.height = height
		input.width = width
		input.iterations = iterations
		if input.seed, err = pfvAsUint(r, "seed"); err != nil {
		} else if input.pctIce, err = pfvAsInt(r, "pct_ice"); err != nil {
		} else if input.pctWater, err = pfvAsInt(r, "pct_water"); err != nil {
		} else if input.shiftX, err = pfvAsInt(r, "shift_x"); err != nil {
		} else if input.shiftY, err = pfvAsInt(r, "shift_y"); err != nil {
		} else if input.secret, _ = pfvAsString(r, "secret"); err != nil {
		} else {
			input.fname = fmt.Sprintf("%x.json", input.seed)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}
		log.Printf("%s %s: %+v\n", r.Method, r.URL, input)

		var m *gen.Map

		// does map already exist?
		data, err := os.ReadFile(input.fname)
		if err == nil {
			// use it
			m = &gen.Map{}
			if err = json.Unmarshal(data, m); err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			}
		} else {
			hh := sha1.New()
			hh.Write([]byte(input.secret))
			sis := base64.URLEncoding.EncodeToString(hh.Sum(nil))
			hh = sha1.New()
			hh.Write([]byte(secret))
			sss := base64.URLEncoding.EncodeToString(hh.Sum(nil))
			if sis != sss {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// generate it
			m = gen.New(input.height, input.width, rand.New(rand.NewSource(int64(input.seed))))
			m.RandomFractureCircle(input.iterations)
			m.Normalize()

			// save it
			data, err := json.Marshal(m)
			if err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			} else if err = os.WriteFile(input.fname, data, 0644); err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			} else {
				log.Printf("%s %s: json: created %s\n", r.Method, r.URL, input.fname)
			}
		}

		if m == nil {
			log.Printf("%s %s: map is null\n", r.Method, r.URL)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if input.shiftX != 0 {
			m.ShiftX(-1 * m.Width() * input.shiftX / 100)
		}
		if input.shiftY != 0 {
			m.ShiftY(m.Height() * input.shiftY / 100)
		}

		// generate color map
		cm := cmap.FromHistogram(m.Histogram(), input.pctWater, input.pctIce, cmap.Water, cmap.Terrain, cmap.Ice)

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

// helper functions
func pfvAsInt(r *http.Request, key string) (int, error) {
	raw := r.PostFormValue(key)
	if raw == "" {
		return 0, fmt.Errorf("%q: missing", key)
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%q: %w", key, err)
	}
	return val, nil
}

func pfvAsString(r *http.Request, key string) (string, error) {
	raw := r.PostFormValue(key)
	if raw == "" {
		return "", fmt.Errorf("%q: missing", key)
	}
	return raw, nil
}

func pfvAsUint(r *http.Request, key string) (uint64, error) {
	raw := r.PostFormValue(key)
	if raw == "" {
		return 0, fmt.Errorf("%q: missing", key)
	}
	val, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("%q: %w", key, err)
	}
	return val, nil
}
