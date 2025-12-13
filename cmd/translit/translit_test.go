package main

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/terminal"
)

func TestReadingFromStdin(t *testing.T) {
	if os.Getenv("DO_TEST_X") == "1" {
		*dictionary.L2cPtr = true
		*dictionary.C2lPtr = false
		*dictionary.HtmlPtr = false
		*dictionary.TextPtr = true
		*dictionary.InputPathPtr = ""
		flag.Parse()
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReadingFromStdin")
	cmd.Env = append(os.Environ(), "DO_TEST_X=1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	input := "Pitamo se da li će uspeti?\nNadamo se da hoće…\n"
	output := "Питамо се да ли ће успети?\nНадамо се да хоће…\nPASS\n"

	func() {
		defer stdin.Close()
		fmt.Fprintf(os.Stderr, "Тест улаз филтера је:\n%s\n", input)
		io.WriteString(stdin, input)
	}()
	capture, _ := io.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Добијени излаз из филтера је:\n%s\n", string(capture))
	if output == string(capture) {
		fmt.Fprintln(os.Stderr, "Ухваћен је очекивани излаз!")
		clearData()
		return
	}

	clearData()
	t.Fatalf("Није успело читање са стандардног улаза: %v", err)

}

func TestL2CHtmlInputFileFromInternet(t *testing.T) {
	*dictionary.L2cPtr = true
	*dictionary.C2lPtr = false
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
	if os.Getenv("DO_TEST") == "1" {
		*dictionary.L2cPtr = true
		*dictionary.C2lPtr = false
		*dictionary.InputPathPtr = "https://zadovoljna.nova.rs/fitnes-i-ishrana/francuski-tost-przenice-iz-rerne/"
		flag.Parse()
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestL2CHtmlInputFileFromInternetWithTrailingSlash")
	cmd.Env = append(os.Environ(), "DO_TEST=1")
	out, err := cmd.CombinedOutput()
	fmt.Printf("Излаз: %s", out)
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		//Program exits with code 1 and prints expected message: Тренутно није дозвољено да се URL завршава са /
		clearData()
		return
	}

	clearData()
	t.Fatalf("Процес је бацио грешку %v, а требало је да статус изласка из програма буде 1", err)

}

func TestL2CTextInputFile(t *testing.T) {
	*dictionary.L2cPtr = true
	*dictionary.C2lPtr = false
	*dictionary.InputPathPtr = "../../test/testdata/rec_godine.txt"
	flag.Parse()

	expectedOutput, _ := filepath.Abs("../../test/testdata/rec_godine_izlaz.txt")

	main()

	exist := isOutputFileExist()
	clearData()

	if !exist {
		t.Fatalf(`Транслит није направио фајл %q`, getOutputFileName())
	} else {
		fmt.Fprintln(os.Stderr, "Пронађен!")
	}

	transliterated, err := os.Open(getOutputFileName())
	if err != nil {
		log.Fatal(err)
	}
	defer transliterated.Close()
	expected, err := os.Open(expectedOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer expected.Close()

	checksumTransliterate := checksumSHA256(transliterated)
	checksumExpected := checksumSHA256(expected)

	if strings.Compare(checksumTransliterate, checksumExpected) != 0 {
		t.Fatalf("Садржај пресловљеног фајла се разликује од очекиваног!")
	} else {
		fmt.Fprintln(os.Stderr, "Садржај пресловљеног фајла је исправан")
	}
}

func clearData() {
	terminal.InputFilenames = nil
	terminal.InputFilePaths = nil
	terminal.OutputFilePaths = nil
}

func checksumSHA256(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getOutputFileName() string {
	outDirPath := filepath.Join(filepath.Dir(*dictionary.InputPathPtr), terminal.OutputDir)
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
