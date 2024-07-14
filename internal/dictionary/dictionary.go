package dictionary

import (
	"bufio"
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/porfirion/trie"
)

var (
	Whitepref   = regexp.MustCompile(`^[\s\p{Zs}]+`)
	Whitesuff   = regexp.MustCompile(`[\s\p{Zs}]+$`)
	Pref        = regexp.MustCompile(`^[('“”"‘…]+`)
	Suff        = regexp.MustCompile(`[)'“"‘…,!?\.]+$`)
	Prefmap     = strings.NewReplacer("(", "(", "'", "’", "“", "„", "”", "„", "\"", "„", "‘", "’", "…", "…")
	Suffmap     = strings.NewReplacer(")", ")", "'", "’", "“", "”", "\"", "”", "‘", "’", "…", "…", "!", "!", ",", ",", "?", "?", ".", ".")
	Fixdigraphs = regexp.MustCompile(`\p{Lu}*(Dž|Nj|Lj)\p{Lu}+(Dž|Nj|Lj)?\p{Lu}*`)

	Doit    = true
	Version = "v0.3.0"

	Rdr             = bufio.NewReader(os.Stdin)
	Out             = bufio.NewWriter(os.Stdout)
	InputFilenames  []string
	InputFilePaths  []string
	OutputFilePaths []string
	OutputDir       = "../../output"

	L2cPtr       = flag.Bool("l2c", false, "`Смер` пресловљавања је латиница у ћирилицу")
	C2lPtr       = flag.Bool("c2l", false, "`Смер` пресловљавања је ћирилица у латиницу")
	HtmlPtr      = flag.Bool("html", false, "`Формат` улаза је (X)HTML")
	TextPtr      = flag.Bool("text", false, "`Формат` улаза је прости текст")
	InputPathPtr = flag.String("i", "", "Путања улазног фајла или директоријума")

	Tbl = trie.BuildFromMap(map[string]string{
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
	Tbl1 = map[string]string{
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
	SerbianWordsWithForeignCharacterCombinations = []string{
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
	CommonForeignWords = []string{
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
	WholeForeignWords = []string{
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
		"green",
		"h2o",
		"hair",
		"have",
		"home",
		"idj",
		"idjtv",
		"ii",
		"iii",
		"inclusive",
		"like",
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
	ForeignCharacterCombinations = []string{
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
	DigraphExceptions = map[string][]string{
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
	DigraphReplacements = map[string]map[string]string{
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
