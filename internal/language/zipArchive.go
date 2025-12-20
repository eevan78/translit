package language

import (
	"github.com/eevan78/translit/internal/archive"
	"github.com/eevan78/translit/internal/terminal"
)

type ZipArchive struct {
	inputFilePath  string
	outputFilePath string
	documents      []Document
}

func (document *ZipArchive) open() {
	unzipDir, translitDir := terminal.PrepareZipDirectories(document.inputFilePath)
	archive.Unzip(document.inputFilePath, unzipDir)
	inputFilePaths := terminal.PrepareInputDirectoryForZip(unzipDir)
	outputFilePaths := terminal.PrepareOutputDirectoryForZip(unzipDir, inputFilePaths, translitDir)
	document.documents = CreateZipDocuments(inputFilePaths, outputFilePaths)
}

func (document *ZipArchive) transliterate() {
	Transliterate(document.documents)
}

func (document *ZipArchive) getInputFilePath() string {
	return document.inputFilePath
}

func (document *ZipArchive) getOuputFilePath() string {
	return document.outputFilePath
}
