// Line filter that transliterates UTF-8 coded plain text or (X)HTML between
// Serbian latin and Serbian cyrillic script. It properly handles foreign words,
// latin digraph splitting, units, and fixes punctuation. For plain text input
// it preserves the line indentation, but normalizes the rest of the whitespace
// in the line to one single space between each word.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/porfirion/trie"
	"golang.org/x/net/html"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	whitepref   = regexp.MustCompile(`^[\s\p{Zs}]+`)
	whitesuff   = regexp.MustCompile(`[\s\p{Zs}]+$`)
	pref        = regexp.MustCompile(`^[('“"‘…]+`)
	suff        = regexp.MustCompile(`[)'“"‘…,!?\.]+$`)
	prefmap     = strings.NewReplacer("(", "(", "'", "’", "“", "„", "\"", "„", "‘", "’", "…", "…")
	suffmap     = strings.NewReplacer(")", ")", "'", "’", "“", "”", "\"", "”", "‘", "’", "…", "…", "!", "!", ",", ",", "?", "?", ".", ".")
	fixdigraphs = regexp.MustCompile(`\p{Lu}*(Dž|Nj|Lj)\p{Lu}+(Dž|Nj|Lj)?\p{Lu}*`)

	doit    = true
	version = "v0.2.0"

	rdr  = bufio.NewReader(os.Stdin)
	out  = bufio.NewWriter(os.Stdout)
	out1 = os.Stdout

	l2cPtr  = flag.Bool("l2c", false, "`Смер` пресловљавања је латиница у ћирилицу")
	c2lPtr  = flag.Bool("c2l", false, "`Смер` пресловљавања је ћирилица у латиницу")
	htmlPtr = flag.Bool("html", false, "`Формат` улаза је (X)HTML")
	textPtr = flag.Bool("text", false, "`Формат` улаза је прости текст")

	tbl = trie.BuildFromMap(map[string]string{
		"A":   "А",
		"B":   "Б",
		"V":   "В",
		"G":   "Г",
		"D":   "Д",
		"Đ":   "Ђ",
		"Ð":   "Ђ",
		"ᴆ":   "Ђ",
		"DJ":  "Ђ",
		"DЈ":  "Ђ", // D + cyrillic J
		"Dj":  "Ђ",
		"Dј":  "Ђ", // D + cyrillic j
		"E":   "Е",
		"Ž":   "Ж",
		"Ž":  "Ж", // Z with caron
		"Z":   "З",
		"I":   "И",
		"J":   "Ј",
		"K":   "К",
		"L":   "Л",
		"LJ":  "Љ",
		"LЈ":  "Љ", // L + cyrillic J
		"Ǉ":   "Љ",
		"Lj":  "Љ",
		"Lј":  "Љ", // L + cyrillic j
		"ǈ":   "Љ",
		"M":   "М",
		"N":   "Н",
		"NJ":  "Њ",
		"NЈ":  "Њ", // N + cyrillic J
		"Ǌ":   "Њ",
		"Nj":  "Њ",
		"Nј":  "Њ", // N + cyrillic j
		"ǋ":   "Њ",
		"O":   "О",
		"P":   "П",
		"R":   "Р",
		"S":   "С",
		"T":   "Т",
		"Ć":   "Ћ",
		"Ć":  "Ћ", // C with acute accent
		"U":   "У",
		"F":   "Ф",
		"H":   "Х",
		"C":   "Ц",
		"Č":   "Ч",
		"Č":  "Ч", // C with caron
		"DŽ":  "Џ",
		"Ǆ":   "Џ",
		"DŽ": "Џ", // D + Z with caron
		"Dž":  "Џ",
		"ǅ":   "Џ",
		"Dž": "Џ", // D + z with caron
		"Š":   "Ш",
		"Š":  "Ш", // S with caron
		"a":   "а",
		"æ":   "ае",
		"b":   "б",
		"v":   "в",
		"g":   "г",
		"d":   "д",
		"đ":   "ђ",
		"dj":  "ђ",
		"dј":  "ђ", // d + cyrillic j
		"e":   "е",
		"ž":   "ж",
		"ž":  "ж", // z with caron
		"z":   "з",
		"i":   "и",
		"ĳ":   "иј",
		"j":   "ј",
		"k":   "к",
		"l":   "л",
		"lj":  "љ",
		"lј":  "љ", // l + cyrillic j
		"ǉ":   "љ",
		"m":   "м",
		"n":   "н",
		"nj":  "њ",
		"nј":  "њ", // n + cyrillic j
		"ǌ":   "њ",
		"o":   "о",
		"œ":   "ое",
		"p":   "п",
		"r":   "р",
		"s":   "с",
		"ﬆ":   "ст",
		"t":   "т",
		"ć":   "ћ",
		"ć":  "ћ", // c with acute accent
		"u":   "у",
		"f":   "ф",
		"ﬁ":   "фи",
		"ﬂ":   "фл",
		"h":   "х",
		"c":   "ц",
		"č":   "ч",
		"č":  "ч", // c with caron
		"dž":  "џ",
		"ǆ":   "џ",
		"dž": "џ", // d + z with caron
		"š":   "ш",
		"š":  "ш", // s with caron
	})
	tbl1 = map[string]string{
		"А": "A",
		"Б": "B",
		"В": "V",
		"Г": "G",
		"Д": "D",
		"Ђ": "Đ",
		"Е": "E",
		"Ж": "Ž",
		"З": "Z",
		"И": "I",
		"Ј": "J",
		"К": "K",
		"Л": "L",
		"Љ": "Lj",
		"М": "M",
		"Н": "N",
		"Њ": "Nj",
		"О": "O",
		"П": "P",
		"Р": "R",
		"С": "S",
		"Т": "T",
		"Ћ": "Ć",
		"У": "U",
		"Ф": "F",
		"Х": "H",
		"Ц": "C",
		"Ч": "Č",
		"Џ": "Dž",
		"Ш": "Š",
		"а": "a",
		"б": "b",
		"в": "v",
		"г": "g",
		"д": "d",
		"ђ": "đ",
		"е": "e",
		"ж": "ž",
		"з": "z",
		"и": "i",
		"ј": "j",
		"к": "k",
		"л": "l",
		"љ": "lj",
		"м": "m",
		"н": "n",
		"њ": "nj",
		"о": "o",
		"п": "p",
		"р": "r",
		"с": "s",
		"т": "t",
		"ћ": "ć",
		"у": "u",
		"ф": "f",
		"х": "h",
		"ц": "c",
		"ч": "č",
		"џ": "dž",
		"ш": "š",
	}
	serbianWordsWithForeignCharacterCombinations = []string{
		"alchajmer",
		"ammar",
		"amss",
		"aparthejd",
		"ddor",
		"dss",
		"dvadesettrog",
		"epp",
		"fss",
		"gss",
		"interreakc",
		"interresor",
		"izzzdiovns",
		"kss",
		"llsls",
		"mmf",
		"naddu",
		"natha",
		"natho",
		"ommetar",
		"penthaus",
		"poddirektor",
		"poddisciplin",
		"poddomen",
		"poddres",
		"posthumn",
		"posttrans",
		"posttraum",
		"pothodni",
		"pothranj",
		"preddijabetes",
		"prethod",
		"ptt",
		"sbb",
		"sdss",
		"ssp",
		"ssrnj",
		"sssr",
		"superračun",
		"šopingholi",
		"tass",
		"transseks",
		"transsibir",
		"tridesettrog",
		"uppr",
		"vannastav",
	}
	commonForeignWords = []string{
		"administration",
		"adobe",
		"advanced",
		"advertising",
		"autocad",
		"bitcoin",
		"book",
		"boot",
		"cancel",
		"canon",
		"career",
		"carlsberg",
		"cisco",
		"clio",
		"cloud",
		"coca-col",
		"cookie",
		"cooking",
		"cool",
		"covid",
		"dacia",
		"default",
		"develop",
		"e-mail",
		"edge",
		"email",
		"emoji",
		"english",
		"facebook",
		"fashion",
		"food",
		"foundation",
		"gaming",
		"gmail",
		"gmt",
		"good",
		"google",
		"hdmi",
		"image",
		"iphon",
		"ipod",
		"javascript",
		"jazeera",
		"joomla",
		"league",
		"like",
		"linkedin",
		"look",
		"macbook",
		"mail",
		"manager",
		"maps",
		"mastercard",
		"mercator",
		"microsoft",
		"mitsubishi",
		"notebook",
		"nvidia",
		"online",
		"outlook",
		"panasonic",
		"pdf",
		"peugeot",
		"podcast",
		"postpaid",
		"printscreen",
		"procredit",
		"project",
		"punk",
		"renault",
		"rock",
		"screenshot",
		"seen",
		"selfie",
		"share",
		"shift",
		"shop",
		"smartphone",
		"space",
		"steam",
		"stream",
		"subscrib",
		"timeout",
		"tool",
		"topic",
		"trailer",
		"ufc",
		"unicredit",
		"username",
		"viber",
	}
	wholeForeignWords = []string{
		"about",
		"air",
		"alpha",
		"and",
		"back",
		"bitcoin",
		"brainz",
		"celebrities",
		"co2",
		"conditions",
		"cpu",
		"creative",
		"disclaimer",
		"discord",
		"dj",
		"electronics",
		"entertainment",
		"files",
		"fresh",
		"fun",
		"geographic",
		"gmbh",
		"h2o",
		"hair",
		"have",
		"home",
		"idj",
		"idjtv",
		"ii",
		"iii",
		"life",
		"latest",
		"live",
		"login",
		"made",
		"makeup",
		"must",
		"national",
		"previous",
		"public",
		"reserved",
		"score",
		"screen",
		"terms",
		"url",
		"vii",
		"viii",
		"visa",
	}
	foreignCharacterCombinations = []string{
		`q`,
		`w`,
		`x`,
		`y`,
		`ü`,
		`ö`,
		`ä`,
		`ø`,
		`ß`,
		`&`,
		`@`,
		`#`,
		`bb`,
		`cc`,
		`dd`,
		`ff`,
		`gg`,
		`hh`,
		`kk`,
		`ll`,
		`mm`,
		`nn`,
		`pp`,
		`rr`,
		`ss`,
		`tt`,
		`zz`,
		`ch`,
		`gh`,
		`th`,
		`'s`,
		`'t`,
		`.com`,
		`.net`,
		`.info`,
		`.rs`,
		`.org`,
		`©`,
		`®`,
		`™`,
	}
	digraphExceptions = map[string][]string{
		"dj": {
			"adjektiv",
			"adjunkt",
			"autodjel",
			"bazdje",
			"bdje",
			"bezdje",
			"blijedje",
			"bludje",
			"bridjе",
			"vidjel",
			"vidjet",
			"vindjakn",
			"višenedje",
			"vrijedje",
			"gdje",
			"gudje",
			"gdjir",
			"daždje",
			"dvonedje",
			"devetonedje",
			"desetonedje",
			"djb",
			"djeva",
			"djevi",
			"djevo",
			"djed",
			"djejstv",
			"djel",
			"djenem",
			"djeneš",
			//"djene" rare (+ Дјене (town)), but it would colide with ђене-ђене, ђеневљанка, ђенерал итд.
			"djenu",
			"djet",
			"djec",
			"dječ",
			"djuar",
			"djubison",
			"djubouz",
			"djuer",
			"djui",
			// "djuk", djuk (engl. Duke) косило би се нпр. са Djukanović
			"djuks",
			"djulej",
			"djumars",
			"djupont",
			"djurant",
			"djusenberi",
			"djuharst",
			"djuherst",
			"dovdje",
			"dogrdje",
			"dodjel",
			"drvodje",
			"drugdje",
			"elektrosnabdje",
			"žudje",
			"zabludje",
			"zavidje",
			"zavrijedje",
			"zagudje",
			"zadjev",
			"zadjen",
			"zalebdje",
			"zaludje",
			"zaodje",
			"zapodje",
			"zarudje",
			"zasjedje",
			"zasmrdje",
			"zastidje",
			"zaštedje",
			"zdje",
			"zlodje",
			"igdje",
			"izbledje",
			"izblijedje",
			"izvidje",
			"izdjejst",
			"izdjelj",
			"izludje",
			"isprdje",
			"jednonedje",
			"kojegdje",
			"kudjelj",
			"lebdje",
			"ludjel",
			"ludjet",
			"makfadjen",
			"marmadjuk",
			"međudjel",
			"nadjaha",
			"nadjača",
			"nadjeb",
			"nadjev",
			"nadjenul",
			"nadjenuo",
			"nadjenut",
			"negdje",
			"nedjel",
			"nadjunač",
			"nenadjača",
			"nenadjebi",
			"nenavidje",
			"neodje",
			"nepodjarm",
			"nerazdje",
			"nigdje",
			"obdjel",
			"obnevidje",
			"ovdje",
			"odjav",
			"odjah",
			"odjaš",
			"odjeb",
			"odjev",
			"odjed",
			"odjezd",
			"odjek",
			"odjel",
			"odjen",
			"odjeć",
			"odjec",
			"odjur",
			"odsjedje",
			"ondje",
			"opredje",
			"osijedje",
			"osmonedje",
			"pardju",
			"perdju",
			"petonedje",
			"poblijedje",
			"povidje",
			"pogdjegdje",
			"pogdje",
			"podjakn",
			"podjamč",
			"podjastu",
			"podjemč",
			"podjar",
			"podjeb",
			"podjed",
			"podjezič",
			"podjel",
			"podjen",
			"podjet",
			"pododjel",
			"pozavidje",
			"poludje",
			"poljodjel",
			"ponegdje",
			"ponedjelj",
			"porazdje",
			"posijedje",
			"posjedje",
			"postidje",
			"potpodjel",
			"poštedje",
			"pradjed",
			"prdje",
			"preblijedje",
			"previdje",
			"predvidje",
			"predjel",
			"preodjen",
			"preraspodje",
			"presjedje",
			"pridjev",
			"pridjen",
			"prismrdje",
			"prištedje",
			"probdje",
			"problijedje",
			"prodjen",
			"prolebdje",
			"prosijedje",
			"prosjedje",
			"protivdjel",
			"prošlonedje",
			"radjard",
			"razvidje",
			"razdjev",
			"razdjel",
			"razodje",
			"raspodje",
			"rasprdje",
			"remekdjel",
			"rudjen",
			"rudjet",
			"sadje",
			"svagdje",
			"svidje",
			"svugdje",
			"sedmonedjelj",
			"sijedje",
			"sjedje",
			"smrdje",
			"snabdje",
			"snovidje",
			"starosjedje",
			"stidje",
			"studje",
			"sudjel",
			"tronedje",
			"ublijedje",
			"uvidje",
			"udjel",
			"udjen",
			"uprdje",
			"usidjel",
			"usjedje",
			"usmrdje",
			"uštedje",
			"cjelonedje",
			"četvoronedje",
			"čukundjed",
			"šestonedjelj",
			"štedje",
			"štogdje",
			"šukundjed",
		},
		"dž": {
			"feldžandarm",
			"nadžanj",
			"nadždrel",
			"nadžel",
			"nadžeo",
			"nadžet",
			"nadživ",
			"nadžinj",
			"nadžnj",
			"nadžrec",
			"nadžup",
			"odžali",
			"odžari",
			"odžel",
			"odžive",
			"odživljava",
			"odžubor",
			"odžvaka",
			"odžval",
			"odžvać",
			"podžanr",
			"podžel",
			"podže",
			"podžig",
			"podžiz",
			"podžil",
			"podžnje",
			"podžupan",
			"predželu",
			"predživot",
		},
		"nj": {
			"anjon",
			"injaric",
			"injekc",
			"injekt",
			"injicira",
			"injurij",
			"kenjon",
			"konjug",
			"konjunk",
			"nekonjug",
			"nekonjunk",
			"ssrnj",
			"tanjug",
			"vanjezičk",
		},
	}

	// See: https://en.wikipedia.org/wiki/Zero-width_non-joiner
	digraphReplacements = map[string]map[string]string{
		"dj": {
			"dj": "d\u200Cj",
			"Dj": "D\u200Cj",
			"DJ": "D\u200CJ",
		},
		"dž": {
			"dž": "d\u200Cž",
			"Dž": "D\u200Cž",
			"DŽ": "D\u200CŽ",
		},
		"nj": {
			"nj": "n\u200Cj",
			"Nj": "N\u200Cj",
			"NJ": "N\u200CJ",
		},
	}
)

