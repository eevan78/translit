package language

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/beevik/etree"
	"github.com/eevan78/translit/internal/dictionary"
	"github.com/eevan78/translit/internal/exit"
	"github.com/eevan78/translit/internal/terminal"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/net/html"
)

var (
	acceptedMime = map[string]string{
		"text":  "text/plain; charset=utf-8",
		"html":  "text/html; charset=utf-8",
		"xhtml": "application/xhtml+xml",
		"xml":   "text/xml; charset=utf-8",
	}
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
	lowercaseStr := strings.ToLower(trimExcessiveCharacters((str)))
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

func fixPunctuation(w string) string {
	w = dictionary.Pref.ReplaceAllStringFunc(w, dictionary.Prefmap.Replace)
	w = dictionary.Suff.ReplaceAllStringFunc(w, dictionary.Suffmap.Replace)
	return w
}

func l2c(s string) string {
	result := ""
	w := 0
	s = fixPunctuation(s)
	s = splitDigraphs(s)
	for i, runeValue := range s {
		if w > 1 {
			w -= 1
			continue
		}
		value, prefixLen, ok := dictionary.Tbl.SearchPrefixInString(s[i:])
		if ok {
			result += value
			w = utf8.RuneCountInString(s[i : i+prefixLen])
		} else {
			result += string(runeValue)
		}
	}
	// Remove ZWNJ from the transliterated word
	result = strings.ReplaceAll(result, "\u200C", "")
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
	s = fixPunctuation(s)

	for _, runeValue := range s {
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

func traverseHtmlNode(n *html.Node) {
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
		traverseHtmlNode(c)
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

// Traverses through the XML starting from the given xml element (node). Firstly, it transliterates text which can be mixed
// with other inner xml elements within this node. Then, it goes through the node and recursively do the traversal.
// CDATA section will be skipped and not transliterated.
func traverseXmlNode(node *etree.Element) {
	// iterates through element's Childs and transliterates only the text childs
	// these childs are any part of the text file including new line characters, inline text fields and xml elements
	for _, child := range node.Child {
		if childData, ok := child.(*etree.CharData); ok {
			// we ignore CDATA section
			if childData.IsCData() {
				continue
			}
			// if a child consists of a text transliterate it
			line := childData.Data
			if !allWhite(line) {
				childData.Data = transliterateXmlText(childData.Data)
			}
		}
	}
	// iterates through the Child elements which represent only xml elements
	for _, childElement := range node.ChildElements() {
		traverseXmlNode(childElement)
	}
}
func transliterateXmlText(line string) string {

	lineprefix := dictionary.Whitepref.FindString(line)
	linesuffix := dictionary.Whitesuff.FindString(line)
	words := strings.Fields(line)

	for word := range words {
		if *dictionary.L2cPtr {
			index := transliterationIndexOfWordStartsWith(strings.ToLower(words[word]), dictionary.WholeForeignWords, "-")
			if index >= 0 {
				words[word] = string(words[word][:index]) + l2c(string(words[word][index:]))
			} else if !looksLikeForeignWord(words[word]) {
				words[word] = l2c(words[word])
			}
		} else { // *c2lPtr
			words[word] = c2l(words[word])
		}
	}

	// Preserve the whitespace at the beginning and at the end of the line
	words[0] = lineprefix + words[0]
	words[len(words)-1] += linesuffix
	line = strings.Join(words, " ")
	return line
}

func Transliterate(documents []Document) []Document {
	if !isStdIn() {
		fmt.Println("Пресловљавање")
	}
	for i := range documents {
		documents[i].open()
		documents[i].transliterate()
		if !isStdIn() {
			fmt.Printf("Успешно: %s \nу %s\n", documents[i].getInputFilePath(), documents[i].getOuputFilePath())
		}
	}

	return documents
}

func CreateDocuments() []Document {
	documents := []Document{}

	if isStdIn() {
		documents = append(documents, &StdIn{})
	} else {
		for i := range terminal.InputFilenames {
			mediaType, _ := detectFileType(terminal.InputFilePaths[i])

			switch mediaType {
			case acceptedMime["text"]:
				documents = append(documents,
					&TextDocument{inputFilePath: terminal.InputFilePaths[i],
						outputFilePath: terminal.OutputFilePaths[i]})
			case acceptedMime["html"]:
				documents = append(documents,
					&HtmlDocument{inputFilePath: terminal.InputFilePaths[i],
						outputFilePath: terminal.OutputFilePaths[i]})
			case acceptedMime["xml"], acceptedMime["xhtml"]:
				documents = append(documents,
					&XmlDocument{inputFilePath: terminal.InputFilePaths[i],
						outputFilePath: terminal.OutputFilePaths[i]})
			default:
				fmt.Printf("Упозорење - тип фајла %s није подржан: %s\n", mediaType, terminal.InputFilePaths[i])

			}
		}
	}

	return documents
}

func detectFileType(filePath string) (string, string) {
	mimeType, err := mimetype.DetectFile(filePath)
	if err != nil {
		exit.ExitWithError(err, filePath)
	}

	// converting to lower case not to worry about the case of retrieved string value
	mediaType := strings.ToLower(mimeType.String())

	return mediaType, mimeType.Extension()
}

func isStdIn() bool {
	return *dictionary.InputPathPtr == ""
}
