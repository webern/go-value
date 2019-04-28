// go-value, Copyright (c) 2019 by Matthew James Briggs

package value

import (
	"bytes"
	"sort"
	"strings"
	"unicode"
)

func ToCamelFromWords(input []string) string {
	return doCaseFromWords(input, false)
}

func ToPascalFromWords(input []string) string {
	return doCaseFromWords(input, true)
}

func doCaseFromWords(input []string, isPascal bool) string {
	buf := bytes.Buffer{}
	isFirst := true

	for _, word := range input {

		if isFirst && !isPascal {
			buf.WriteString(Symbolize(strings.ToLower(word), false))
		} else if isFirst && isPascal {
			buf.WriteString(Symbolize(Title(word), false))
		} else {
			buf.WriteString(Symbolize(Title(word), false))
		}

		isFirst = false
	}

	result := buf.String()
	return result
}

func doParse(input string, dictionary []string) []string {
	tokens := Tokenize(input)
	words := make([]string, 0)

	for _, token := range tokens {
		currentWords := Words(token, dictionary)
		words = append(words, currentWords...)
	}

	return words
}

func Title(input string) string {
	if len(input) <= 1 {
		return strings.ToUpper(input)
	}
	s := strings.ToLower(input)
	start := strings.ToUpper(s[:1])
	rest := s[1:]
	return start + rest
}

func Symbolize(input string, canStartWithDigit bool) string {
	if len(input) == 0 {
		return ""
	}

	first := input[:1]

	var firstRune rune
	buf := bytes.Buffer{}
	for ix, c := range input {
		if ix == 0 {
			firstRune = c
		} else {
			if unicode.IsLetter(c) || unicode.IsNumber(c) {
				buf.WriteRune(c)
			}
		}
	}

	if !unicode.IsLetter(firstRune) && !unicode.IsNumber(firstRune) {
		first = "Column"
	}

	rest := buf.String()

	if !canStartWithDigit {
		switch first {
		case "0":
			first = "Zero"
		case "1":
			first = "One"
		case "2":
			first = "Two"
		case "3":
			first = "Three"
		case "4":
			first = "Four"
		case "5":
			first = "Five"
		case "6":
			first = "Six"
		case "7":
			first = "Seven"
		case "8":
			first = "Eight"
		case "9":
			first = "Nine"
		}
	}

	return first + rest
}

func ToCamel(input string, dictionary []string) string {
	words := doParse(input, dictionary)
	return ToCamelFromWords(words)
}

func ToPascal(input string, dictionary []string) string {
	words := doParse(input, dictionary)
	return ToPascalFromWords(words)
}

