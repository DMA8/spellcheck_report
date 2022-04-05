package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"reportSpeller/internal/errorGenerator"
	"reportSpeller/internal/yandexspeller"

	"git.wildberries.ru/oer/tokenizer/normalize"
	"github.com/Saimunyz/speller" //спеллер
)

const (
	nTests           = 3000
	testCasesPerWord = 5
)

func countRightSuggest(right, suggest string) int {
	rightSplitted := strings.Split(right, " ")
	suggestSplitted := strings.Split(suggest, " ")
	rightC := 0
	if len(suggestSplitted) != len(rightSplitted) {
		return 0
	}
	for i := 0; i < len(rightSplitted); i++ {
		if rightSplitted[i] == suggestSplitted[i] {
			rightC++
		}
	}
	return rightC
}

type fullSentenceTestCounters struct {
	AllTested                   int
	SpellerRight                int
	SpellerNormalizedRight      int
	YaNormalizedRight           int
	YandexRight                 int
	YandexWrong                 int
	SpellerRightWhenYandexWrong int
}

type wordsTestCounters struct {
	allTested                         int
	spellerCorrected                  int
	yandexCorrected                   int
	spellerSuggestAnotherWordFreqDict int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func differnetRunes(errorWord, suggested string) int {
	counter := 0
	s1 := []rune(errorWord)
	s2 := []rune(suggested)
	minLen := min(len(s1), len(s2))
	diffLen := max(len(s1), len(s2)) - minLen
	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			counter++
		}
	}
	counter += diffLen
	return counter
}

func ya(inp string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	yandexSpellerClient := yandexspeller.New(
		yandexspeller.Config{
			Lang: "RU",
		},
		&http.Client{Timeout: time.Second * 20},
	)
	ans, err := yandexSpellerClient.SpellCheck(inp)
	if err != nil {
		ans, err = yandexSpellerClient.SpellCheck(inp)
		if err != nil {
			return
		}
	}
	out <- ans
}

