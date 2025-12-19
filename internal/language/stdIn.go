package language

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/exit"
)

type StdIn struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func (document *StdIn) open() {
	document.reader = bufio.NewReader(os.Stdin)
	document.writer = bufio.NewWriter(os.Stdout)
}

func (document *StdIn) transliterate() {

loop:
	for {
		switch line, err := document.reader.ReadString('\n'); err {
		case nil:
			lineprefix := dictionary.Whitepref.FindString(line)
			words := strings.Fields(line)
			doit := true
			for n := range words {
				if strings.HasPrefix(words[n], "<|") {
					doit = false                                  // Do not transliterate
					words[n] = strings.TrimPrefix(words[n], "<|") // Remove marker of the beginning
					words[n] = fixPunctuation(words[n])
				}
				if strings.HasSuffix(words[n], "|>") {
					doit = true                                   // Transliterate after this word
					words[n] = strings.TrimSuffix(words[n], "|>") // Remove marker of the end
					words[n] = fixPunctuation(words[n])
					continue // Move to the next word
				}
				if !doit {
					words[n] = fixPunctuation(words[n])
					continue
				}
				if *dictionary.L2cPtr {
					index := transliterationIndexOfWordStartsWith(strings.ToLower(words[n]), dictionary.WholeForeignWords, "-")
					if index >= 0 {
						words[n] = string(words[n][:index]) + l2c(string(words[n][index:]))
					} else if !looksLikeForeignWord(words[n]) {
						words[n] = l2c(words[n])
					}
				} else if *dictionary.C2lPtr {
					words[n] = c2l(words[n])
				}
			}

			if lineprefix != "" && lineprefix != "\n" && len(words) != 0 {
				words[0] = lineprefix + words[0]
			}
			outl := strings.Join(words, " ")
			outl += "\n"
			if _, err = document.writer.WriteString(outl); err != nil {
				exit.ExitWithError(err)
			}
			_ = document.writer.Flush()

		case io.EOF:
			break loop

		default:
			exit.ExitWithError(err)
		}
	}
}

func (document *StdIn) getInputFilePath() string {
	return ""
}

func (document *StdIn) getOuputFilePath() string {
	return ""
}
