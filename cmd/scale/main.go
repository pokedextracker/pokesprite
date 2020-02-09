package main

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"golang.org/x/image/draw"
)

const (
	threshold = 100
)

func main() {
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()
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

		if bounds.Dx() > threshold || bounds.Dy() > threshold {
			r := image.Rect(0, 0, bounds.Max.X/2, bounds.Max.Y/2)
			rgba := image.NewRGBA(r)
			draw.NearestNeighbor.Scale(rgba, r, img, bounds, draw.Over, nil)
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
}
