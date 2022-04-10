package reader

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

type pngProcessor struct {
}

//color that use to separat frames
const sepColor = 0xFF00FE00

func (p *pngProcessor) Ext() string {
	return ".png"
}

func (p *pngProcessor) Process(fileName string, inDir string, outDir string) error {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	reader, err := os.Open(filepath.Join(inDir, fileName))
	if err != nil {
		return fmt.Errorf("error: file could not be opened: %w", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return fmt.Errorf("error: decoding image: %w", err)
	}

	frames := image_2_array_pix_frame(m)

	err = WriteDartArray(fileName, outDir, frames)
	if err != nil {
		return fmt.Errorf("error: writing output: %w", err)
	}

	return nil

}

func image_2_array_pix_frame(src image.Image) [][][]int64 {

	frames := make([][][]int64, 0)

	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	src_rgba := image.NewRGBA(src.Bounds())
	draw.Copy(src_rgba, image.Point{}, src, src.Bounds(), draw.Src, nil)

	curFrame := 0
	for y := 0; y < height; y++ {
		curFrame = 0
		row := make([]int64, 0)
		for x := 0; x < width; x++ {
			idx_s := (y*width + x) * 4
			pix := src_rgba.Pix[idx_s : idx_s+4]

			color := int64(0xFF000000)
			//TRANSPARENT --> USE BLACK
			if pix[3] != 0x00 {
				color = (int64(pix[3]) << 24) + (int64(pix[0]) << 16) + (int64(pix[1]) << 8) + int64(pix[2])
			}

			if color == sepColor {
				for len(frames) < curFrame+1 {
					frames = append(frames, make([][]int64, height))
				}

				frames[curFrame][y] = row
				row = make([]int64, 0)
				curFrame++
				continue
			}

			row = append(row, color)

		}

		if len(row) > 0 {
			for len(frames) < curFrame+1 {
				frames = append(frames, make([][]int64, height))
			}

			frames[curFrame][y] = row
		}
	}

	return frames
}
