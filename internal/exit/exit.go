package exit

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/eevan78/translit/internal/dictionary"
)

func Pomoc() {
	fmt.Fprintf(flag.CommandLine.Output(), "Ово је филтер %s верзија %s\nСаставио eevan78, 2024-%v\n\n", os.Args[0], dictionary.ProgramVersion, (time.Now()).Year())
	fmt.Fprintf(flag.CommandLine.Output(), "Филтер чита UTF-8 кодирани текст са стандардног улаза или из наведеног фајла и исписује га на\nстандардни излаз или у излазни фајл, пресловљен сагласно са следећим заставицама:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nКада се наведе -c, не сме да се наведе ниједна друга заставица. Програм се подешава читањем\nконфигурације. У супротном, мора да се наведе по једна и само једна заставица из обе групе\nСмер и Формат. Када се наведе заставица за улазни фајл потребно је да се наведе само заставица смера.\nЦеле речи између „<|” и „|>” у простом тексту се не пресловљавају.\nТекст унутар <span lang=\"sr-Latn\"></span> елемента у (X)HTML се не пресловљава у ћирилицу,\nа текст унутар <span lang=\"sr-Cyrl\"></span> се не пресловљава у латиницу.\n\nПримери:\n%s -l2c -html\t\tпреслови (X)HTML у ћирилицу\n%s -text -c2l\t\tпреслови прости текст у латиницу\n%s -c\t\t\tпрограм чита подешавања из фајла конфигурације\n", os.Args[0], os.Args[0], os.Args[0])
}

func ExitWithError(err error, filename string) {
	fmt.Fprintln(os.Stderr, "Грешка у раду са: ", filename, err)
	os.Exit(1)
}

func ExitWithHelp() {
	Pomoc()
	os.Exit(1)
}