func correctLearnData() {
	a := make(chan string, 100)
	file, err := os.Open("queriesRU.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file2, err := os.Create("cleanedSentences.txt")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	reader := bufio.NewScanner(file)
	//var mu sync.Mutex
	var wg sync.WaitGroup
	go func(){
		for {
			select {
			case msg := <- a:
				fmt.Fprintf(file2, "%s\n", msg)
			}
		}
	}()
	ok := true
	for ok{
		for i := 0; i < 100 && ok; i++ {
			ok = reader.Scan()
			wg.Add(1)
			go ya(reader.Text(), a, &wg)
		}
		wg.Wait()
	}
}

func main() {
	var mu sync.Mutex
	tokenizer := normalize.NewNormalizer()
	err := tokenizer.LoadDictionariesLocal("./data/words.csv.gz", "./data/spellcheck1.csv")
	if err != nil {
		log.Fatal(err)
	}

	sentenceCounter := fullSentenceTestCounters{}
	wordsCounter := wordsTestCounters{}
	done := make(chan struct{})
	set := make(map[string]struct{})
	freqMapFile, err := os.Open("datasets/freq.txt")
	freqMap := make(map[string]int)
	if err != nil {
		panic(err)
	}
	reader2 := bufio.NewScanner(freqMapFile)
	for ok := reader2.Scan(); ok; {
		splitted := strings.Split(reader2.Text(), " ")
		if len(splitted) != 2 {
			ok = reader2.Scan()
			continue
		}
		freq, _ := strconv.Atoi(splitted[1])
		freqMap[splitted[0]] = freq
		ok = reader2.Scan()
	}
	speller := speller.NewSpeller("config.yaml")
	yandexSpellerClient := yandexspeller.New(
		yandexspeller.Config{
			Lang: "RU",
		},
		&http.Client{Timeout: time.Second * 20},
	)

	// load model
	err = speller.LoadModel("models/model-without_singleWords.gz") 
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		done <- struct{}{}
		panic(err)
	}

	testFile, err := os.Open("without_singleWords.txt") // PROVIDE A PATH TO THE QUERIES
	if err != nil {
		panic(err)
	}
	spellerRightWhenYandexWrond, err := os.Create("spellerRightYandexWrong.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewScanner(testFile)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		select {
		case <-c:
			finish(&mu, sentenceCounter, wordsCounter, spellerRightWhenYandexWrond)
		case <-done:
			finish(&mu, sentenceCounter, wordsCounter, spellerRightWhenYandexWrond)
		}
	}()
	failLogs, err := os.Create("yaRigth_spellWrong_log.txt")
	if err != nil {
		panic(err)
	}
	bothWrong, err := os.Create("bothWrongLog.txt")
	if err != nil {
		panic(err)
	}
	normalizeSuccess, err := os.Create("normalizeSuccess.txt")
	if err != nil {
		panic(err)
	}
	normalizeFail, err := os.Create("normalizeFail.txt")
	if err != nil {
		panic(err)
	}

	normalizeSuccessYA, err := os.Create("normalizeSuccessYA.txt")
	if err != nil {
		panic(err)
	}
	normalizeFailYA, err := os.Create("normalizeFailYA.txt")
	if err != nil {
		panic(err)
	}

	for ok := reader.Scan(); ok; {
		set[strings.ToLower(reader.Text())] = struct{}{}
		ok = reader.Scan()
	}

	fmt.Fprint(failLogs, "(word -> error) | yaSucced: *word* | spellerFail: *spellerSuggest*\n\n")
	for msg, _ := range set {
		var flagLine, flagLine2, flagLine3, flagLine4, flagLine5, flagLine6 bool
		//msg := reader.Text()
		if len([]rune(msg)) < 3 {
			continue
		}
		set[msg] = struct{}{}
		msg = strings.Trim(msg, "\n")
		if !isCyrillic(msg) {
			continue
		}
		myErrors := errorGenerator.GenerateTwoErrorNTimes(msg, testCasesPerWord)

		mu.Lock()
		for RightWord, generatedErrors := range myErrors {
			spelRight, yaRigth := 0, 0
			fmt.Printf("Tested word is | %s |\n", RightWord)
			for _, generatedError := range generatedErrors {
				yandexResult, _ := yandexSpellerClient.SpellCheck(generatedError)
				spellerResult := speller.SpellCorrect(generatedError)
				wordsCounter.allTested += len(strings.Split(generatedError, " "))
				wordsCounter.spellerCorrected += countRightSuggest(RightWord, spellerResult)
				wordsCounter.yandexCorrected += countRightSuggest(RightWord, yandexResult)
				sentenceCounter.AllTested++
				if spellerResult == RightWord {
					sentenceCounter.SpellerRight++
					spelRight++
				} else {
					normalizedSpellerSuggest, normalizedRightWord := normalizeDiffWords(spellerResult, RightWord, tokenizer)

					if normalizedSpellerSuggest == normalizedRightWord {
						if !flagLine3 {
							fmt.Fprintf(normalizeSuccess, "Right: \"%s\" NormForm: \"%s\"|\n", RightWord, normalizedRightWord)
						}
						flagLine3 = true
						sentenceCounter.SpellerNormalizedRight++
						fmt.Fprintf(normalizeSuccess, "Speller: \"%s\" SpellerNormForm \"%s\" |(error: \"%s\")\n", spellerResult, normalizedSpellerSuggest, generatedError)
					} else {
						if !flagLine4 {
							fmt.Fprintf(normalizeFail, "Right: \"%s\" NormForm: \"%s\"|\n", RightWord, normalizedRightWord)
						}
						flagLine4 = true
						fmt.Fprintf(normalizeFail, "Speller: \"%s\" SpellerNormForm \"%s\" |(error: \"%s\")\n", spellerResult, normalizedSpellerSuggest, generatedError)
					}
				}
				if yandexResult == RightWord {
					sentenceCounter.YandexRight++
					yaRigth++
				} else {
					normYa, normalizedRightWord := normalizeDiffWords(yandexResult, RightWord, tokenizer)
					if normYa == normalizedRightWord {
						if !flagLine5 {
							fmt.Fprintf(normalizeSuccessYA, "Right: \"%s\" NormForm: \"%s\"|\n", RightWord, normalizedRightWord)
						}
						sentenceCounter.YaNormalizedRight++
						flagLine5 = true
						fmt.Fprintf(normalizeSuccessYA, "YandexSug: \"%s\" YaNorm: \"%s\" |(error: \"%s\")\n", yandexResult, normYa, generatedError)
					} else {
						if !flagLine6 {
							fmt.Fprintf(normalizeFailYA, "Right: \"%s\" NormForm: \"%s\"|\n", RightWord, normalizedRightWord)
						}
						flagLine6 = true
						fmt.Fprintf(normalizeFailYA, "YandexSug: \"%s\" YaNorm: \"%s\" |(error: \"%s\")\n", yandexResult, normYa, generatedError)
					}
					sentenceCounter.YandexWrong++
					if spellerResult == RightWord {
						flagLine = true
						sentenceCounter.SpellerRightWhenYandexWrong++
						fmt.Fprintf(spellerRightWhenYandexWrond, "W: %s E: %s Y: %s S: %s\n", RightWord, generatedError, yandexResult, spellerResult)
					}

				}
				if yandexResult == RightWord && spellerResult != RightWord {
					fmt.Fprintf(failLogs, "(%s -> %s) | yaSucced: %s | spellerFail: %s\n", RightWord, generatedError, yandexResult, spellerResult)
					ySplt := strings.Split(yandexResult, " ")
					sSplt := strings.Split(spellerResult, " ")
					rightSplt := strings.Split(RightWord, " ")
					errSplt := strings.Split(generatedError, " ")
					if len(ySplt) == len(sSplt) {
						for i := 0; i < len(ySplt); i++ {
							if ySplt[i] != sSplt[i] {
								fmt.Fprintf(failLogs, "Error: %s Expected: %s (freq: %d diffRunes: %d), SpellerSuggest: %s (freq: %d diffRunes: %d)\n",
									errSplt[i], rightSplt[i], freqMap[rightSplt[i]], differnetRunes(errSplt[i], rightSplt[i]),
									sSplt[i], freqMap[sSplt[i]], differnetRunes(errSplt[i], sSplt[i]))
								if freqMap[sSplt[i]] != 0 {
									wordsCounter.spellerSuggestAnotherWordFreqDict++
								}
							}
						}
					}
					fmt.Fprintf(failLogs, "------------------------------------------\n")
				}
				if yandexResult != RightWord && spellerResult != RightWord {
					flagLine2 = true
					fmt.Fprintf(bothWrong, "Error: %s Expected: %s SpellerSuggest: %s YandexSuggest: %s\n", generatedError, RightWord, spellerResult, yandexResult)
				}
				fmt.Printf("generated error is: %s; | S: %s %v |", generatedError, spellerResult, spellerResult == RightWord)
				fmt.Printf(" Y: %s %v |\n", yandexResult, yandexResult == RightWord)

			}
			fmt.Printf("spellerRight: %d, yaRight %d \n", spelRight, yaRigth)
			fmt.Println("------------------------------------------------------------------")
			if flagLine {
				fmt.Fprintf(spellerRightWhenYandexWrond, "-------------------------------------\n")
				flagLine = false
			}
			if flagLine2 {
				fmt.Fprintf(bothWrong, "-------------------------------------\n")
				flagLine2 = false
			}
			if flagLine3 {
				fmt.Fprintf(normalizeSuccess, "-------------------------------------\n")
				flagLine3 = false
			}
			if flagLine4 {
				fmt.Fprintf(normalizeFail, "-------------------------------------\n")
				flagLine4 = false
			}
			if flagLine5 {
				fmt.Fprintf(normalizeSuccessYA, "-------------------------------------\n")
				flagLine5 = false
			}
			if flagLine6 {
				fmt.Fprintf(normalizeFailYA, "-------------------------------------\n")
				flagLine6 = false
			}
			if sentenceCounter.AllTested > nTests {
				mu.Unlock()
				done <- struct{}{}
				time.Sleep(time.Second * 10)
			}
		}
		mu.Unlock()
	}
}

