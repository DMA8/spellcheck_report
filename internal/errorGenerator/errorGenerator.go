package errorGenerator

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

var keyboardBase, notCloseErrors, closeErrors map[string]string
var randomSeed int

const (
	PermitationCode   = -1
	MissingLetterCode = -2
	DoublingCode      = -3
)

//главная функция для генерации разнообразных ошибок
func Engine(baseWord string, errorType, ErrorFreqQuery, ErrorFreqWord, testCasesPerWord int) map[string][]string {
	switch errorType {
	case 0:
		keyboardBase = closeErrors
		return NErrorPerEveryNWords(baseWord, ErrorFreqQuery, ErrorFreqWord, testCasesPerWord)
	case 1:
		keyboardBase = notCloseErrors
		return NErrorPerEveryNWords(baseWord, ErrorFreqQuery, ErrorFreqWord, testCasesPerWord)
	case 2:
		return NErrorPerEveryNWords(baseWord, ErrorFreqQuery, PermitationCode, testCasesPerWord)
	case 3:
		return NErrorPerEveryNWords(baseWord, ErrorFreqQuery, MissingLetterCode, testCasesPerWord)
	case 4:
		return NErrorPerEveryNWords(baseWord, ErrorFreqQuery, DoublingCode, testCasesPerWord)
	case 5:
		ans := make(map[string][]string)
		keyboardBase = closeErrors
		ans[baseWord] = append(ans[baseWord], NErrorPerEveryNWords(baseWord, ErrorFreqQuery, ErrorFreqWord, testCasesPerWord)[baseWord]...)
		keyboardBase = notCloseErrors
		ans[baseWord] = append(ans[baseWord], NErrorPerEveryNWords(baseWord, ErrorFreqQuery, ErrorFreqWord, testCasesPerWord)[baseWord]...)
		ans[baseWord] = append(ans[baseWord], NErrorPerEveryNWords(baseWord, ErrorFreqQuery, PermitationCode, testCasesPerWord)[baseWord]...)
		ans[baseWord] = append(ans[baseWord], NErrorPerEveryNWords(baseWord, ErrorFreqQuery, MissingLetterCode, testCasesPerWord)[baseWord]...)
		ans[baseWord] = append(ans[baseWord], NErrorPerEveryNWords(baseWord, ErrorFreqQuery, DoublingCode, testCasesPerWord)[baseWord]...)
		return ans
	}
	return nil
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
	if _, ok := keyboardBase[string(wRunes[indxToChange])]; !ok {
		return inpWord
	}
	indxToGet := rand.Intn(len([]rune(keyboardBase[string(wRunes[indxToChange])]))) //changed here
	wrongLetter := []rune(keyboardBase[string(wRunes[indxToChange])])[indxToGet]
	wRunes[indxToChange] = wrongLetter
	return string(wRunes)
}

func OneRandomPermutation(inpWord string) string {
	inpWord = strings.ToLower(inpWord)
	wRunes := []rune(inpWord)
	if len(wRunes) <= 3 {
		return inpWord
	}
	rand.Seed(int64(randomSeed))
	randomSeed++
	indxToChange := rand.Intn(len(wRunes))
	if indxToChange < len(wRunes)-1 {
		wRunes[indxToChange], wRunes[indxToChange+1] = wRunes[indxToChange+1], wRunes[indxToChange]
	} else {
		wRunes[indxToChange], wRunes[indxToChange-1] = wRunes[indxToChange-1], wRunes[indxToChange]
	}
	return string(wRunes)
}

func OneRandomMissing(inpWord string) string {
	inpWord = strings.ToLower(inpWord)
	wRunes := []rune(inpWord)
	if len(wRunes) <= 3 {
		return inpWord
	}
	rand.Seed(int64(randomSeed))
	randomSeed++
	indxToChange := rand.Intn(len(wRunes))
	wRunes = append(wRunes[:indxToChange], wRunes[indxToChange+1:]...)
	return string(wRunes)
}

func OneRandomDoubling(inpWord string) string {
	inpWord = strings.ToLower(inpWord)
	wRunes := []rune(inpWord)
	if len(wRunes) <= 3 {
		return inpWord
	}
	rand.Seed(int64(randomSeed))
	randomSeed++
	indxToChange := rand.Intn(len(wRunes))
	wRunes = append(wRunes[:indxToChange+1], wRunes[indxToChange:]...)
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
	if _, ok := keyboardBase[string(wRunes[indxToChange1])]; !ok {
		return inpWord
	}
	if _, ok := keyboardBase[string(wRunes[indxToChange2])]; !ok {
		return inpWord
	}
	wrongLetter1 := []rune(keyboardBase[string(wRunes[indxToChange1])])[indxToGet1]
	wrongLetter2 := []rune(keyboardBase[string(wRunes[indxToChange2])])[indxToGet2]
	wRunes[indxToChange1] = wrongLetter1
	wRunes[indxToChange2] = wrongLetter2

	rand.Seed(time.Now().Unix())
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
		if len(splt)%errorEveryNWords != 0 {
			nErrors = len(splt)/errorEveryNWords + 1
		} else {
			nErrors = len(splt) / errorEveryNWords
		}
	}
	// if NErrorsInWord >= 2 {
	// 	NErrorsInWord = 2
	// } else {
	// 	NErrorsInWord = 1
	// }
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
			switch NErrorsInWord {
			case 1:
				errorQuery[key] = OneRandomError(errorQuery[key])
			case 2:
				errorQuery[key] = TwoRandomError(errorQuery[key])
			case PermitationCode:
				errorQuery[key] = OneRandomPermutation(errorQuery[key])
			case MissingLetterCode:
				errorQuery[key] = OneRandomMissing(errorQuery[key])
			case DoublingCode:
				errorQuery[key] = OneRandomDoubling(errorQuery[key])
			}

		}
		out[inpWord] = append(out[inpWord], strings.Join(errorQuery, " "))

	}
	return out
}

func TwoErrorPerEveryNWords(inpWord string, errorEveryNWords, num int) map[string][]string {
	//min 1 error
	splt := strings.Fields(inpWord)
	out := make(map[string][]string)
	var nErrors int
	if errorEveryNWords == 0 || errorEveryNWords == 1 {
		nErrors = len(splt)
	} else {
		nErrors = len(splt)/errorEveryNWords + 1
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

func init() {
	//все соседние буквы на клавиатуре
	notCloseErrors = map[string]string{
		"й": "яыу",
		"ц": "фывк",
		"у": "ыва",
		"к": "вап",
		"е": "апр",
		"н": "про",
		"г": "рол",
		"ш": "олд",
		"щ": "лдж",
		"з": "джэ",
		"х": "жэ\\",
		"ъ": "жэ\\",

		"ф": "ця",
		"ы": "йцуячс",
		"в": "цукчсм",
		"а": "укесми",
		"п": "кенмит",
		"р": "енгить",
		"о": "нгштьб",
		"л": "гшщьбю",
		"д": "шщзбю.",
		"ж": "щзхбю.",
		"э": "зхъ.",

		"я": "ыв",
		"ч": "фыв",
		"с": "ыва",
		"м": "вап",
		"и": "апр",
		"т": "про",
		"ь": "рол",
		"б": "лдж",
		"ю": "джэ",
	}

	//только соседние буквы на одной строке
	closeErrors = map[string]string{
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
}
