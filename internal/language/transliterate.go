package language

import (
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/exit"
	"github.com/eevan78/translit/internal/terminal"
	"golang.org/x/net/html"
)

func looksLikeForeignWord(word string) bool {
	trimmedWord := trimExcessiveCharacters(word)
	processed := strings.ToLower(trimmedWord)
	if processed == "" {
		return false
	}

	if wordStartsWith(processed, dictionary.SerbianWordsWithForeignCharacterCombinations) {
		return false
	}

	if wordContainsString(processed, dictionary.ForeignCharacterCombinations) {
		return true
	}

	if wordStartsWith(processed, dictionary.CommonForeignWords) {
		return true
	}

	if wordIsEqualTo(processed, dictionary.WholeForeignWords) {
		return true
	}

	if wordContainsMeasurementUnit(trimmedWord) {
		return true
	}

	return false
}

func wordStartsWith(word string, array []string) bool {
	for _, arrayWord := range array {
		if strings.HasPrefix(word, arrayWord) {
			return true
		}
	}
	return false
}

func wordContainsString(word string, array []string) bool {
	for _, arrayWord := range array {
		if strings.Contains(word, arrayWord) {
			return true
		}
	}
	return false
}

func wordIsEqualTo(word string, array []string) bool {
	for _, arrayWord := range array {
		if word == arrayWord {
			return true
		}
	}
	return false
}

func transliterationIndexOfWordStartsWith(word string, array []string, charSeparator string) int {
	processed := strings.ToLower(trimExcessiveCharacters(word))
	if processed == "" {
		return -1
	}
	for _, arrayWord := range array {
		if strings.HasPrefix(word, arrayWord+charSeparator) {
			return len(arrayWord + charSeparator)
		}
	}

	return -1
}

func trimExcessiveCharacters(word string) string {
	const excessiveChars = "[\\s!?,:;.\\*\\-—~`'\"„”“”‘’(){}\\[\\]<>«»\\/\\\\]"
	regExp := regexp.MustCompile("^(" + excessiveChars + ")+|(" + excessiveChars + ")+$")

	return regExp.ReplaceAllString(word, "")
}

func wordContainsMeasurementUnit(word string) bool {
	const unitAdjacentToSth = `([zafpnμmcdhKMGTPEY]?([BVWJFSHCΩATNhlmg]|m[²³]?|s[²]?|cd|Pa|Wb|Hz))`
	const unitOptionalyAdjacentToSth = `(°[FC]|[kMGTPZY](B|Hz)|[pnμmcdhk]m[²³]?|m[²³]|[mcdh][lg]|kg|km)`
	const number = `(((\d+([\.,]\d)*)|(\d*[½⅓¼⅕⅙⅐⅛⅑⅒⅖¾⅗⅜⅘⅚⅝⅞])))`
	regExp := regexp.MustCompile("^(" + number + unitAdjacentToSth + ")|(" + number + "?(" + unitOptionalyAdjacentToSth + "|" + unitAdjacentToSth + "/" + unitAdjacentToSth + "))$")

	return regExp.MatchString(word)
}

func splitDigraphs(str string) string {
	lowercaseStr := strings.ToLower(str)
	strout := strings.Clone(str)
	for digraph := range dictionary.DigraphExceptions {
		if !strings.Contains(lowercaseStr, digraph) {
			continue
		}
		for _, word := range dictionary.DigraphExceptions[digraph] {
			if !strings.HasPrefix(lowercaseStr, word) {
				continue
			}
			// Split all possible occurrences, regardless of case.
			for key, word := range dictionary.DigraphReplacements[digraph] {
				strout = strings.Replace(strout, key, word, 1)
			}
			break
		}
	}
	return strout
}

func l2c(s string) string {
	result := ""
	w := 0
	stopafter := false
	word := strings.Clone(s)
	if strings.HasSuffix(word, "|>") {
		stopafter = true
		word = strings.TrimSuffix(word, "|>")
	}
	// Fix punctuation
	word = dictionary.Pref.ReplaceAllStringFunc(word, dictionary.Prefmap.Replace)
	word = dictionary.Suff.ReplaceAllStringFunc(word, dictionary.Suffmap.Replace)

	word = splitDigraphs(word)
	if dictionary.Doit {
		for i, runeValue := range word {
			if w > 1 {
				w -= 1
				continue
			}
			value, prefixLen, ok := dictionary.Tbl.SearchPrefixInString(word[i:])
			if ok {
				result += string(value)
				w = utf8.RuneCountInString(string(word[i : i+prefixLen]))
			} else {
				result += string(runeValue)
			}
		}
	} else {
		result = word
	}
	// Remove ZWNJ from the transliterated word
	result = strings.ReplaceAll(result, "\u200C", "")
	if stopafter {
		dictionary.Doit = true
	}
	return result
}

func Uppercase(s string) string {
	up := ""
	for _, runeValue := range s {
		up += string(unicode.ToUpper(runeValue))
	}
	return up
}

func c2l(s string) string {
	result := ""
	word := strings.Clone(s)
	// Fix punctuation
	word = dictionary.Pref.ReplaceAllStringFunc(word, dictionary.Prefmap.Replace)
	word = dictionary.Suff.ReplaceAllStringFunc(word, dictionary.Suffmap.Replace)

	for _, runeValue := range word {
		value, ok := dictionary.Tbl1[string(runeValue)]
		if ok {
			result += value
		} else {
			result += string(runeValue)
		}
	}
	if dictionary.Fixdigraphs.MatchString(result) {
		return Uppercase(result)
	} else {
		return result
	}
}

