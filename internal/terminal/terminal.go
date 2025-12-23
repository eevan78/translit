package terminal

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/eevan78/translit/internal/exit"
)

var (
	L2cPtr       = flag.Bool("l2c", false, "`Смер` пресловљавања је латиница у ћирилицу")
	C2lPtr       = flag.Bool("c2l", false, "`Смер` пресловљавања је ћирилица у латиницу")
	HtmlPtr      = flag.Bool("html", false, "`Формат` улаза је (X)HTML")
	TextPtr      = flag.Bool("text", false, "`Формат` улаза је прости текст")
	ConfigPtr    = flag.Bool("c", false, "Користи се конфигурација")
	InputPathPtr = flag.String("i", "", "Путања улазног фајла или директоријума")

	OutputDir = "output"
)

func OpenInputFile(filename string) (*os.File, *bufio.Reader) {
	inputFile, err := os.Open(filename)
	if err != nil {
		exit.ExitWithError(err, filename)
	}

	rdr := bufio.NewReader(inputFile)
	return inputFile, rdr
}

func CreateOutputFile(filename string) (*os.File, *bufio.Writer) {
	outputFile, err := os.Create(filename)
	if err != nil {
		exit.ExitWithError(err, filename)
	}

	out := bufio.NewWriter(outputFile)
	return outputFile, out
}

func PrepareInputDirectory(inputDirectoryPath string) (inputFilePaths []string) {
	inputDir, err := os.Open(inputDirectoryPath)
	if err != nil {
		exit.ExitWithError(err, inputDirectoryPath)
	}

	inputFileNames, err := inputDir.Readdirnames(0)
	if err != nil {
		exit.ExitWithError(err, inputDirectoryPath)
	}

	absPath, _ := filepath.Abs(inputDirectoryPath)
	for i := range inputFileNames {
		inputFilePaths = append(inputFilePaths, filepath.Join(absPath, inputFileNames[i]))
	}

	return inputFilePaths
}

func PrepareOutputDirectory(inputDirectoryPath string, inputFilePaths []string, outputDirectoryPath string) (outputFilePaths []string) {
	if _, err := os.Stat(outputDirectoryPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(outputDirectoryPath, os.ModePerm)
		if err != nil {
			exit.ExitWithError(err, outputDirectoryPath)
		}
	}

	absPath, _ := filepath.Abs(outputDirectoryPath)
	var inputFileNames []string
	for i := range inputFilePaths {
		inputFileNames = append(inputFileNames, filepath.Base(inputFilePaths[i]))
		outputFilePaths = append(outputFilePaths, filepath.Join(absPath, inputFileNames[i]))
	}

	return outputFilePaths
}

func PrepareZipDirectories(inputFilePath string) (tempDir string, outputDir string) {
	// directory to place all archived files has the same name as the archive
	dirName := strings.Split(filepath.Base(inputFilePath), ".")[0]

	// Create a temporary directory with a custom prefix
	tempDir, err := os.MkdirTemp("", dirName)
	if err != nil {
		exit.ExitWithError(err, "Error creating temporary directory"+dirName)
	}

	outputDir, err = os.MkdirTemp("", "output")
	if err != nil {
		exit.ExitWithError(err, "Error creating temporary directory: output")
	}

	return tempDir, outputDir
}

func prepareInputFile() (inputFilePaths []string) {
	if strings.HasPrefix(*InputPathPtr, "http") {
		prepareInputFileFromInternet()
	}

	// strip directories from the input filepath if exist
	absPath, _ := filepath.Abs(*InputPathPtr)
	inputFilePaths = append(inputFilePaths, absPath)
	return inputFilePaths
}

func prepareInputFileFromInternet() {
	var err error
	if strings.HasSuffix(*InputPathPtr, "/") {
		err = errors.New("тренутно није дозвољено да се URL завршава са /")
		exit.ExitWithError(err, *InputPathPtr)
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
	response, err = grab.Get(tmpDir, *InputPathPtr)
	if err != nil {
		exit.ExitWithError(err, *InputPathPtr)
	}
	*InputPathPtr = response.Filename
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
	if *InputPathPtr != "" {
		// file no matter config
		if *L2cPtr == *C2lPtr || *HtmlPtr || *TextPtr {
			exit.ExitWithHelp()
		}
	} else {
		// std in
		arguments := os.Args[1:]
		if *ConfigPtr {
			// config
			if len(arguments) == 1 {
				// program called only with -c flag so we test config
				if *L2cPtr == *C2lPtr || *HtmlPtr == *TextPtr {
					exit.ExitWithHelp()
				}
			} else {
				// program called with multiple flags
				exit.ExitWithHelp()
			}
		} else {
			// no config
			if *L2cPtr == *C2lPtr || *HtmlPtr == *TextPtr {
				exit.ExitWithHelp()
			}
		}
	}
}

func ProcessFilePaths() (inputFilePaths []string, outputFilePaths []string) {
	if *InputPathPtr != "" {
		isDirectory, err := isDirectory(*InputPathPtr)
		if err != nil {
			exit.ExitWithError(err, *InputPathPtr)
		}
		if isDirectory {
			inputFilePaths = PrepareInputDirectory(*InputPathPtr)
		} else {
			inputFilePaths = prepareInputFile()
		}

		outputDirectoryPath := filepath.Join(filepath.Dir(*InputPathPtr), OutputDir)
		outputFilePaths = PrepareOutputDirectory(*InputPathPtr, inputFilePaths, outputDirectoryPath)
	}
	return inputFilePaths, outputFilePaths
}
