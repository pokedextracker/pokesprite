package size

import (
	"image"
	"os"

	"github.com/pkg/errors"
)

func Max(files []os.FileInfo) (int, int, error) {
	height := 0
	width := 0

	for _, file := range files {
		name := file.Name()

		file, err := os.Open("./images/" + name)
		if err != nil {
			return 0, 0, err
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			return 0, 0, errors.Wrapf(err, "image decode: %s: ", name)
		}
		bounds := img.Bounds()

		if bounds.Dy() > height {
			height = bounds.Dy()
		}
		if bounds.Dx() > width {
			width = bounds.Dx()
		}
	}

	return height, width, nil
}
