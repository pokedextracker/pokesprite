// This script is used to take an existing spritesheet and chop it up into
// individual sprites. It reads from the filename passed as the first argument
// to get all the necessary information about this spritesheet. This is a JSON
// file that we created, and it's structure is defined by the `Data` type below.
// This was first created to handle the combine spritesheet for Legends: Arceus
// found at https://www.spriters-resource.com/fullview/168439/ (regular) and
// https://www.spriters-resource.com/fullview/168440/ (shiny), but it shouldn't
// be exclusive to those.
package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
)

type Data struct {
	Filename string    `json:"filename"`
	Columns  int       `json:"columns"`
	Rows     int       `json:"rows"`
	Outline  int       `json:"outline_px_size"`
	Padding  int       `json:"padding_px_size"`
	Suffix   *string   `json:"suffix"`
	Pokemon  []Pokemon `json:"pokemon"`
}

type Pokemon struct {
	ID   int     `json:"id"`
	Form *string `json:"form"`
	Skip bool    `json:"skip"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: task chop -- <filename.json>")
		fmt.Println("Example: task chop -- ./data/spritesheet.json")
		os.Exit(1)
	}

	filename := os.Args[1]

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data := Data{}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		panic(err)
	}

	spritesheet, err := os.Open(data.Filename)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(spritesheet)
	spritesheet.Close()
	if err != nil {
		panic(err)
	}

	// height = (total height of spritesheet - ((rows + 1) * outline size)) / rows
	height := (img.Bounds().Size().Y - ((data.Rows + 1) * data.Outline)) / data.Rows
	width := (img.Bounds().Size().X - ((data.Columns + 1) * data.Outline)) / data.Columns

	for index, pokemon := range data.Pokemon {
		if pokemon.Skip {
			continue
		}

		// Calculate which row and column we're on based on the index.
		row := index / data.Columns
		column := index % data.Columns

		// Create new image data.
		r := image.Rectangle{image.Point{0, 0}, image.Point{width - 2*data.Padding, height - 2*data.Padding}}
		rgba := image.NewRGBA(r)
		draw.Draw(rgba, r.Bounds(), img, image.Point{column*height + (column+1)*data.Outline + data.Padding, row*width + (row+1)*data.Outline + data.Padding}, draw.Src)

		// Generate the new filename.
		filename := fmt.Sprintf("./images/%03d", pokemon.ID)
		if data.Suffix != nil {
			filename += "-" + *data.Suffix
		}
		if pokemon.Form != nil {
			filename += "-" + *pokemon.Form
		}
		filename += ".png"

		// Write the new chopped up png out.
		out, err := os.Create(filename)
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
