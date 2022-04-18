package errorGenerator

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

var nearKeyboardLetters map[string]string
var randomSeed int

func init() {
	randomSeed = 0
	//все соседние буквы на клавиатуре
	combinations := map[string]string{
		"й": "фыц",
		"ц": "йфыву",
		"у": "цывак",
		"к": "увапе",
		"е": "капрн",
		"н": "епрог",
		"г": "нролш",
		"ш": "голдщ",
		"щ": "шлджз",
		"з": "щджэх",
		"х": "зжээ\\ъ",
		"ъ": "хэ\\",

		"ф": "йцыя",
		"ы": "йфячвуц",
		"в": "цычсаку",
		"а": "увсмпек",
		"п": "амирнек",
		"р": "епитогн",
		"о": "ртьлшгн",
		"л": "гоьбдщш",
		"д": "шлбюжзщ",
		"ж": "дю.эхзщ",
		"э": "ж.\\ъхз",

		"я": "фыч",
		"ч": "яфывс",
		"с": "чывам",
		"м": "свапи",
		"и": "мапрт",
		"т": "ипроь",
		"ь": "тролб",
		"б": "ьолдю",
		"ю": "блджэ.",
	}

	//только соседние буквы на одной строке
	combinations2 := map[string]string{
		"й": "цф", "q": "wa",
		"ц": "йу", "w": "qe",
		"у": "цк", "e": "wr",
		"к": "уе", "r": "et",
		"е": "кн", "t": "ry",
		"н": "ег", "y": "tu",
		"г": "нш", "u": "yi",
		"ш": "гщ", "i": "uo",
		"щ": "шз", "o": "ip",
		"з": "щх", "p": "o[",
		"х": "зъ",
		"ъ": "хэ",

		"ф": "ый", "a": "sq",
		"ы": "фв", "s": "ad",
		"в": "ыа", "d": "sf",
		"а": "вп", "f": "dg",
		"п": "ар", "g": "fh",
		"р": "по", "h": "gj",
		"о": "рл", "j": "hk",
		"л": "од", "k": "jl",
		"д": "лж", "l": "k;",
		"ж": "дэ",
		"э": "жх",

		"я": "чф", "z": "ax",
		"ч": "яс", "x": "zc",
		"с": "чм", "c": "xv",
		"м": "си", "v": "cb",
		"и": "мт", "b": "vn",
		"т": "иь", "n": "bm",
		"ь": "тб", "m": "n,",
		"б": "ью",
		"ю": "б.",
	}
	nearKeyboardLetters = combinations2
	combinations["0"] = "a"
}

func OneRandomError(inpWord string) string {
	inpWord = strings.ToLower(inpWord)
	wRunes := []rune(inpWord)
	if len(wRunes) <= 3 {
		return inpWord
	}
	rand.Seed(int64(randomSeed))
	randomSeed++
	indxToChange := rand.Intn(len(wRunes))
	indxToGet := rand.Intn(2)
	if _, ok := nearKeyboardLetters[string(wRunes[indxToChange])]; !ok {
		return inpWord
	}
	wrongLetter := []rune(nearKeyboardLetters[string(wRunes[indxToChange])])[indxToGet]
	wRunes[indxToChange] = wrongLetter
	return string(wRunes)
}

func TwoRandomError(inpWord string) string {
	inpWord = strings.ToLower(inpWord)
	wRunes := []rune(inpWord)
	if len(wRunes) < 5 {
		return OneRandomError(inpWord)
	}
	rand.Seed(int64(randomSeed))
	randomSeed++
	indxToChange1 := rand.Intn(len(wRunes))
	var indxToChange2 int
	for indxToChange2 = rand.Intn(len(wRunes)); indxToChange2 == indxToChange1; {
		indxToChange2 = rand.Intn(len(wRunes))
	}
	indxToGet1 := rand.Intn(2)
	indxToGet2 := rand.Intn(2)
	if _, ok := nearKeyboardLetters[string(wRunes[indxToChange1])]; !ok {
		return inpWord
	}
	if _, ok := nearKeyboardLetters[string(wRunes[indxToChange2])]; !ok {
		return inpWord
	}
	wrongLetter1 := []rune(nearKeyboardLetters[string(wRunes[indxToChange1])])[indxToGet1]
	wrongLetter2 := []rune(nearKeyboardLetters[string(wRunes[indxToChange2])])[indxToGet2]
	wRunes[indxToChange1] = wrongLetter1
	wRunes[indxToChange2] = wrongLetter2

	rand.Seed(time.Now().Unix())

	//
	

	return string(wRunes)
}

