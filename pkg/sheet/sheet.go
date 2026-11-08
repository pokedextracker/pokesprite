// Package sheet decides which spritesheet each image in the `images` directory
// goes on and where it sits. The spritesheet and css scripts both use it, so
// the positions in the stylesheet always match the spritesheets.
package sheet

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	// Dir is where the source images live.
	Dir = "./images"
	// LoveBall is the one image that isn't a Pokémon. It's served as its own
	// file instead of being on a spritesheet, so the home page can show it
	// without downloading a whole spritesheet.
	LoveBall = "love-ball.png"

	columns = 32
	// gutter is the empty space between icons, so a neighbor's edge can't show
	// when the browser scales the page by a fraction (zoom, or a device pixel
	// ratio like 1.25).
	gutter = 2
	prefix = "pokesprite"
)

// separateFamilies are the game families whose icons get spritesheets of
// their own, so pages for other games never download them. Their icons are a
// different style (Legends: Arceus icons are drawn in a circle), so they're
// only ever shown in dexes for that game. A family here must also be in
// nameRE, or its images would be read as a form.
var separateFamilies = map[string]bool{
	"legends_arceus": true,
}

var nameRE = regexp.MustCompile(`^(\d+)(-shiny)?(-legends_arceus|-home)?(-.*)?\.png$`)

// Icon is one image and where it sits on its spritesheet.
type Icon struct {
	File     string
	Selector string
	X        int
	Y        int
	Width    int
	Height   int
}

// Sheet is one spritesheet and the icons on it.
type Sheet struct {
	// Name is the file name without an extension, e.g. pokesprite-shiny.
	Name   string
	Icons  []Icon
	Width  int
	Height int
}

// Load reads every image in the `images` directory except the love ball and
// lays them out on spritesheets: one for regular icons and one for shiny
// icons, plus a regular and a shiny one for each separate game family.
func Load() ([]*Sheet, error) {
	files, err := os.ReadDir(Dir)
	if err != nil {
		return nil, err
	}

	sheets := map[string]*Sheet{}
	for _, file := range files {
		name := file.Name()
		if name == LoveBall || strings.HasPrefix(name, ".") {
			continue
		}

		matches := nameRE.FindStringSubmatch(name)
		if matches == nil {
			return nil, fmt.Errorf("%s: name doesn't match NNN[-shiny][-game_family][-form].png", name)
		}
		id := matches[1]
		shiny := matches[2] == "-shiny"
		gameFamily := strings.TrimPrefix(matches[3], "-")
		form := strings.TrimPrefix(matches[4], "-")

		selector := ".pkicon.pkicon-" + id
		if form != "" {
			selector += ".form-" + form
		}
		if gameFamily != "" {
			selector += ".game-family-" + gameFamily
		}
		if shiny {
			selector += ".color-shiny"
		}

		sheetName := prefix
		if separateFamilies[gameFamily] {
			sheetName += "-" + strings.ReplaceAll(gameFamily, "_", "-")
		}
		if shiny {
			sheetName += "-shiny"
		}

		width, height, err := Size(name)
		if err != nil {
			return nil, err
		}

		s := sheets[sheetName]
		if s == nil {
			s = &Sheet{Name: sheetName}
			sheets[sheetName] = s
		}
		s.Icons = append(s.Icons, Icon{File: name, Selector: selector, Width: width, Height: height})
	}

	var out []*Sheet
	for _, s := range sheets {
		s.layout()
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// layout places the icons in a grid, in file name order, with every cell the
// size of the sheet's largest icon plus the gutter.
func (s *Sheet) layout() {
	sort.Slice(s.Icons, func(i, j int) bool { return s.Icons[i].File < s.Icons[j].File })

	cellWidth, cellHeight := 0, 0
	for _, icon := range s.Icons {
		if icon.Width > cellWidth {
			cellWidth = icon.Width
		}
		if icon.Height > cellHeight {
			cellHeight = icon.Height
		}
	}
	cellWidth += gutter
	cellHeight += gutter

	for i := range s.Icons {
		s.Icons[i].X = i % columns * cellWidth
		s.Icons[i].Y = i / columns * cellHeight
	}

	cols := columns
	if len(s.Icons) < cols {
		cols = len(s.Icons)
	}
	rows := (len(s.Icons) + columns - 1) / columns
	// The last column and row don't need a gutter after them.
	s.Width = cols*cellWidth - gutter
	s.Height = rows*cellHeight - gutter
}

// Size returns the width and height of an image in the `images` directory.
func Size(name string) (int, int, error) {
	file, err := os.Open(filepath.Join(Dir, name))
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", name, err)
	}
	return config.Width, config.Height, nil
}