func finish(mut *sync.Mutex, c fullSentenceTestCounters, w wordsTestCounters, logFile *os.File) {
	mut.Lock()
	defer mut.Unlock()
	fmt.Printf("\nResults:\n TotalTests: %d\n SpellerRate %.2f%% (Norm: %.2f%%),  YandexRate %.2f%% (Norm: %.2f%%)\n",
		c.AllTested, float64(c.SpellerRight)/float64(c.AllTested)*100, float64(c.SpellerRight+c.SpellerNormalizedRight)/float64(c.AllTested)*100,
		float64(c.YandexRight)/float64(c.AllTested)*100, float64(c.YandexRight+c.YaNormalizedRight)/float64(c.AllTested)*100)
	fmt.Fprintf(logFile, "YandexFails %d SpellerRight %d SpellerRate %.2f\n", c.YandexWrong, c.SpellerRightWhenYandexWrong,
		(float64(c.SpellerRightWhenYandexWrong)/float64(c.YandexWrong))*100)
	fmt.Printf("Total words: %d, SpellerRate %.2f, YandexRate %.2f\n", w.allTested, (float64(w.spellerCorrected)/float64(w.allTested))*100, (float64(w.yandexCorrected)/float64(w.allTested))*100)
	os.Exit(0)
}

func normalizeDiffWords(suggest, right string, tk *normalize.Normalizer) (string, string) {
	var outputRight []string
	var outputSuggest []string

	splittedSuggest := strings.Split(suggest, " ")
	splittedRight := strings.Split(right, " ")
	if len(splittedRight) != len(splittedSuggest) {
		return suggest, right
	}
	for i := range splittedRight {
		if splittedRight[i] != splittedSuggest[i] {
			lemmaSuggest := tk.NormalizeWithoutMeta(splittedSuggest[i])[0][0].Lemma
			outputSuggest = append(outputSuggest, lemmaSuggest)
			lemmaRight := tk.NormalizeWithoutMeta(splittedRight[i])[0][0].Lemma
			outputRight = append(outputRight, lemmaRight)
		} else {
			outputRight = append(outputRight, splittedRight[i])
			outputSuggest = append(outputSuggest, splittedSuggest[i])
		}
	}
	suggestAns := strings.Join(outputSuggest, " ")
	rightAns := strings.Join(outputRight, " ")
	return suggestAns, rightAns
}

func getBrandRU(csvReader *csv.Reader) (string, error) {
	pattern := `^[а-яА-Я\s-]*$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}
	for {
		records, err := csvReader.Read()
		if err != nil || len(records) < 2 {
			log.Print(err)
			return "", err
		}
		brand := records[1]
		if r.MatchString(brand) {
			return strings.ToLower(brand), err
		}
	}
}

func getBrandEN(csvReader *csv.Reader) (string, error) {
	pattern := `^[a-zA-Z\s-]*$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}
	for {
		records, err := csvReader.Read()
		if err != nil || len(records) < 2 {
			log.Print(err)
			//records, err = csvReader.Read()
			return "", err
		}
		brand := records[1]
		if r.MatchString(brand) {
			return strings.ToLower(brand), nil
		}
	}
}

func isCyrillic(word string) bool {
	words := strings.Split(word, " ")
	for _, j := range words {
		for _, v := range j {
			if v >= 'а' && v <= 'я' {
				continue
			}
			return false
		}
	}
	return true
}
