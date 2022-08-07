// This script looks at the `images` directory and generates a .scss file with
// the classes and the correct `background-position`s so that the spritesheet
// can be used for PokedexTracker.
package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/pokedextracker/pokesprite/pkg/size"
	"github.com/pokedextracker/pokesprite/pkg/sorter"
)

const (
	columns = 32

	preamble = `.pkicon {
  @include crisp-rendering();

  display: inline-block;
  background-image: url('/pokesprite.png');
  background-repeat: no-repeat;

  &.game-family-legends_arceus {
    border-radius: 50%;
  }
}

.pkicon.pkicon-ball-love { %s }

`
)

var (
	height = 0
	width  = 0
)

var nameRE = regexp.MustCompile(`(\d+)(-shiny)?(-legends_arceus)?(-.*)?\.png`)

func main() {
	var buf bytes.Buffer

	files, err := ioutil.ReadDir("./images")
	if err != nil {
		panic(err)
	}

	height, width, err = size.Max(files)
	if err != nil {
		panic(err)
	}

	loveBallStyles, err := generateStyles("love-ball.png", 0, 0)
	if err != nil {
		panic(err)
	}
	_, err = buf.WriteString(fmt.Sprintf(preamble, loveBallStyles))
	if err != nil {
		panic(err)
	}

	// Sort files alphabetically.
	sort.Sort(sorter.New(files))

	for i, file := range files {
		name := file.Name()
		// Skip love ball styles since it's already been written.
		if name == "love-ball.png" {
			continue
		}

		matches := nameRE.FindAllStringSubmatch(name, -1)
		id := matches[0][1]
		shiny := matches[0][2] == "-shiny"
		gameFamily := strings.Trim(matches[0][3], "-")
		form := strings.Trim(matches[0][4], "-")

		class := ".pkicon.pkicon-" + id

		if form != "" {
			class += ".form-" + form
		}

		if gameFamily != "" {
			class += ".game-family-" + gameFamily
		}

		if shiny {
			class += ".color-shiny"
		}

		column := int(math.Mod(float64(i), float64(columns)))
		row := i/columns + 1

		styles, err := generateStyles(name, column, row)
		if err != nil {
			panic(err)
		}

		_, err = buf.WriteString(fmt.Sprintf("%s { %s }\n", class, styles))
		if err != nil {
			panic(err)
		}
	}

	err = ioutil.WriteFile("./output/pokesprite.scss", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func generateStyles(name string, column, row int) (string, error) {
	file, err := os.Open("./images/" + name)
	if err != nil {
		return "", err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	bounds := img.Bounds()
	x := column * width
	y := row * height

	// Due to the way background-position works, these values need to be
	// negative.
	styles := fmt.Sprintf(
		"width: %dpx; height: %dpx; background-position: %dpx %dpx;",
		bounds.Dx(),
		bounds.Dy(),
		x*-1,
		y*-1,
	)

	return styles, nil
}
