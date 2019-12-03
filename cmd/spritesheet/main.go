// This script looks at the `images` directory and stitches all the images into
// a single spritesheet. If there are new images copied from
// https://github.com/msikma/pokesprite that are not in the `images` directory,
// you should run the `rename` script first. It will put the love ball sprite on
// the first row, and then put all Pokemon on subsequent rows. The `height` and
// `width` variables should be the max height and width possible. If a new
// generate produces larger sprites, you should update those values. While this
// script produces a PNG with the best compression, it's still recommended to
// use `pngcrush` to get it even smaller.
package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/pokedextracker/pokesprite/pkg/sorter"
)

const (
	height  = 30
	width   = 50
	columns = 32
)

func main() {
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		panic(err)
	}

	// Minus 1 to exclude the love ball, and plus 1 to include it again since it
	// will be on its own line.
	rows := int(math.Ceil(float64(len(files)-1)/columns)) + 1
	r := image.Rectangle{image.Point{0, 0}, image.Point{columns * width, rows * height}}
	rgba := image.NewRGBA(r)

	// Draw the love ball on its own row.
	err = drawImage(rgba, "love-ball.png", 0, 0)
	if err != nil {
		panic(err)
	}

	// Sort files alphabetically.
	sort.Sort(sorter.New(files))

	for i, file := range files {
		name := file.Name()
		// Skip drawing the love ball since we already drew it on its own row.
		if name == "love-ball.png" {
			continue
		}

		column := int(math.Mod(float64(i), float64(columns)))
		row := i/columns + 1

		err := drawImage(rgba, name, column, row)
		if err != nil {
			panic(err)
		}
	}

	// Write the output png to a file.
	out, err := os.Create("./output/pokesprite.png")
	if err != nil {
		panic(err)
	}
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = encoder.Encode(out, rgba)
	if err != nil {
		panic(err)
	}
}

func drawImage(rgba draw.Image, name string, column, row int) error {
	file, err := os.Open("./images/" + name)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	wdiff := width - img.Bounds().Dx()
	if wdiff < 0 {
		return fmt.Errorf("width (%dpx) is too small for %s (%dpx)", width, name, img.Bounds().Dx())
	}
	hdiff := height - img.Bounds().Dy()
	if hdiff < 0 {
		return fmt.Errorf("height (%dpx) is too small for %s (%dpx)", height, name, img.Bounds().Dy())
	}

	woffset := wdiff / 2
	hoffset := hdiff / 2
	x := column*width + woffset
	y := row*height + hoffset

	rect := image.Rectangle{image.Point{x, y}, image.Point{x + img.Bounds().Dx(), y + img.Bounds().Dy()}}
	draw.Draw(rgba, rect, img, image.Point{0, 0}, draw.Src)

	return nil
}