// Tokenize breaks up a string at any non alphanumeric boundries such as whitespaces, dashes or underscores
func Tokenize(input string) []string {
	buf := bytes.Buffer{}
	result := make([]string, 0, 0)

	for _, c := range input {
		partOfWord := unicode.IsLetter(c) || unicode.IsNumber(c)

		if partOfWord {
			buf.WriteString(string(c))
		} else {
			current := buf.String()

			if len(current) > 0 {
				result = append(result, current)
			}

			buf = bytes.Buffer{}
		}
	}

	current := buf.String()

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

// Words takes a list of strings, and further breaks those strings into words by searching the given dictionary for
// matches, it is recursive to find and split all available words in the string
func Words(input string, dictionary []string) []string {
	result := make([]string, 0)
	words(input, dictionary, &result)
	return result
}

func words(input string, dictionary []string, result *[]string) {
	temp := SplitIntoWords(input, dictionary)

	if len(temp) == 1 {
		*result = append(*result, temp...)
		return
	}

	for _, t := range temp {
		recursiveResult := Words(t, dictionary)
		*result = append(*result, recursiveResult...)
	}

	return
}

func SplitIntoWords(input string, dictionary []string) []string {
	lower := strings.ToLower(input)
	result := make([]string, 0)
	for _, dictWord := range dictionary {
		if len(dictWord) > len(input) {
			continue
		}

		ix := strings.Index(lower, dictWord)

		if ix < 0 {
			continue
		}

		if ix == 0 && len(dictWord) == len(lower) {
			result = append(result, input)
			return result
		}

		foundWordEnd := ix + len(dictWord)
		word1 := input[:ix]
		word2 := input[ix:foundWordEnd]
		word3 := ""

		if foundWordEnd < len(lower) {
			word3 = input[foundWordEnd:]
		}

		if len(word1) > 0 {
			result = append(result, word1)
		}

		result = append(result, word2)

		if len(word3) > 0 {
			result = append(result, word3)
		}

		return result
	}

	// not found
	result = append(result, input)
	return result
}

// CreateDictionary adds wordsToAdd and removes wordsToRemove from the default dictionary, which is the list of the top
// 1000 English words words with length greater than 2. wordsToAdd will be placed at the highest priority order. This is
// inefficient, you should call this function infrequently in your program.
//
func CreateDictionary(wordsToAdd, wordsToRemove []string) []string {
	s := defaultDictionary[:]

	for ix := range wordsToAdd {
		wordsToAdd[ix] = strings.ToLower(wordsToAdd[ix])
	}

	for ix := range wordsToRemove {
		wordsToRemove[ix] = strings.ToLower(wordsToRemove[ix])
	}

	for ix := range s {
		s[ix] = strings.ToLower(s[ix])
	}

	if wordsToAdd != nil {
		sort.Sort(byReverseLength(wordsToAdd))
		temp := s
		temp = append(wordsToAdd, s...)
		s = temp
	}

	if wordsToRemove != nil {
		for _, wordToRemove := range wordsToRemove {
			for ix, currentWord := range s {
				if wordToRemove == currentWord {
					s = append(s[:ix], s[ix+1:]...)
					break
				}
			}
		}
	}

	SortByLengthReverse(s)
	return s
}

// top 1000 words in english according to https://www.ef.edu/english-resources/english-vocabulary/top-1000-words/
var defaultDictionary = [...]string{
	"ability",
	"able",
	"about",
	"above",
	"accept",
	"according",
	"account",
	"across",
	"act",
	"action",
	"activity",
	"actually",
	"add",
	"address",
	"administration",
	"admit",
	"adult",
	"affect",
	"after",
	"again",
	"against",
	"age",
	"agency",
	"agent",
	"ago",
	"agree",
	"agreement",
	"ahead",
	"air",
	"all",
	"allow",
	"almost",
	"alone",
	"along",
	"already",
	"also",
	"although",
	"always",
	"American",
	"among",
	"amount",
	"analysis",
	"animal",
	"another",
	"answer",
	"any",
	"anyone",
	"anything",
	"appear",
	"applied",
	"apply",
	"approach",
	"area",
	"argue",
	"arm",
	"around",
	"arrive",
	"art",
	"article",
	"artist",
	"ask",
	"assume",
	"attack",
	"attention",
	"attorney",
	"audience",
	"author",
	"authority",
	"available",
	"avoid",
	"away",
	"baby",
	"back",
	"bad",
	"bag",
	"ball",
	"bank",
	"bar",
	"base",
	"beat",
	"beautiful",
	"because",
	"become",
	"bed",
	"before",
	"begin",
	"behavior",
	"behind",
	"believe",
	"benefit",
	"best",
	"better",
	"between",
	"beyond",
	"big",
	"bill",
	"billion",
	"bit",
	"black",
	"blend",
	"blood",
	"blue",
	"board",
	"body",
	"book",
	"born",
	"both",
	"box",
	"boy",
	"break",
	"bring",
	"brother",
	"brown",
	"budget",
	"build",
	"building",
	"business",
	"but",
	"buy",
	"call",
	"camera",
	"campaign",
	"can",
	"cancer",
	"candidate",
	"capital",
	"card",
	"care",
	"career",
	"carrier",
	"carry",
	"case",
	"catch",
	"cause",
	"cell",
	"center",
	"central",
	"century",
	"certain",
	"certainly",
	"chair",
	"challenge",
	"chance",
	"change",
	"character",
	"charge",
	"check",
	"child",
	"choice",
	"choose",
	"church",
	"citizen",
	"city",
	"civil",
	"claim",
	"class",
	"clear",
	"clearly",
	"close",
	"coach",
	"cold",
	"collection",
	"college",
	"color",
	"come",
	"commercial",
	"common",
	"community",
	"company",
	"compare",
	"computer",
	"concern",
	"condition",
	"conference",
	"Congress",
	"consider",
	"consumer",
	"contain",
	"continue",
	"control",
	"cost",
	"could",
	"country",
	"couple",
	"course",
	"court",
	"cover",
	"create",
	"crime",
	"cultural",
	"culture",
	"cup",
	"current",
	"customer",
	"cut",
	"dark",
	"data",
	"daughter",
	"day",
	"dead",
	"deal",
	"death",
	"debate",
	"decade",
	"decide",
	"decision",
	"deep",
	"defense",
	"degree",
	"Democrat",
	"democratic",
	"describe",
	"design",
	"despite",
	"detail",
	"detailed",
	"determine",
	"develop",
	"development",
	"die",
	"difference",
	"different",
	"difficult",
	"dinner",
	"direction",
	"director",
	"discover",
	"discuss",
	"discussion",
	"disease",
	"doctor",
	"dog",
	"door",
	"down",
	"draw",
	"dream",
	"drive",
	"drop",
	"drug",
	"during",
	"each",
	"early",
	"east",
	"easy",
	"eat",
	"economic",
	"economy",
	"edge",
	"education",
	"effect",
	"effort",
	"eight",
	"either",
	"election",
	"else",
	"employee",
	"end",
	"energy",
	"enjoy",
	"enough",
	"enter",
	"entire",
	"environment",
	"environmental",
	"especially",
	"establish",
	"even",
	"evening",
	"event",
	"ever",
	"every",
	"everybody",
	"everyone",
	"everything",
	"evidence",
	"exactly",
	"example",
	"executive",
	"exist",
	"expect",
	"experience",
	"expert",
	"explain",
	"eye",
	"face",
	"fact",
	"factor",
	"fail",
	"fall",
	"family",
	"far",
	"fast",
	"father",
	"fear",
	"federal",
	"feel",
	"feeling",
	"few",
	"field",
	"fight",
	"figure",
	"fill",
	"film",
	"final",
	"finally",
	"financial",
	"find",
	"fine",
	"finger",
	"finish",
	"fire",
	"firm",
	"first",
	"fish",
	"five",
	"floor",
	"fly",
	"focus",
	"follow",
	"food",
	"foot",
	"for",
	"force",
	"foreign",
	"forget",
	"form",
	"former",
	"forward",
	"four",
	"fox",
	"free",
	"friend",
	"from",
	"front",
	"full",
	"fund",
	"future",
	"game",
	"garden",
	"gas",
	"general",
	"generation",
	"get",
	"girl",
	"give",
	"glass",
	"goal",
	"good",
	"government",
	"great",
	"green",
	"ground",
	"group",
	"grow",
	"growth",
	"guess",
	"gun",
	"guy",
	"hair",
	"half",
	"hammer",
	"hmmer",
	"hmr",
	"Hmrr",
	"hand",
	"hang",
	"happen",
	"happy",
	"hard",
	"have",
	"head",
	"health",
	"hear",
	"heart",
	"heat",
	"heavy",
	"hello",
	"help",
	"her",
	"here",
	"herself",
	"high",
	"him",
	"himself",
	"his",
	"history",
	"hit",
	"hold",
	"home",
	"hope",
	"hospital",
	"hot",
	"hotel",
	"hour",
	"house",
	"how",
	"however",
	"huge",
	"human",
	"hundred",
	"husband",
	"idea",
	"identify",
	"image",
	"imagine",
	"impact",
	"important",
	"improve",
	"include",
	"including",
	"increase",
	"indeed",
	"index",
	"indicate",
	"individual",
	"industry",
	"information",
	"inside",
	"instead",
	"institution",
	"interest",
	"interesting",
	"international",
	"interview",
	"into",
	"investment",
	"involve",
	"issue",
	"item",
	"its",
	"itself",
	"job",
	"join",
	"just",
	"keep",
	"key",
	"kid",
	"kill",
	"kind",
	"kitchen",
	"know",
	"knowledge",
	"land",
	"language",
	"large",
	"last",
	"late",
	"later",
	"laugh",
	"law",
	"lawyer",
	"lay",
	"lead",
	"leader",
	"learn",
	"least",
	"leave",
	"left",
	"leg",
	"legal",
	"less",
	"let",
	"letter",
	"level",
	"lie",
	"life",
	"light",
	"like",
	"likely",
	"line",
	"list",
	"listen",
	"little",
	"live",
	"local",
	"long",
	"look",
	"lose",
	"loss",
	"lot",
	"love",
	"low",
	"machine",
	"magazine",
	"main",
	"maintain",
	"major",
	"majority",
	"make",
	"man",
	"manage",
	"management",
	"manager",
	"many",
	"market",
	"marriage",
	"material",
	"matter",
	"may",
	"maybe",
	"mean",
	"measure",
	"media",
	"medical",
	"meet",
	"meeting",
	"member",
	"memory",
	"mention",
	"message",
	"method",
	"middle",
	"might",
	"military",
	"million",
	"mind",
	"minute",
	"miss",
	"mission",
	"model",
	"modern",
	"moment",
	"money",
	"month",
	"more",
	"morning",
	"most",
	"mother",
	"mouth",
	"move",
	"movement",
	"movie",
	"Mrs",
	"much",
	"music",
	"must",
	"myself",
	"name",
	"nation",
	"national",
	"natural",
	"nature",
	"near",
	"nearly",
	"necessary",
	"need",
	"network",
	"never",
	"new",
	"news",
	"newspaper",
	"next",
	"nice",
	"night",
	"none",
	"nor",
	"north",
	"not",
	"note",
	"nothing",
	"notice",
	"now",
	"number",
	"occur",
	"off",
	"offer",
	"office",
	"office",
	"officer",
	"official",
	"often",
	"oil",
	"old",
	"once",
	"one",
	"only",
	"onto",
	"open",
	"operation",
	"opportunity",
	"option",
	"order",
	"organization",
	"other",
	"others",
	"our",
	"out",
	"outside",
	"over",
	"own",
	"owner",
	"page",
	"pain",
	"painting",
	"paper",
	"parent",
	"part",
	"participant",
	"particular",
	"particularly",
	"partner",
	"party",
	"pass",
	"past",
	"patient",
	"pattern",
	"pay",
	"peace",
	"people",
	"per",
	"perform",
	"performance",
	"perhaps",
	"period",
	"person",
	"personal",
	"phone",
	"physical",
	"pick",
	"picture",
	"piece",
	"place",
	"plan",
	"plant",
	"play",
	"player",
	"point",
	"police",
	"policy",
	"political",
	"politics",
	"poor",
	"popular",
	"population",
	"position",
	"positive",
	"possible",
	"power",
	"practice",
	"prepare",
	"present",
	"president",
	"pressure",
	"pretty",
	"prevent",
	"price",
	"private",
	"probably",
	"problem",
	"process",
	"produce",
	"product",
	"production",
	"professional",
	"professor",
	"program",
	"project",
	"property",
	"protect",
	"prove",
	"provide",
	"public",
	"pull",
	"purpose",
	"push",
	"put",
	"quality",
	"question",
	"quick",
	"quickly",
	"quite",
	"race",
	"radio",
	"raise",
	"range",
	"rate",
	"rather",
	"reach",
	"read",
	"ready",
	"real",
	"reality",
	"realize",
	"really",
	"reason",
	"receive",
	"recent",
	"recently",
	"recognize",
	"record",
	"red",
	"reduce",
	"reflect",
	"region",
	"relate",
	"relationship",
	"religious",
	"remain",
	"remember",
	"remove",
	"report",
	"represent",
	"Republican",
	"require",
	"research",
	"resource",
	"respond",
	"response",
	"responsibility",
	"rest",
	"result",
	"retention",
	"return",
	"reveal",
	"revenue",
	"rich",
	"right",
	"rise",
	"risk",
	"road",
	"rock",
	"role",
	"room",
	"route",
	"routing",
	"row",
	"rule",
	"run",
	"safe",
	"same",
	"save",
	"say",
	"scene",
	"school",
	"science",
	"scientist",
	"score",
	"sea",
	"season",
	"seat",
	"second",
	"section",
	"security",
	"see",
	"seek",
	"seem",
	"sell",
	"send",
	"senior",
	"sense",
	"series",
	"serious",
	"serve",
	"service",
	"set",
	"seven",
	"several",
	"sex",
	"sexual",
	"shake",
	"share",
	"she",
	"shoot",
	"short",
	"shot",
	"should",
	"shoulder",
	"show",
	"side",
	"sign",
	"significant",
	"similar",
	"simple",
	"simply",
	"since",
	"sing",
	"single",
	"sister",
	"sit",
	"site",
	"situation",
	"six",
	"size",
	"skill",
	"skin",
	"small",
	"smile",
	"social",
	"society",
	"soldier",
	"some",
	"somebody",
	"someone",
	"something",
	"sometimes",
	"son",
	"song",
	"soon",
	"sort",
	"sound",
	"source",
	"south",
	"southern",
	"space",
	"speak",
	"special",
	"specific",
	"speech",
	"spend",
	"sport",
	"spring",
	"staff",
	"stage",
	"stand",
	"standard",
	"star",
	"start",
	"state",
	"statement",
	"station",
	"stay",
	"step",
	"still",
	"stock",
	"stop",
	"store",
	"story",
	"strategy",
	"street",
	"strong",
	"structure",
	"student",
	"study",
	"stuff",
	"style",
	"subject",
	"success",
	"successful",
	"such",
	"suddenly",
	"suffer",
	"suggest",
	"suggested",
	"summer",
	"support",
	"sure",
	"surface",
	"system",
	"table",
	"take",
	"talk",
	"tandem",
	"target",
	"tariff",
	"task",
	"tax",
	"teach",
	"teacher",
	"team",
	"technology",
	"television",
	"tell",
	"ten",
	"tend",
	"term",
	"test",
	"than",
	"thank",
	"that",
	"the",
	"their",
	"them",
	"themselves",
	"then",
	"theory",
	"there",
	"these",
	"they",
	"thing",
	"think",
	"third",
	"this",
	"those",
	"though",
	"thought",
	"thousand",
	"threat",
	"three",
	"through",
	"throughout",
	"throw",
	"thus",
	"time",
	"today",
	"together",
	"tonight",
	"too",
	"top",
	"total",
	"tough",
	"toward",
	"town",
	"trade",
	"traditional",
	"training",
	"transport",
	"travel",
	"treat",
	"treatment",
	"tree",
	"trial",
	"trip",
	"trouble",
	"TRUE",
	"truth",
	"try",
	"turn",
	"two",
	"type",
	"under",
	"understand",
	"unit",
	"until",
	"upon",
	"use",
	"usually",
	"value",
	"various",
	"very",
	"victim",
	"view",
	"violence",
	"visit",
	"voice",
	"vote",
	"wait",
	"walk",
	"wall",
	"want",
	"war",
	"watch",
	"water",
	"way",
	"weapon",
	"wear",
	"week",
	"weight",
	"well",
	"west",
	"western",
	"what",
	"whatever",
	"when",
	"where",
	"whether",
	"which",
	"while",
	"white",
	"who",
	"whole",
	"whom",
	"whose",
	"why",
	"wide",
	"wife",
	"will",
	"win",
	"wind",
	"window",
	"wish",
	"with",
	"within",
	"without",
	"woman",
	"wonder",
	"word",
	"work",
	"worker",
	"world",
	"worry",
	"would",
	"write",
	"writer",
	"wrong",
	"yard",
	"yeah",
	"year",
	"yes",
	"yet",
	"you",
	"young",
	"your",
	"yourself",
}
