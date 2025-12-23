package language

import (
	"fmt"
	"os"

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
	inputFilePaths := terminal.PrepareInputDirectory(document.unzipDir)
	outputFilePaths := terminal.PrepareOutputDirectory(document.unzipDir, inputFilePaths, document.translitDir)
	document.innerDocuments = CreateDocuments(inputFilePaths, outputFilePaths)
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
	inputDir := document.translitDir
	transliteratedFiles, _ := os.ReadDir(inputDir)

	if len(transliteratedFiles) > 0 {
		if err := archive.Zip(inputDir, document.outputFilePath); err != nil {
			exit.ExitWithError(err, document.inputFilePath)
		}
		fmt.Printf("Успешно: %s \nу %s\n", document.inputFilePath, document.outputFilePath)
	} else {
		fmt.Println("Неуспешно: ниједан фајл у улазној zip архиви није успешно пресловљен.", document.inputFilePath)
	}
}
