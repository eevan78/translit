package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestL2CHtmlInputFileFromInternet(t *testing.T) {
	*l2cPtr = true
	*c2lPtr = false
	*htmlPtr = true
	*textPtr = false
	*inputPathPtr = "https://www.k1info.rs/kultura-i-umetnost/knjige/33243/bestseler-debitantski-roman-melisa-da-kosta/vest"
	flag.Parse()

	main()

	exist := isOutputFileExist()
	clearData()
	if !exist {
		t.Fatalf(`Translit nije napravio fajl %q`, getOutputFileName())
	}
}

func TestL2CHtmlInputFileFromInternetWithTrailingSlash(t *testing.T) {
	*l2cPtr = true
	*c2lPtr = false
	*htmlPtr = true
	*textPtr = false
	*inputPathPtr = "https://zadovoljna.nova.rs/fitnes-i-ishrana/francuski-tost-przenice-iz-rerne/"
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
	*l2cPtr = true
	*c2lPtr = false
	*htmlPtr = false
	*textPtr = true
	*inputPathPtr = "test/rec_godine.txt"
	flag.Parse()

	main()

	exist := isOutputFileExist()
	clearData()

	if !exist {
		t.Fatalf(`Translit nije napravio fajl %q`, getOutputFileName())
	}
}

func clearData() {
	inputFilenames = nil
	inputFilePaths = nil
	outputFilePaths = nil
}

func getOutputFileName() string {

	absDirectoryPath, _ := filepath.Abs(outputDir)
	absFilenamePath := filepath.Join(absDirectoryPath, filepath.Base(*inputPathPtr))

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
