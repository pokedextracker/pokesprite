// This script is used to rename icon files from
// https://github.com/msikma/pokesprite (using the `data/pokemon.json` file and
// the `pokemon-*` directories) into an `images` directory. This is the first
// step before generating a new spritesheet. It copies all Pokemon and forms
// (except right-facing ones and female ones) and their shiny versions as well.
// The naming format it uses is `<number>-<shiny>-<form>.png` where `<shiny>`
// and `<form>` are omitted if they are regular. For example, `001.png` and
// `026-shiny-alola.png`. It assumes that you have the repo cloned into
// ../msikma-pokesprite. To make things easier, you can delete all irrelevant
// Pokemon from the JSON file, but either way, once you've run this command, you
// probably want to just remove any git modifications. This is mostly useful for
// new sprites.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type pokemon struct {
	Slug struct {
		English string `json:"eng"`
	} `json:"slug"`
	Gen7 gen `json:"gen-7"`
	Gen8 gen `json:"gen-8"`
}

type gen struct {
	Forms map[string]struct {
		IsAliasOf *string `json:"is_alias_of"`
	} `json:"forms"`
}

func main() {
	data, err := ioutil.ReadFile("../msikma-pokesprite/data/pokemon.json")
	if err != nil {
		panic(err)
	}
	pkmn := map[string]*pokemon{}

	err = json.Unmarshal(data, &pkmn)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./images", 0744)
	if err != nil {
		panic(err)
	}

	for id, pokemon := range pkmn {
		for form, opts := range pokemon.Gen8.Forms {
			if opts.IsAliasOf != nil {
				continue
			}
			// Copy regular sprites.
			src := "../msikma-pokesprite/pokemon-gen8/regular/" + pokemon.Slug.English
			dest := "./images/" + id
			if form != "$" {
				src += "-" + form
				dest += "-" + form
			}
			src += ".png"
			dest += ".png"
			_, err := copy(src, dest)
			if err != nil {
				panic(err)
			}
			// Copy shiny sprites.
			src = "../msikma-pokesprite/pokemon-gen8/shiny/" + pokemon.Slug.English
			dest = "./images/" + id + "-shiny"
			if form != "$" {
				src += "-" + form
				dest += "-" + form
			}
			src += ".png"
			dest += ".png"
			_, err = copy(src, dest)
			if err != nil {
				panic(err)
			}
		}
	}

	// Move love ball sprite.
	// _, err = copy("./icons/pokeball/love.png", "./images/love-ball.png")
	// if err != nil {
	// 	panic(err)
	// }
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
