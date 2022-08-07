package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"golang.org/x/image/draw"
)

var (
	threshold    int
	scale        float64
	height       int
	width        int
	scalerString string

	scalers = map[string]draw.Scaler{
		"pixel":   draw.NearestNeighbor,
		"regular": draw.CatmullRom,
	}
)

func main() {
	// Flags
	flag.IntVar(&threshold, "threshold", 100, "Threshold (in px) for whether to scale the image")
	flag.Float64Var(&scale, "scale", 0.5, "Factor to scale by (should be used with threshold)")
	flag.IntVar(&height, "height", -1, "Explicit height (in px) to scale to (use with width; takes precedence over scale)")
	flag.IntVar(&width, "width", -1, "Explicit width (in px) to scale to (use with height; takes precedence over scale)")
	flag.StringVar(&scalerString, "scaler", "pixel", "The scaler to use e.g. pixel, regular")
	flag.Parse()

	if height == -1 && width != -1 ||
		width == -1 && height != -1 {
		fmt.Println("Error: Both height and width need to be specified together")
		fmt.Println("Example: task scale -- -height 100 -width 100")
		os.Exit(1)
	}

	scaler, ok := scalers[scalerString]
	if !ok {
		fmt.Printf("Error: Invalid scaler %q\n", scalerString)
		fmt.Println("Example: task scale -- -scaler regular")
		os.Exit(1)
	}

	// Logic
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
			var r image.Rectangle
			if height == -1 && width == -1 {
				r = image.Rect(0, 0, int(math.Round(float64(bounds.Max.X)*scale)), int(math.Round(float64(bounds.Max.Y)*scale)))
			} else {
				r = image.Rect(0, 0, width, height)
			}
			rgba := image.NewRGBA(r)
			scaler.Scale(rgba, r, img, bounds, draw.Over, nil)
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
