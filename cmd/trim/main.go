package main

import (
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
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

		// fmt.Printf("(%03d, %03d)    (%03d, %03d)\n", left, top, right, top)
		// fmt.Printf("(%03d, %03d)    (%03d, %03d)\n", left, bottom, right, bottom)

		r := image.Rectangle{image.Point{0, 0}, image.Point{right - left, bottom - top}}
		rgba := image.NewRGBA(r)

		draw.Draw(rgba, r.Bounds(), img, image.Point{left, top}, draw.Src)

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
