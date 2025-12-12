package language

import (
	"github.com/beevik/etree"
	"github.com/eevan78/translit/internal/terminal"
)

type XmlDocument struct {
	inputFilePath  string
	outputFilePath string
	fop            *terminal.FileOperator
}

func (document *XmlDocument) open() {
	document.fop = &terminal.FileOperator{}
	document.fop.Open(document.inputFilePath)
	document.fop.Create(document.outputFilePath)
}

func (document *XmlDocument) transliterate() {
	xmlDocument := etree.NewDocument()
	if _, err := xmlDocument.ReadFrom(document.fop.Reader); err != nil {
		panic(err)
	}
	traverseXmlNode(&xmlDocument.Element)
	xmlDocument.WriteTo(document.fop.Writer)

	_ = document.fop.Writer.Flush()
}
