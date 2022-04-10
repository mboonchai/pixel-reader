package reader

type gifProcessor struct {
}

func (p *gifProcessor) Ext() string {
	return ".gif"
}

func (p *gifProcessor) Process(fileName string, inDir string, outDir string) error {
	return nil
}
