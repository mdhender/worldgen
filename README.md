# World Map Generator

A map generator.

## Options
* Seed:
  This the "seed" used to initialize the random number generator.
  If you use the same seed, you'll get the same output.
  This is helpful if you want to customize a map.
* Percent Water:
  This number controls the "sea level" in the map.
  Higher values result in more oceans/seas and less land.
  Weird values are either rejected or ignored.
* Percent Ice:
  Higher values result in more ice covered lands.
  Ice starts at the highest elevations.
  Lower elevations are covered as this value increases.
* Shift X:
  Moves the generated image left by the percentage amount.
  Helpful if you want to center a certain area of the images.
* Shift Y:
  Moves the generated image up by the percentage amount.
  Helpful if you want to center a certain area of the images.
* Secret:
  For combating spam.
  If you're running on a server, then you have to type in the secret to create new maps.
  (If you're running locally, this is field is not displayed and secrets are not required.)

# Building
1. Clone the repository.
2. Build using `go build ./...`.

# Running

## Mac or Linux
Start the `cmd/wg/wg` executable.

## Windows
Start the `cmd\wg\wg.exe` executable..

# Viewing
Open `http://localhost:8080/` in your browser.

## First time
Note that it takes about ten seconds to generate a new image for any seed:

    2023/06/14 17:32:18 POST /generate: json: created c0ffeecafe.json
    2023/06/14 17:32:18 POST /generate: elapsed 6.342958125sn

The results are cached so that viewing or customizing for the same seed value happens in a fraction of a second:

    2023/06/14 17:32:39 POST /generate: entering
    2023/06/14 17:32:39 POST /generate: elapsed 251.442791msn
