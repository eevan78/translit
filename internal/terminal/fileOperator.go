package terminal

import (
	"bufio"
	"os"
)

type FileOperator struct {
	InputFile  *os.File
	OutputFile *os.File
	Reader     *bufio.Reader
	Writer     *bufio.Writer
}

func (fop *FileOperator) Open(filePath string) {
	fop.InputFile, fop.Reader = OpenInputFile(filePath)
}

func (fop *FileOperator) Create(filePath string) {
	fop.OutputFile, fop.Writer = CreateOutputFile(filePath)
}

func (fop *FileOperator) Close() {
	fop.InputFile.Close()
	fop.OutputFile.Close()
}
