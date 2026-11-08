// This script looks at the `images` directory and stitches the images into
// spritesheets: one for regular icons, one for shiny icons, and a regular and
// shiny one for each game family whose icons are kept separate (see
// pkg/sheet). That way a page only downloads the icons it can show. If there
// are new images copied from https://github.com/msikma/pokesprite that are not
// in the `images` directory, you should run the `rename` script first.
//
// Each spritesheet is written as a lossless WebP using `cwebp`, which
// compresses much better than PNG. The love ball isn't on a spritesheet. It's
// copied to the output directory as is, so the home page can show it without
// downloading a spritesheet.
package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pokedextracker/pokesprite/pkg/sheet"
)

const outputDir = "./output"

func main() {
	sheets, err := sheet.Load()
	if err != nil {
		panic(err)
	}

	// Remove spritesheets from earlier runs, so a sheet that no longer exists
	// doesn't get copied.
	old, err := filepath.Glob(filepath.Join(outputDir, "*.webp"))
	if err != nil {
		panic(err)
	}
	for _, path := range old {
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}

	for _, s := range sheets {
		if err := write(s); err != nil {
			panic(err)
		}
	}

	loveBall, err := os.ReadFile(filepath.Join(sheet.Dir, sheet.LoveBall))
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(outputDir, sheet.LoveBall), loveBall, 0644)
	if err != nil {
		panic(err)
	}
}

// write draws a spritesheet and saves it as a lossless WebP. It's drawn to a
// temporary PNG first since that's what `cwebp` reads.
func write(s *sheet.Sheet) error {
	rgba := image.NewNRGBA(image.Rect(0, 0, s.Width, s.Height))
	for _, icon := range s.Icons {
		if err := drawImage(rgba, icon); err != nil {
			return err
		}
	}

	pngPath := filepath.Join(outputDir, s.Name+".png")
	out, err := os.Create(pngPath)
	if err != nil {
		return err
	}
	defer os.Remove(pngPath)
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	err = encoder.Encode(out, rgba)
	out.Close()
	if err != nil {
		return err
	}

	// -lossless with -z 9 is the slowest, smallest lossless setting. -exact
	// is left off on purpose: without it, cwebp may change the color of fully
	// transparent pixels, which can't be seen, and the file is smaller.
	cmd := exec.Command("cwebp", "-quiet", "-mt", "-lossless", "-z", "9", pngPath, "-o", filepath.Join(outputDir, s.Name+".webp"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func drawImage(rgba draw.Image, icon sheet.Icon) error {
	file, err := os.Open(filepath.Join(sheet.Dir, icon.File))
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rect := image.Rect(icon.X, icon.Y, icon.X+icon.Width, icon.Y+icon.Height)
	draw.Draw(rgba, rect, img, img.Bounds().Min, draw.Src)
	return nil
}