func OneErrorQuery(inpWord string) string {
	var out bytes.Buffer
	defer out.Reset()
	words := strings.Split(inpWord, " ")
	for i := range words {
		words[i] = OneRandomError(words[i]) //
	}
	return strings.Join(words, " ")
}

func TwoErrorQuery(inpWord string) string {
	words := strings.Split(inpWord, " ")
	for i := range words {
		words[i] = TwoRandomError(words[i])
	}
	return strings.Join(words, " ")
}

func GenerateOneErrorNTimes(inpWord string, num int) map[string][]string {
	out := make(map[string][]string)
	for i := 0; i < num; i++ {
		out[inpWord] = append(out[inpWord], OneErrorQuery(inpWord))
	}
	return out
}

func GenerateTwoErrorNTimes(inpWord string, num int) map[string][]string {
	out := make(map[string][]string)
	for i := 0; i < num; i++ {
		out[inpWord] = append(out[inpWord], TwoErrorQuery(inpWord))
	}
	return out
}

func NErrorPerEveryNWords(inpWord string, errorEveryNWords, NErrorsInWord, numTestCases int) map[string][]string {
	//min 1 error
	splt := strings.Fields(inpWord)
	out := make(map[string][]string)
	var nErrors int
	if errorEveryNWords == 0 || errorEveryNWords == 1 {
		nErrors = len(splt)
	} else {
		if len(splt) % errorEveryNWords != 0 {
			nErrors = len(splt) / errorEveryNWords + 1
		} else {
			nErrors = len(splt) / errorEveryNWords
		}
	}
	if NErrorsInWord >= 2 {
		NErrorsInWord = 2
	} else {
		NErrorsInWord = 1
	}
	for i := 0; i < numTestCases; i++ {
		errorQuery := strings.Fields(inpWord)
		errorIndxs := make(map[int]struct{})
		rand.Seed(int64(randomSeed))
		randomSeed++
		for len(errorIndxs) < nErrors {
			indx := rand.Intn(len(splt))
			errorIndxs[indx] = struct{}{}
		}
		for key := range errorIndxs {
			if NErrorsInWord == 1 {
				errorQuery[key] = OneRandomError(errorQuery[key])
			} else {
				errorQuery[key] = TwoRandomError(errorQuery[key])
			}
		}
		out[inpWord] = append(out[inpWord], strings.Join(errorQuery, " "))

	}
	return out
}

func TwoErrorPerEveryNWords(inpWord string, errorEveryNWords,num int) map[string][]string {
	//min 1 error
	splt := strings.Fields(inpWord)
	out := make(map[string][]string)
	var nErrors int
	if errorEveryNWords == 0 || errorEveryNWords == 1 {
		nErrors = len(splt)
	} else {
		nErrors = len(splt) / errorEveryNWords + 1
	}
	for i := 0; i < num; i++ {
		errorQuery := strings.Fields(inpWord)
		errorIndxs := make(map[int]struct{})
		rand.Seed(int64(randomSeed))
		randomSeed++
		for len(errorIndxs) < nErrors {
			indx := rand.Intn(len(splt))
			errorIndxs[indx] = struct{}{}
		}
		for key := range errorIndxs {
			errorQuery[key] = OneRandomError(errorQuery[key])
		}
		out[inpWord] = append(out[inpWord], strings.Join(errorQuery, " "))

	}
	return out
}
