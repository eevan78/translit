package language

import (
	"github.com/eevan78/translit/internal/archive"
	"github.com/eevan78/translit/internal/exit"
	"github.com/eevan78/translit/internal/terminal"
)

type ZipArchive struct {
	inputFilePath  string
	outputFilePath string
	innerDocuments []Document
	translitDir    string // where to place transliterated files
	unzipDir       string // where to place unzipped files
}

func (document *ZipArchive) open() {
	document.unzipDir, document.translitDir = terminal.PrepareZipDirectories(document.inputFilePath)
	archive.Unzip(document.inputFilePath, document.unzipDir)
	inputFilePaths := terminal.PrepareInputDirectoryForZip(document.unzipDir)
	outputFilePaths := terminal.PrepareOutputDirectoryForZip(document.unzipDir, inputFilePaths, document.translitDir)
	document.innerDocuments = CreateZipDocuments(inputFilePaths, outputFilePaths)
}

func (document *ZipArchive) transliterate() {
	Transliterate(document.innerDocuments)
}

func (document *ZipArchive) getInputFilePath() string {
	return document.inputFilePath
}

func (document *ZipArchive) getOuputFilePath() string {
	return document.outputFilePath
}

func (document *ZipArchive) finalize() {
	if err := archive.Zip(document.translitDir, document.outputFilePath); err != nil {
		exit.ExitWithError(err, document.inputFilePath)
	}
}