func allWhite(s string) bool {
	result := true
	for _, runevalue := range s {
		if !unicode.IsSpace(runevalue) {
			result = false
			break
		}
	}
	return result
}

func traverseNode(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		// Properly adjust the lang attribute, or add it if it's missing
		if n.Data == "html" {
			namespace := ""
			notexist := true
			if *dictionary.L2cPtr {
				for i, attrib := range n.Attr {
					if attrib.Key == "lang" || attrib.Key == "xml:lang" {
						n.Attr[i].Val = "sr-Cyrl-t-sr-Latn"
						notexist = false
					}
					if attrib.Key == "xml:lang" || attrib.Key == "xmlns" {
						namespace = "xml"
					}
				}
				if notexist {
					n.Attr = append(n.Attr, html.Attribute{Namespace: namespace, Key: "lang", Val: "sr-Cyrl-t-sr-Latn"})
				}
			} else if *dictionary.C2lPtr {
				for i, attrib := range n.Attr {
					if attrib.Key == "lang" || attrib.Key == "xml:lang" {
						n.Attr[i].Val = "sr-Latn-t-sr-Cyrl"
						notexist = false
					}
					if attrib.Key == "xml:lang" || attrib.Key == "xmlns" {
						namespace = "xml"
					}
				}
				if notexist {
					n.Attr = append(n.Attr, html.Attribute{Namespace: namespace, Key: "lang", Val: "sr-Latn-t-sr-Cyrl"})
				}
			}
		}
	case html.TextNode:
		// Transliterate if text is not inside a script or a style element
		if !allWhite(n.Data) && n.Parent.Type == html.ElementNode && shouldTransliterate(n) {
			nodeprefix := dictionary.Whitepref.FindString(n.Data)
			nodesuffix := dictionary.Whitesuff.FindString(n.Data)
			words := strings.Fields(n.Data)

			for w := range words {
				if *dictionary.L2cPtr {
					index := transliterationIndexOfWordStartsWith(strings.ToLower(words[w]), dictionary.WholeForeignWords, "-")
					if index >= 0 {
						words[w] = string(words[w][:index]) + l2c(string(words[w][index:]))
					} else if !looksLikeForeignWord(words[w]) {
						words[w] = l2c(words[w])
					}
				} else { // *c2lPtr
					words[w] = c2l(words[w])
				}
			}

			// Preserve the whitespace at the beginning and at the end of the node data
			words[0] = nodeprefix + words[0]
			words[len(words)-1] += nodesuffix
			n.Data = strings.Join(words, " ")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c)
	}
}

// Checks whether a text node should be transliterated. Returns true if it should, and false otherwise.
// A node should not be transliterated if its parent node is script or style.
// Also, node should not be transliterated if it has a lang attribute with a value of Latin script
// when transliteration is to the Cyrillic script, or if it has a lang attribute with a value of Cyrillic script
// when transliteration is to the Latin script.
func shouldTransliterate(n *html.Node) bool {

	if n.Parent.Data == "script" || n.Parent.Data == "style" {
		return false
	}

	attr := ""

	if *dictionary.L2cPtr {
		attr = "sr-Latn"
	} else { // *c2lPtr
		attr = "sr-Cyrl"
	}

	shouldTranslit := true
	if n.Parent.Data == "span" {
		for _, k := range n.Parent.Attr {
			if (k.Key == "lang" || k.Key == "xml:lang") && k.Val == attr {
				shouldTranslit = false
				break
			}
		}

	}

	return shouldTranslit
}

func transliterateHtmlFile() {
	doc, err := html.Parse(dictionary.Rdr)
	if err != nil {
		exit.ExitWithError(err)
	}
	traverseNode(doc)
	if err := html.Render(dictionary.Out, doc); err != nil {
		exit.ExitWithError(err)
	}
	_ = dictionary.Out.Flush()
}

func transliterateTextFile() {

loop:
	for {
		switch line, err := dictionary.Rdr.ReadString('\n'); err {
		case nil:
			lineprefix := dictionary.Whitepref.FindString(line)
			words := strings.Fields(line)
			for n := range words {
				if strings.HasPrefix(words[n], "<|") {
					dictionary.Doit = false
					words[n] = strings.TrimPrefix(words[n], "<|")
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
			if _, err = dictionary.Out.WriteString(outl); err != nil {
				exit.ExitWithError(err)
			}
			_ = dictionary.Out.Flush()

		case io.EOF:
			break loop

		default:
			exit.ExitWithError(err)
		}
	}
}

func TransliterateHtml() {
	for i := range dictionary.InputFilenames {
		terminal.OpenInputFile(dictionary.InputFilePaths[i])
		terminal.CreateOutputFile(dictionary.OutputFilePaths[i])
		transliterateHtmlFile()
	}
}

func TransliterateText() {
	for i := range dictionary.InputFilenames {
		terminal.OpenInputFile(dictionary.InputFilePaths[i])
		terminal.CreateOutputFile(dictionary.OutputFilePaths[i])
		transliterateTextFile()
	}
}
