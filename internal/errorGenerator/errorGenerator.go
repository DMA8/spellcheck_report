package errorGenerator

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
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
		"й": "цф",
		"ц": "йу",
		"у": "цк",
		"к": "уе",
		"е": "кн",
		"н": "ег",
		"г": "нш",
		"ш": "гщ",
		"щ": "шз",
		"з": "щх",
		"х": "зъ",
		"ъ": "хэ",

		"ф": "ый",
		"ы": "фв",
		"в": "ыа",
		"а": "вп",
		"п": "ар",
		"р": "по",
		"о": "рл",
		"л": "од",
		"д": "лж",
		"ж": "дэ",
		"э": "жх",

		"я": "чф",
		"ч": "яс",
		"с": "чм",
		"м": "си",
		"и": "мт",
		"т": "иь",
		"ь": "тб",
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
	return string(wRunes)
}

func OneErrorQuery(inpWord string) string {
	var out bytes.Buffer
	defer out.Reset()
	words := strings.Split(inpWord, " ")
	for i := range words {
		words[i] = TwoRandomError(words[i])
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
