package errorGenerator

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var nearKeyboardLetters map[string]string
var randomSeed int

func init() {
	fmt.Println(1)
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
	if len(wRunes) < 3 {
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

	// indxToChange3 := rand.Intn(len(wRunes) - 1)
	// option := rand.Intn(4)
	// if option == 1 { //сдваиваем
	// 	r := wRunes[indxToChange3]
	// 	wRunes = wRunes[:len(wRunes) + 1]
	// 	copy(wRunes[indxToChange3 + 1 :], wRunes[indxToChange3:])
	// 	wRunes[indxToChange3] = r
	// 	return string(wRunes)
	// } else if option == 2 { //пропускаем символ
	// 	copy(wRunes[indxToChange3:], wRunes[indxToChange3 + 1:])
	// 	return string(wRunes[:len(wRunes) - 1])
	// }
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
