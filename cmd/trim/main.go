package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// This is used for debugging to only attempt to process a single file. This should be the full file name without
	// the directory e.g. "999-shiny.png".
	debugFile := os.Getenv("DEBUG_FILE")

	files, err := ioutil.ReadDir("./images")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()
		if debugFile != "" && name != debugFile {
			// We have a debug file set and this file isn't it, so we skip it.
			continue
		}
		if strings.Contains(name, "legends_arceus") {
			// We intentionally want whitespace for the Legends Arceus sprites.
			debug("%s: skipping because legends arceus sprites need whitespace", name)
			continue
		}

		// Open the file so we can get its dimensions.
		file, err := os.Open("./images/" + name)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			panic(err)
		}
		bounds := img.Bounds()
		top := 0
		bottom := 0
		left := 0
		right := 0

	top:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				transparent := a == 0
				if !transparent {
					top = y
					break top
				}
			}
		}
	bottom:
		for y := bounds.Max.Y; y >= bounds.Min.Y; y-- {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				transparent := a == 0
				if !transparent {
					bottom = y + 1
					break bottom
				}
			}
		}
	left:
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				_, _, _, a := img.At(x, y).RGBA()
				transparent := a == 0
				if !transparent {
					left = x
					break left
				}
			}
		}
	right:
		for x := bounds.Max.X; x >= bounds.Min.X; x-- {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				_, _, _, a := img.At(x, y).RGBA()
				transparent := a == 0
				if !transparent {
					right = x + 1
					break right
				}
			}
		}

		isCorrectSize := top == bounds.Min.Y && bottom == bounds.Max.Y && left == bounds.Min.X && right == bounds.Max.X
		debug("%s: bounds.Min.Y: %d, bounds.Max.Y: %d, bounds.Min.X: %d, bounds.Max.X: %d", name, bounds.Min.Y, bounds.Max.Y, bounds.Min.X, bounds.Max.X)
		debug("%s: top: %d, bottom: %d, left: %d, right: %d", name, top, bottom, left, right)
		debug("%s: isCorrectSize: %v", name, isCorrectSize)

		if isCorrectSize {
			// This file doesn't need to be cropped at all, so we just move onto the next one. If we don't do this, we
			// often see a lot of images thrashing when we rewrite it even when it theoretically shouldn't have changed.
			debug("%s: skipping because it's already the right size", name)
			continue
		}

		// This rectangle represents the full size of the final image. Since we just have the coordinates, we need to
		// calculate width and height by subtracting.
		r := image.Rectangle{image.Point{0, 0}, image.Point{right - left, bottom - top}}
		rgba := image.NewRGBA(r)

		// This point that we're passing in is the offset point to start the cropping.
		draw.Draw(rgba, r.Bounds(), img, image.Point{left, top}, draw.Src)

		// Open the file again so we can write to it.
		out, err := os.Create("./images/" + name)
		if err != nil {
			panic(err)
		}
		encoder := png.Encoder{}
		err = encoder.Encode(out, rgba)
		out.Close()
		if err != nil {
			panic(err)
		}
	}
}

func debug(format string, a ...any) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf(format+"\n", a...)
	}
}
