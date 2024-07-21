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

func Pomoc() {
	fmt.Fprintf(flag.CommandLine.Output(), "Ово је филтер %s верзија %s\nСаставио eevan78, 2024\n\n", os.Args[0], dictionary.Version)
	fmt.Fprintf(flag.CommandLine.Output(), "Филтер чита UTF-8 кодирани текст са стандардног улаза или из наведеног фајла и исписује га на\nстандардни излаз или наведени фајл, пресловљен сагласно са следећим заставицама:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nМора да се наведе по једна и само једна заставица из обе групе Смер и Формат.\nЗаставице за путању улазног и излазног фајла, као и излазног директоријума нису у обавезне.\nЦеле речи између „<|” и „|>” у простом тексту се не пресловљавају у ћирилицу.\nТекст унутар <span lang=\"sr-Latn\"></span> елемента у (X)HTML се не пресловљава у\nћирилицу, а текст унутар <span lang=\"sr-Cyrl\"></span> се не пресловљава у латиницу.\n\nПримери:\n%s -l2c -html\t\tпреслови (X)HTML у ћирилицу\n%s -text -c2l\t\tпреслови прости текст у латиницу\n", os.Args[0], os.Args[0])
}

func OpenInputFile(filename string) {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	dictionary.Rdr = bufio.NewReader(inputFile)
}

func prepareInputDirectory() {
	isDirectory, errors := isDirectory(*dictionary.InputPathPtr)
	if errors != nil {
		panic(errors)
	}

	if isDirectory {
		inputDir, err := os.Open(*dictionary.InputPathPtr)
		if err != nil {
			panic(err)
		}

		var error error

		dictionary.InputFilenames, error = inputDir.Readdirnames(0)
		if error != nil {
			panic(error)
		}

		absPath, _ := filepath.Abs(*dictionary.InputPathPtr)
		for i := range dictionary.InputFilenames {
			dictionary.InputFilePaths = append(dictionary.InputFilePaths, filepath.Join(absPath, dictionary.InputFilenames[i]))
		}

		fmt.Println("Улазни фајлови:")
		fmt.Println(dictionary.InputFilePaths)
	}
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

func ProcessFlags() {
	flag.Usage = Pomoc
	flag.Parse()
	if *dictionary.L2cPtr == *dictionary.C2lPtr || *dictionary.HtmlPtr == *dictionary.TextPtr {
		Pomoc()
		os.Exit(0)
	}

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
