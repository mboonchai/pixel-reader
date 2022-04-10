package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Reader interface {
	Process(fileName string) error
}
type processor interface {
	Ext() string
	Process(fileName string, inDir string, outDir string) error
}

func New(inDir string, outDir string) Reader {
	return &reader{
		inDir:  inDir,
		outDir: outDir,
		processors: []processor{
			&gifProcessor{},
			&pngProcessor{},
		},
	}

}

type reader struct {
	inDir  string
	outDir string

	processors []processor
}

func (r *reader) Process(fileName string) error {

	ext := filepath.Ext(fileName)

	for _, p := range r.processors {
		if strings.ToLower(ext) == p.Ext() {
			err := p.Process(fileName, r.inDir, r.outDir)
			if err != nil {
				return fmt.Errorf("eror processing %s : %w", fileName, err)
			}

			return nil
		}
	}

	return fmt.Errorf("no valid processor to process file: %s", fileName)
}

func WriteDartArray(fileName string, outDir string, frames [][][]int64) error {

	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	outName := filepath.Join(outDir, fmt.Sprintf("%s.dart", name))

	writer, err := os.Create(outName)
	if err != nil {
		return fmt.Errorf("error: file could not be opened: %w", err)
	}
	defer writer.Close()

	writer.WriteString(fmt.Sprintf("final List<List<List<int>>> %s =[\n", name))

	frame := make([]string, 0)
	for _, f := range frames {
		col := make([]string, 0)
		for _, y := range f {
			row := make([]string, 0)
			for _, x := range y {
				row = append(row, fmt.Sprintf("0x%X", x))
			}

			col = append(col, fmt.Sprintf("  [\n    %s\n  ]\n", strings.Join(row, ",")))
		}

		frame = append(frame, fmt.Sprintf("[\n%s\n]\n", strings.Join(col, ",")))
	}

	writer.WriteString(strings.Join(frame, ","))
	writer.WriteString("];\n")

	return nil
}
