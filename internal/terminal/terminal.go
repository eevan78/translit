package terminal

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/exit"
)

func pomoc() {
	fmt.Fprintf(flag.CommandLine.Output(), "Ово је филтер %s верзија %s\nСаставио eevan78, 2024-%v\n\n", os.Args[0], dictionary.ProgramVersion, (time.Now()).Year())
	fmt.Fprintf(flag.CommandLine.Output(), "Филтер чита UTF-8 кодирани текст са стандардног улаза или из наведеног фајла и исписује га на\nстандардни излаз или наведени фајл, пресловљен сагласно са следећим заставицама:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nКада се наведе -c, не сме да се наведе ниједна друга заставица. Програм се подешава читањем\nконфигурације. У супротном, мора да се наведе по једна и само једна заставица из обе групе\nСмер и Формат. Заставице за путању улазног и излазног фајла, као и излазног директоријума\nнису обавезне. Целе речи између „<|” и „|>” у простом тексту се не пресловљавају у ћирилицу.\nТекст унутар <span lang=\"sr-Latn\"></span> елемента у (X)HTML се не пресловљава у ћирилицу,\nа текст унутар <span lang=\"sr-Cyrl\"></span> се не пресловљава у латиницу.\n\nПримери:\n%s -l2c -html\t\tпреслови (X)HTML у ћирилицу\n%s -text -c2l\t\tпреслови прости текст у латиницу\n%s -c\t\t\tпрограм чита подешавања из фајла конфигурације\n", os.Args[0], os.Args[0], os.Args[0])
}

func OpenInputFile(filename string) {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	dictionary.Rdr = bufio.NewReader(inputFile)
}

func prepareInputDirectory() {
	inputDir, err := os.Open(*dictionary.InputPathPtr)
	if err != nil {
		panic(err)
	}

	dictionary.InputFilenames, err = inputDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	for i := range dictionary.InputFilenames {
		dictionary.InputFilePaths = append(dictionary.InputFilePaths, filepath.Join(absPath, dictionary.InputFilenames[i]))
	}

	fmt.Println("Улазни фајлови:")
	fmt.Println(dictionary.InputFilePaths)
}

func prepareOutputDirectory() {
	outDirName := filepath.Join(filepath.Dir(*dictionary.InputPathPtr), dictionary.OutputDir)
	if _, err := os.Stat(outDirName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outDirName, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	absPath, _ := filepath.Abs(outDirName)
	for i := range dictionary.InputFilenames {
		dictionary.OutputFilePaths = append(dictionary.OutputFilePaths, filepath.Join(absPath, dictionary.InputFilenames[i]))
	}

	fmt.Println("Излазни фајлови:")
	fmt.Println(dictionary.OutputFilePaths)
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

	dictionary.InputFilenames = append(dictionary.InputFilenames, filepath.Base(*dictionary.InputPathPtr))
	absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
	dictionary.InputFilePaths = append(dictionary.InputFilePaths, absPath)

	fmt.Println("Улазни фајл:")
	fmt.Println(dictionary.InputFilePaths)
}

func CreateOutputFile(filename string) {
	outputFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	dictionary.Out = bufio.NewWriter(outputFile)
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

func arrayContainsSubstring(array []string, value string) bool {
	for _, element := range array {
		if strings.Contains(element, value) {
			return true
		}
	}
	return false
}

// Reset opposite flag for the one added as a command line argument.
func resetOppositeFlags() {
	arguments := os.Args[1:]

	if arrayContainsSubstring(arguments, "l2c") {
		*dictionary.C2lPtr = false
	}
	if arrayContainsSubstring(arguments, "c2l") {
		*dictionary.L2cPtr = false
	}
	if arrayContainsSubstring(arguments, "html") {
		*dictionary.TextPtr = false
	}
	if arrayContainsSubstring(arguments, "text") {
		*dictionary.HtmlPtr = false
	}
}

func ProcessFlags() {
	flag.Usage = pomoc
	resetOppositeFlags()
	flag.Parse()
	if !*dictionary.ConfigPtr && (*dictionary.L2cPtr == *dictionary.C2lPtr || *dictionary.HtmlPtr == *dictionary.TextPtr) ||
		*dictionary.ConfigPtr && (*dictionary.L2cPtr || *dictionary.C2lPtr || *dictionary.HtmlPtr || *dictionary.TextPtr) {
		pomoc()
		os.Exit(1)
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
