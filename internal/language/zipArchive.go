package language

import (
	"fmt"

	"github.com/eevan78/translit/internal/archive"
	"github.com/eevan78/translit/internal/terminal"
)

type ZipArchive struct {
	inputFilePath  string
	outputFilePath string
	fop            *terminal.FileOperator
	inputFilePaths []string
}

func (document *ZipArchive) open() {
	tempDir := unzip(document.inputFilePath)
	filePaths := terminal.PrepareInputDirectory2(tempDir)

	fmt.Println(filePaths)

}

func (document *ZipArchive) transliterate() {

}

func (document *ZipArchive) getInputFilePath() string {
	return document.inputFilePath
}

func (document *ZipArchive) getOuputFilePath() string {
	return document.outputFilePath
}

func unzip(inputFilePath string) (tempDir string) {
	tempDir, _ = terminal.PrepareZipDirectories(inputFilePath)
	archive.Unzip(inputFilePath, tempDir)
	return tempDir
}
