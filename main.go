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

	"github.com/Saimunyz/speller" //спеллер
)

const (
	nTests           = 8000
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
	AllTested int
	SpellerRight int
	YandexRight int
	YandexWrong int
	SpellerRightWhenYandexWrong int
}

type wordsTestCounters struct {
	allTested int
	spellerCorrected int
	yandexCorrected int
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
	counter+=diffLen
	return counter
}


func main() {
	var mu sync.Mutex
	sentenceCounter := fullSentenceTestCounters{}
	wordsCounter := wordsTestCounters{}
	done := make(chan struct{})
	set := make(map[string]struct{})
	freqMapFile, err := os.Open("datasets/with_brands/brand-freq-dict.txt")
	freqMap := make(map[string]int)
	if err != nil {
		log.Fatal(err)
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
	err = speller.LoadModel("models/with_brand.gz")
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		done <- struct{}{}
		log.Println(err)
	}

	testFile, err := os.Open("sentences.txt") // PROVIDE A PATH TO THE QUERIES
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
	for ok := reader.Scan(); ok; {
		set[strings.ToLower(reader.Text())] = struct{}{}
		ok = reader.Scan()
	}

	fmt.Fprint(failLogs, "(word -> error) | yaSucced: *word* | spellerFail: *spellerSuggest*\n\n")
	for msg, _ := range set {
		var flagLine bool
		//msg := reader.Text()
		if  len([]rune(msg)) < 3 {
			continue
		}
		set[msg] = struct{}{}
		msg = strings.Trim(msg, "\n")
		if !isCyrillic(msg) {
			continue
		}
		//msg = speller.SpellCorrect(msg) //токенизация тестового слова
		myErrors := errorGenerator.GenerateTwoErrorNTimes(msg, testCasesPerWord)
		//не обработан случай, когда сгенерированная ошибка превращается в слово без орфографических ошибок.
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
				}
				if yandexResult == RightWord {
					sentenceCounter.YandexRight++
					yaRigth++
				} else {
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
				fmt.Printf("generated error is: %s; S: %s %v |", generatedError, spellerResult, spellerResult == RightWord)
				fmt.Printf(" Y: %s %v\n", yandexResult, yandexResult == RightWord)
			}
			fmt.Printf("spellerRight: %d, yaRight %d\n", spelRight, yaRigth)
			fmt.Println("------------------------------------------------------------------")
			if flagLine {
				fmt.Fprintf(spellerRightWhenYandexWrond, "-------------------------------------\n")
				flagLine = false
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
	fmt.Printf("\nResults:\n TotalTests: %d\n SpellerRate %.2f%%, YandexRate %.2f%%\n",
		c.AllTested, float64(c.SpellerRight)/float64(c.AllTested)*100,
		float64(c.YandexRight)/float64(c.AllTested)*100)
	fmt.Fprintf(logFile, "YandexFails %d SpellerRight %d SpellerRate %.2f\n", c.YandexWrong, c.SpellerRightWhenYandexWrong,
		(float64(c.SpellerRightWhenYandexWrong)/float64(c.YandexWrong))*100)
	fmt.Printf("Total words: %d, SpellerRate %.2f, YandexRate %.2f\n", w.allTested, (float64(w.spellerCorrected)/float64(w.allTested))*100, (float64(w.yandexCorrected)/float64(w.allTested))*100)
	os.Exit(0)
}

// func finishWhenFail(CounterAllTested, CounterSpellerRight, CounterYandexRight, CounterYandexWrong, CounterSpellerRightWhenYandexWrong int, logFile *os.File) {
// 	fmt.Printf("\nResults:\n TotalTests: %d\n SpellerRate %.2f%%, YandexRate %.2f%%\n",
// 		CounterAllTested, float64(CounterSpellerRight)/float64(CounterAllTested)*100,
// 		float64(CounterYandexRight)/float64(CounterAllTested)*100)

// 	fmt.Fprintf(logFile, "YandexFails %d SpellerRight %d SpellerRate %.2f\n", CounterYandexWrong, CounterSpellerRightWhenYandexWrong,
// 		(float64(CounterSpellerRightWhenYandexWrong)/float64(CounterYandexWrong))*100)
// 	os.Exit(0)
// }

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
