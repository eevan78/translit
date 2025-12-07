package terminal

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/exit"
)

var (
	Rdr             = bufio.NewReader(os.Stdin)
	Out             = bufio.NewWriter(os.Stdout)
	InputFilenames  []string
	InputFilePaths  []string
	OutputFilePaths []string
	OutputDir       = "output"
)

func OpenInputFile(filename string) (*os.File, *bufio.Reader) {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	Rdr = bufio.NewReader(inputFile)
	return inputFile, Rdr
}

func CreateOutputFile(filename string) (*os.File, *bufio.Writer) {
	outputFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	Out = bufio.NewWriter(outputFile)
	return outputFile, Out
}

func prepareInputDirectory() {
	inputDir, err := os.Open(*dictionary.InputPathPtr)
	if err != nil {
		panic(err)
	}

	InputFilenames, err = inputDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	for i := range InputFilenames {
		InputFilePaths = append(InputFilePaths, filepath.Join(absPath, InputFilenames[i]))
	}

	fmt.Println("Улазни фајлови:")
	printFilePaths(InputFilePaths)
}

func prepareOutputDirectory() {
	outDirName := filepath.Join(filepath.Dir(*dictionary.InputPathPtr), OutputDir)
	if _, err := os.Stat(outDirName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outDirName, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	absPath, _ := filepath.Abs(outDirName)
	for i := range InputFilenames {
		OutputFilePaths = append(OutputFilePaths, filepath.Join(absPath, InputFilenames[i]))
	}

	if len(OutputFilePaths) > 1 {
		fmt.Println("Излазни фајлови:")
	} else {
		fmt.Println("Излазни фајл:")
	}
	printFilePaths(OutputFilePaths)
}

func prepareInputFile() {
	var err error

	if strings.HasPrefix(*dictionary.InputPathPtr, "http") {
		if strings.HasSuffix(*dictionary.InputPathPtr, "/") {
			err = errors.New("тренутно није дозвољено да се URL завршава са /")
			exit.ExitWithError(err)
		}

		tmpDir := "tmp"

		if _, err := os.Stat(tmpDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(tmpDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		var response *grab.Response
		//download file to the tmp directory
		response, err = grab.Get(tmpDir, *dictionary.InputPathPtr)
		if err != nil {
			exit.ExitWithError(err)
		}
		*dictionary.InputPathPtr = response.Filename
	}

	// strip directories from the input filepath if exist

	InputFilenames = append(InputFilenames, filepath.Base(*dictionary.InputPathPtr))
	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	InputFilePaths = append(InputFilePaths, absPath)

	fmt.Println("Улазни фајл:")
	printFilePaths(InputFilePaths)
}

func isDirectory(path string) (bool, error) {
	if strings.HasPrefix(path, "http") {
		return false, nil
	}

	absPath, _ := filepath.Abs(path)

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func ProcessFlags() {
	flag.Usage = exit.Pomoc
	flag.Parse()
}

func CheckFlags() {
	if *dictionary.InputPathPtr != "" {
		// file no matter config
		if *dictionary.L2cPtr == *dictionary.C2lPtr || *dictionary.HtmlPtr || *dictionary.TextPtr {
			exit.ExitWithHelp()
		}
	} else {
		// std in
		arguments := os.Args[1:]
		if *dictionary.ConfigPtr {
			// config
			if len(arguments) == 1 {
				// program called only with -c flag so we test config
				if *dictionary.L2cPtr == *dictionary.C2lPtr || *dictionary.HtmlPtr == *dictionary.TextPtr {
					exit.ExitWithHelp()
				}
			} else {
				// program called with multiple flags
				exit.ExitWithHelp()
			}
		} else {
			// no config
			if *dictionary.L2cPtr == *dictionary.C2lPtr || *dictionary.HtmlPtr == *dictionary.TextPtr {
				exit.ExitWithHelp()
			}
		}
	}
}

func ProcessFilePaths() {
	if *dictionary.InputPathPtr != "" {
		isDirectory, errors := isDirectory(*dictionary.InputPathPtr)
		if errors != nil {
			panic(errors)
		}

		if isDirectory {
			prepareInputDirectory()
		} else {
			prepareInputFile()
		}
		prepareOutputDirectory()
	}
}

func printFilePaths(filePaths []string) {
	for i := range filePaths {
		fmt.Println(filePaths[i])
	}
}
