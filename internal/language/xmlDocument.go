package language

import (
	"github.com/beevik/etree"
	"github.com/eevan78/translit/internal/exit"
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
	// do not consider CDATA section as XML element so we can differentiate them during transliteration.
	xmlDocument.ReadSettings = etree.ReadSettings{PreserveCData: true}
	if _, err := xmlDocument.ReadFrom(document.fop.Reader); err != nil {
		exit.ExitWithError(err, document.getInputFilePath())
	}
	traverseXmlNode(&xmlDocument.Element)
	xmlDocument.WriteTo(document.fop.Writer)

	_ = document.fop.Writer.Flush()
}

func (document *XmlDocument) getInputFilePath() string {
	return document.inputFilePath
}

func (document *XmlDocument) getOuputFilePath() string {
	return document.outputFilePath
}
