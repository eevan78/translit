package terminal

import (
	"bufio"
	"fmt"
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
	err := fop.InputFile.Close()
	if err != nil {
		// could be ommited beacuse it is safe for us to call close twice, but printing for developers
		fmt.Println("Фајл је већ затворен:", fop.InputFile)
	}

	err = fop.OutputFile.Close()
	if err != nil {
		// could be ommited beacuse it is safe for us to call close twice, but printing for developers
		fmt.Println("Фајл је већ затворен:", fop.OutputFile)
	}
}
