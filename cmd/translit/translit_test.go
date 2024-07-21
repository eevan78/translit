package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/eevan78/translit/internal/dictionary"
)

func TestL2CHtmlInputFileFromInternet(t *testing.T) {
	*dictionary.L2cPtr = true
	*dictionary.C2lPtr = false
	*dictionary.HtmlPtr = true
	*dictionary.TextPtr = false
	*dictionary.InputPathPtr = "https://www.k1info.rs/kultura-i-umetnost/knjige/33243/bestseler-debitantski-roman-melisa-da-kosta/vest"
	flag.Parse()

	main()

	exist := isOutputFileExist()
	clearData()
	if !exist {
		t.Fatalf(`Translit nije napravio fajl %q`, getOutputFileName())
	}
}

func TestL2CHtmlInputFileFromInternetWithTrailingSlash(t *testing.T) {
	*dictionary.L2cPtr = true
	*dictionary.C2lPtr = false
	*dictionary.HtmlPtr = true
	*dictionary.TextPtr = false
	*dictionary.InputPathPtr = "https://zadovoljna.nova.rs/fitnes-i-ishrana/francuski-tost-przenice-iz-rerne/"
	flag.Parse()

	if os.Getenv("DO_TEST") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestL2CHtmlInputFileFromInternetWithTrailingSlash")
	cmd.Env = append(os.Environ(), "DO_TEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		//Program exits with code 1 and prints expected message: Тренутно није дозвољено да се URL завршава са /
		clearData()
		return
	}

	clearData()
	t.Fatalf("Процес је бацио грешку %v, а требало је да статус изласка из пробрама буде 1", err)

}

func TestL2CTextInputFile(t *testing.T) {
	*dictionary.L2cPtr = true
	*dictionary.C2lPtr = false
	*dictionary.HtmlPtr = false
	*dictionary.TextPtr = true
	*dictionary.InputPathPtr = "../../test/testdata/rec_godine.txt"
	flag.Parse()

	main()

	exist := isOutputFileExist()
	clearData()

	if !exist {
		t.Fatalf(`Translit nije napravio fajl %q`, getOutputFileName())
	}
}

func clearData() {
	dictionary.InputFilenames = nil
	dictionary.InputFilePaths = nil
	dictionary.OutputFilePaths = nil
}

func getOutputFileName() string {
	outDirPath := filepath.Join(filepath.Dir(*dictionary.InputPathPtr), dictionary.OutputDir)
	absDirectoryPath, _ := filepath.Abs(outDirPath)
	absFilenamePath := filepath.Join(absDirectoryPath, filepath.Base(*dictionary.InputPathPtr))

	return absFilenamePath
}

func isExist(filePath string) bool {

	fmt.Println("Тражим да ли постоји фајл:")
	fmt.Println(filePath)

	_, err := os.Stat(filePath)
	return errors.Is(err, os.ErrNotExist)
}

func isOutputFileExist() bool {
	fileName := getOutputFileName()
	return !isExist(fileName)
}
