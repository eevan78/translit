package language

import (
	"fmt"

	"github.com/eevan78/translit/internal/exit"
	"github.com/eevan78/translit/internal/terminal"
	"golang.org/x/net/html"
)

type HtmlDocument struct {
	inputFilePath  string
	outputFilePath string
	fop            *terminal.FileOperator
}

func (document *HtmlDocument) open() {
	document.fop = &terminal.FileOperator{}
	document.fop.Open(document.inputFilePath)
	document.fop.Create(document.outputFilePath)
}

func (document *HtmlDocument) transliterate() {
	node, err := html.Parse(document.fop.Reader)
	if err != nil {
		exit.ExitWithError(err, document.getInputFilePath())
	}
	traverseHtmlNode(node)
	if err := html.Render(document.fop.Writer, node); err != nil {
		exit.ExitWithError(err, document.getInputFilePath())
	}
	_ = document.fop.Writer.Flush()
}

func (document *HtmlDocument) getInputFilePath() string {
	return document.inputFilePath
}

func (document *HtmlDocument) getOuputFilePath() string {
	return document.outputFilePath
}

func (document *HtmlDocument) finalize() {
	fmt.Printf("Успешно: %s \nу %s\n", document.inputFilePath, document.outputFilePath)
}
