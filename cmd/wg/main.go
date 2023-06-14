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

// Package main implements a web server for the map generator.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mdhender/worldgen/pkg/gen"
	"github.com/mdhender/worldgen/pkg/generator"
	"github.com/mdhender/worldgen/pkg/sliced"
	"github.com/mdhender/worldgen/pkg/smite"
	"github.com/mdhender/worldgen/pkg/tiled"
	"github.com/mdhender/worldgen/pkg/way"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	Seed := 0x638bb317ac47a6ba
	rand.Seed(int64(Seed))

	height, width, iterations := 600, 1_200, 10_000

	router := way.NewRouter()

	router.Handle("GET", "/", &templateHandler{filename: "index.gohtml"})
	router.HandleFunc("GET", "/fracture", nextSeedHandler("fracture"))
	router.HandleFunc("GET", "/fracture/:seed", fractureHandler(height, width, iterations))
	router.HandleFunc("GET", "/smite", nextSeedHandler("smite"))
	router.HandleFunc("GET", "/smite/:seed", smiteHandler(height, width, iterations))
	router.HandleFunc("GET", "/tile", nextSeedHandler("tile"))
	router.HandleFunc("GET", "/tile/:seed", tileHandler(height, width, iterations))
	router.HandleFunc("GET", "/tiled", nextSeedHandler("tiled"))
	router.HandleFunc("GET", "/tiled/:seed", tiledHandler(height, width, iterations))
	router.HandleFunc("GET", "/asteroids", nextSeedHandler("asteroids"))
	router.HandleFunc("GET", "/asteroids/:seed", asteroidsHandler(height, width, iterations))
	router.HandleFunc("GET", "/customize/:seed", customizeHandler("../templates", "customize.gohtml"))
	router.HandleFunc("GET", "/carto/:seed", cartoHandler())
	router.HandleFunc("GET", "/greyscale/:seed", greyscaleHandler())

	log.Fatalln(http.ListenAndServe(":8080", router))
}

func asteroidsHandler(height, width, iterations int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		seed, err := strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		m := gen.New(height, width, rand.New(rand.NewSource(int64(seed))))
		m.RandomFractureCircle(iterations)
		m.Normalize()

		png, err := m.AsPNG(m.AsImage())
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("WG-Seed", fmt.Sprintf("%x", seed))
		w.WriteHeader(http.StatusOK)
		w.Write(png)

		log.Printf("asteroidsHandler: %x elapsed %v\n", seed, time.Now().Sub(started))

		data, err := json.Marshal(m)
		if err != nil {
			log.Printf("json: %v\n", err)
		} else if err = os.WriteFile(fmt.Sprintf("%x-asteroids.json", seed), data, 0644); err != nil {
			log.Printf("json: %v\n", err)
		} else {
			log.Printf("json: created %x-asteroids.json\n", seed)
		}
	}
}

func fractureHandler(height, width, iterations int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img, err := sliced.Generate(height, width, iterations)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		bb := &bytes.Buffer{}
		err = png.Encode(bb, img)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Write(bb.Bytes())
	}
}

// 4192195a3a17473f

func nextSeedHandler(kind string) http.HandlerFunc {
	//rnd := rand.New(rand.NewSource(int64(seed)))
	//var lock sync.Mutex
	//nextSeed := func() uint64 {
	//	lock.Lock()
	//	defer lock.Unlock()
	//	return rnd.Uint64()
	//}
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/%s/%x", kind, rand.Uint64()), http.StatusSeeOther)
	}
}

func smiteHandler(height, width, iterations int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img, err := smite.Generate(height, width, iterations)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		bb := &bytes.Buffer{}
		err = png.Encode(bb, img)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Write(bb.Bytes())
	}
}

func tileHandler(height, width, iterations int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		seed, err := strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		img, err := tiled.Generate(height, width, iterations, rand.New(rand.NewSource(int64(seed))))
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		bb := &bytes.Buffer{}
		err = png.Encode(bb, img)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		//log.Printf("tileHandler: seed %q %x\n", pSeed, seed)

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("WG-Seed", fmt.Sprintf("%x", seed))
		w.WriteHeader(http.StatusOK)
		w.Write(bb.Bytes())

		log.Printf("tileHander: %x elapsed %v\n", seed, time.Now().Sub(started))
	}
}

func tiledHandler(height, width, iterations int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		pSeed := way.Param(r.Context(), "seed")
		if pSeed == "" {
			http.Error(w, "missing seed", http.StatusBadRequest)
			return
		}
		seed, err := strconv.ParseUint(pSeed, 16, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
			return
		}

		m := generator.New(height, width, rand.New(rand.NewSource(int64(seed))))
		m.RandomFractureCircle(iterations)
		m.Normalize()
		png, err := m.AsPNG()
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		//log.Printf("tiledHandler: seed %q %x\n", pSeed, seed)

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("WG-Seed", fmt.Sprintf("%x", seed))
		w.WriteHeader(http.StatusOK)
		w.Write(png)

		log.Printf("tiledHander: %x elapsed %v\n", seed, time.Now().Sub(started))
	}
}
