// This script is used to rename icon files from
// https://github.com/msikma/pokesprite (using the `data/pkmn.json` file and the
// `icons` directory) into an `images` directory. This is the first step before
// generating a new spritesheet. It copies all Pokemon and forms (except
// right-facing ones and female ones) and their shiny versions as well. The
// naming format it uses is `<number>-<shiny>-<form>.png` where `<shiny>` and
// `<form>` are omitted if they are regular. For example, `001.png` and
// `026-shiny-alola.png`.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type pokemon struct {
	Index int `json:"idx"`
	Slug  struct {
		English string `json:"eng"`
	} `json:"slug"`
	Icons map[string]map[string]bool `json:"icons"`
}

func main() {
	data, err := ioutil.ReadFile("./data/pkmn.json")
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
		for form, opts := range pokemon.Icons {
			if opts["is_duplicate"] {
				continue
			}
			// Copy regular sprites.
			src := "./icons/pokemon/regular/" + pokemon.Slug.English
			dest := "./images/" + id
			if form != "." {
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
			src = "./icons/pokemon/shiny/" + pokemon.Slug.English
			dest = "./images/" + id + "-shiny"
			if form != "." {
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
	_, err = copy("./icons/pokeball/love.png", "./images/love-ball.png")
	if err != nil {
		panic(err)
	}
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
