package terminal

import (
	"bufio"
	"os"
)

type FileOperator struct {
	inputFile  *os.File
	outputFile *os.File
	Reader     *bufio.Reader
	Writer     *bufio.Writer
}

func (fop *FileOperator) Open(filePath string) {
	fop.inputFile, fop.Reader = OpenInputFile(filePath)
}

func (fop *FileOperator) Create(filePath string) {
	fop.outputFile, fop.Writer = CreateOutputFile(filePath)
}

func (fop *FileOperator) Close() {
	fop.inputFile.Close()
	fop.outputFile.Close()
}
