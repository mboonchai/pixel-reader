package reader

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"path/filepath"
)

type gifProcessor struct {
}

func (p *gifProcessor) Ext() string {
	return ".gif"
}

func (p *gifProcessor) Process(fileName string, inDir string, outDir string) error {
	reader, err := os.Open(filepath.Join(inDir, fileName))
	if err != nil {
		return fmt.Errorf("error: file could not be opened: %w", err)
	}
	defer reader.Close()

	frames, err := SplitAnimatedGIF(reader)
	if err != nil {
		return fmt.Errorf("error: reading gif: %w", err)
	}

	err = WriteDartArray(fileName, outDir, frames)
	if err != nil {
		return fmt.Errorf("error: writing output: %w", err)
	}

	return nil
}

func SplitAnimatedGIF(reader io.Reader) (frames [][][]int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while decoding: %s", r)
		}
	}()

	gif, err := gif.DecodeAll(reader)

	if err != nil {
		return nil, err
	}

	frames = make([][][]int64, 0)

	width, height := getGifDimensions(gif)

	src_rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(src_rgba, src_rgba.Bounds(), gif.Image[0], image.ZP, draw.Src)

	//hack, read from only these pixels...
	readY := []int{43, 95, 147, 199, 251, 303, 355, 407}
	readX := []int{63, 115, 167, 219, 271, 323, 375, 427}
	//

	for curFrame, srcImg := range gif.Image {
		draw.Draw(src_rgba, src_rgba.Bounds(), srcImg, image.ZP, draw.Over)

		for iy := 0; iy < len(readY); iy++ {
			row := make([]int64, 0)
			for ix := 0; ix < len(readX); ix++ {
				idx_s := (readY[iy]*width + readX[ix]) * 4
				pix := src_rgba.Pix[idx_s : idx_s+4]

				color := int64(0xFF000000)
				//TRANSPARENT --> USE BLACK
				//NEAR BLACK --> USE BLACK
				if pix[3] != 0x00 && (pix[0] > 0x33 || pix[1] > 0x33 || pix[2] > 0x33) {
					color = (int64(pix[3]) << 24) + (int64(pix[0]) << 16) + (int64(pix[1]) << 8) + int64(pix[2])
				}
				row = append(row, color)
			}

			if len(row) > 0 {
				for len(frames) < curFrame+1 {
					frames = append(frames, make([][]int64, len(readY)))
				}
				frames[curFrame][iy] = row
			}
		}

	}

	return frames, nil
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
