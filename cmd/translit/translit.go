// Line filter that transliterates UTF-8 coded plain text or (X)HTML between
// Serbian latin and Serbian cyrillic script. It properly handles foreign words,
// latin digraph splitting, units, and fixes punctuation. For plain text input
// it preserves the line indentation, but normalizes the rest of the whitespace
// in the line to one single space between each word.
package main

import (
	"github.com/eevan78/translit/internal/configuration"
	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/language"
	"github.com/eevan78/translit/internal/terminal"
)

func main() {

	configuration.ConfigInit()

	terminal.ProcessFlags()

	if *dictionary.HtmlPtr {
		language.TransliterateHtml()
	} else if *dictionary.TextPtr {
		language.TransliterateText()
	}

}
