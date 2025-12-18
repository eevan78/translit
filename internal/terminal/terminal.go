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
	"github.com/gabriel-vasile/mimetype"
)

var (
	rdr             = bufio.NewReader(os.Stdin)
	out             = bufio.NewWriter(os.Stdout)
	InputFilenames  []string
	InputFilePaths  []string
	OutputFilePaths []string
	OutputDir       = "output"
	TmpDir          = "tmp"
)

func OpenInputFile(filename string) (*os.File, *bufio.Reader) {
	inputFile, err := os.Open(filename)
	if err != nil {
		exit.ExitWithError(err, filename)
	}

	rdr = bufio.NewReader(inputFile)
	return inputFile, rdr
}

func CreateOutputFile(filename string) (*os.File, *bufio.Writer) {
	outputFile, err := os.Create(filename)
	if err != nil {
		exit.ExitWithError(err, filename)
	}

	out = bufio.NewWriter(outputFile)
	return outputFile, out
}

func prepareInputDirectory() {
	inputDir, err := os.Open(*dictionary.InputPathPtr)
	if err != nil {
		exit.ExitWithError(err, *dictionary.InputPathPtr)
	}

	InputFilenames, err = inputDir.Readdirnames(0)
	if err != nil {
		exit.ExitWithError(err, *dictionary.InputPathPtr)
	}

	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	for i := range InputFilenames {
		InputFilePaths = append(InputFilePaths, filepath.Join(absPath, InputFilenames[i]))
	}
}

func PrepareInputDirectory2(directoryPath string) (filePaths []string) {
	inputDir, err := os.Open(directoryPath)
	if err != nil {
		exit.ExitWithError(err, directoryPath)
	}

	fileNames, err := inputDir.Readdirnames(0)
	if err != nil {
		exit.ExitWithError(err, directoryPath)
	}

	absPath, _ := filepath.Abs(directoryPath)
	for i := range fileNames {
		filePaths = append(filePaths, filepath.Join(absPath, fileNames[i]))
	}

	return filePaths
}

func prepareOutputDirectory() {
	outDirName := filepath.Join(filepath.Dir(*dictionary.InputPathPtr), OutputDir)
	if _, err := os.Stat(outDirName); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(outDirName, os.ModePerm)
		if err != nil {
			exit.ExitWithError(err, outDirName)
		}
	}

	absPath, _ := filepath.Abs(outDirName)
	for i := range InputFilenames {
		OutputFilePaths = append(OutputFilePaths, filepath.Join(absPath, InputFilenames[i]))
	}
}

func PrepareZipDirectories(inputFilePath string) (tempDir string, outputDir string) {
	// directory to place all archived files has the same name as the archive
	dirName := strings.Split(filepath.Base(inputFilePath), ".")[0]

	// Create a temporary directory with a custom prefix
	tempDir, err := os.MkdirTemp("", dirName)
	if err != nil {
		exit.ExitWithError(err, "Error creating temporary directory")
	}

	outputDir = OutputDir + string(os.PathSeparator) + dirName
	outputDir, _ = filepath.Abs(outputDir)

	return tempDir, outputDir
}

func prepareInputFile() {
	if strings.HasPrefix(*dictionary.InputPathPtr, "http") {
		prepareInputFileFromInternet()
	}

	// strip directories from the input filepath if exist
	InputFilenames = append(InputFilenames, filepath.Base(*dictionary.InputPathPtr))
	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	InputFilePaths = append(InputFilePaths, absPath)
}

func prepareInputFileFromInternet() {
	var err error
	if strings.HasSuffix(*dictionary.InputPathPtr, "/") {
		err = errors.New("тренутно није дозвољено да се URL завршава са /")
		exit.ExitWithError(err, *dictionary.InputPathPtr)
	}

	tmpDir := "tmp"

	if _, err := os.Stat(tmpDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tmpDir, os.ModePerm)
		if err != nil {
			exit.ExitWithError(err, tmpDir)
		}
	}

	var response *grab.Response
	//download file to the tmp directory
	response, err = grab.Get(tmpDir, *dictionary.InputPathPtr)
	if err != nil {
		exit.ExitWithError(err, *dictionary.InputPathPtr)
	}
	*dictionary.InputPathPtr = response.Filename
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
		isDirectory, err := isDirectory(*dictionary.InputPathPtr)
		if err != nil {
			exit.ExitWithError(err, *dictionary.InputPathPtr)
		}

		if isDirectory {
			prepareInputDirectory()
		} else {
			prepareInputFile()
		}

		prepareOutputDirectory()
	}
}

func DetectFileType(filePath string) (string, string) {
	mediaType, err := mimetype.DetectFile(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(mediaType.String(), mediaType.Extension())

	return mediaType.String(), mediaType.Extension()
}