func looksLikeForeignWord(word string) bool {
	trimmedWord := trimExcessiveCharacters(word)
	processed := strings.ToLower(trimmedWord)
	if processed == "" {
		return false
	}

	if wordStartsWith(processed, serbianWordsWithForeignCharacterCombinations) {
		return false
	}

	if wordContainsString(processed, foreignCharacterCombinations) {
		return true
	}

	if wordStartsWith(processed, commonForeignWords) {
		return true
	}

	if wordIsEqualTo(processed, wholeForeignWords) {
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
	for digraph := range digraphExceptions {
		if !strings.Contains(lowercaseStr, digraph) {
			continue
		}
		for _, word := range digraphExceptions[digraph] {
			if !strings.HasPrefix(lowercaseStr, word) {
				continue
			}
			// Split all possible occurrences, regardless of case.
			for key, word := range digraphReplacements[digraph] {
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
	word = pref.ReplaceAllStringFunc(word, prefmap.Replace)
	word = suff.ReplaceAllStringFunc(word, suffmap.Replace)

	word = splitDigraphs(word)
	if doit {
		for i, runeValue := range word {
			if w > 1 {
				w -= 1
				continue
			}
			value, prefixLen, ok := tbl.SearchPrefixInString(word[i:])
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
		doit = true
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
	word = pref.ReplaceAllStringFunc(word, prefmap.Replace)
	word = suff.ReplaceAllStringFunc(word, suffmap.Replace)

	for _, runeValue := range word {
		value, ok := tbl1[string(runeValue)]
		if ok {
			result += value
		} else {
			result += string(runeValue)
		}
	}
	if fixdigraphs.MatchString(result) {
		return Uppercase(result)
	} else {
		return result
	}
}

func Pomoc() {
	fmt.Fprintf(flag.CommandLine.Output(), "Ово је филтер %s верзија %s\nСаставио eevan78, 2024\n\n", os.Args[0], version)
	fmt.Fprintf(flag.CommandLine.Output(), "Филтер чита UTF-8 кодирани текст са стандардног улаза и исписује га на\nстандардни излаз пресловљен сагласно са наведеним заставицама:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nМора да се наведе по једна и само једна заставица из обе групе Смер и Формат.\nЦеле речи између „<|” и „|>” у простом тексту се не пресловљавају у ћирилицу.\nТекст унутар <span lang=\"sr-Latn\"></span> елемента у (X)HTML се не пресловљава у\nћирилицу, а текст унутар <span lang=\"sr-Cyrl\"></span> се не пресловљава у латиницу.\n\nПримери:\n%s -l2c -html\t\tпреслови (X)HTML у ћирилицу\n%s -text -c2l\t\tпреслови прости текст у латиницу\n", os.Args[0], os.Args[0])
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
		namespace := ""
		notexist := true
		// Properly adjust the lang attribute, or add it if it's missing
		if n.Data == "html" {
			if *l2cPtr {
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
					n.Attr = append(n.Attr, html.Attribute{namespace, "lang", "sr-Cyrl-t-sr-Latn"})
				}
			} else if *c2lPtr {
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
					n.Attr = append(n.Attr, html.Attribute{namespace, "lang", "sr-Latn-t-sr-Cyrl"})
				}
			}
		}
	case html.TextNode:
		// Transliterate if text is not inside a script or a style element
		if !allWhite(n.Data) && n.Parent.Type == html.ElementNode && (n.Parent.Data != "script" && n.Parent.Data != "style") {
			nodeprefix := whitepref.FindString(n.Data)
			nodesuffix := whitesuff.FindString(n.Data)
			words := strings.Fields(n.Data)
			for w := range words {
				// Do not transliterate to cyrillic if the parent element is span that has a lang attribute with the value "sr-Latn"
				if *l2cPtr == true && !(n.Parent.Data == "span" && (n.Parent.Attr[0].Key == "lang" || n.Parent.Attr[0].Key == "xml:lang") && n.Parent.Attr[0].Val == "sr-Latn") {
					index := transliterationIndexOfWordStartsWith(strings.ToLower(words[w]), wholeForeignWords, "-")
					if index >= 0 {
						words[w] = string(words[w][:index]) + l2c(string(words[w][index:]))
					} else if !looksLikeForeignWord(words[w]) {
						words[w] = l2c(words[w])
					}
					// Do not transliterate to latin if the parent element is span that has a lang attribut with the value "sr-Cyrl"
				} else if *c2lPtr == true && !(n.Parent.Data == "span" && (n.Parent.Attr[0].Key == "lang" || n.Parent.Attr[0].Key == "xml:lang") && n.Parent.Attr[0].Val == "sr-Cyrl") {
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

func main() {
	flag.Usage = Pomoc
	flag.Parse()
	if !(((!*l2cPtr && *c2lPtr) || (*l2cPtr && !*c2lPtr)) && ((!*htmlPtr && *textPtr) || (*htmlPtr && !*textPtr)) && flag.NFlag() == 2) {
		Pomoc()
		os.Exit(0)
	}

	if *htmlPtr {
		doc, err := html.Parse(rdr)
		if err != nil {
			fmt.Fprintln(os.Stdout, "Грешка:", err)
			os.Exit(1)
		}
		traverseNode(doc)
		if err := html.Render(out, doc); err != nil {
			fmt.Fprintln(os.Stdout, "Грешка:", err)
			os.Exit(1)
		}
		_ = out.Flush()
	} else if *textPtr {
		for {
			switch line, err := rdr.ReadString('\n'); err {
			case nil:
				lineprefix := whitepref.FindString(line)
				words := strings.Fields(line)
				for n := range words {
					if strings.HasPrefix(words[n], "<|") {
						doit = false
						words[n] = strings.TrimPrefix(words[n], "<|")
					}
					if *l2cPtr {
						index := transliterationIndexOfWordStartsWith(strings.ToLower(words[n]), wholeForeignWords, "-")
						if index >= 0 {
							words[n] = string(words[n][:index]) + l2c(string(words[n][index:]))
						} else if !looksLikeForeignWord(words[n]) {
							words[n] = l2c(words[n])
						}
					} else if *c2lPtr {
						words[n] = c2l(words[n])
					}
				}
				// Preserve the line indentation
				if lineprefix != "" && lineprefix != "\n" && len(words) != 0 {
					words[0] = lineprefix + words[0]
				}
				outl := strings.Join(words, " ")
				outl += "\n"
				if _, err = out1.WriteString(outl); err != nil {
					fmt.Fprintln(os.Stderr, "Грешка:", err)
					os.Exit(1)
				}

			case io.EOF:
				os.Exit(0)

			default:
				fmt.Fprintln(os.Stderr, "Грешка:", err)
				os.Exit(1)
			}
		}
	}
}
