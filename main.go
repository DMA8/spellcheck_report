package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"reportSpeller/internal/errorGenerator"

	"github.com/Saimunyz/speller" //спеллер
)

const (
	serviceURL = "http://speller.yandex.net/services/spellservice.json/checkText"
	nTests = 10000
	testCasesPerWord = 5
)

type Misspell struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
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

func testYA(testWord string) string {
	resp, err := http.PostForm(serviceURL, url.Values{
		"text":   {testWord},
		"lang":   {"ru"},
		"format": {"plain"},
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var misspells []Misspell
	if err = json.Unmarshal(body, &misspells); err != nil {
		panic(err)
	}

	if len(misspells) > 0 && len(misspells[0].Suggestions) > 0 {
		sugges := testWord
		for _, misspel := range misspells {
			if len(misspel.Suggestions) > 0 {
				sugges = strings.Replace(sugges, misspel.Word, misspel.Suggestions[0], 1)
			}
		}
		return sugges
	}
	return testWord
}

func finish(mut *sync.Mutex, CounterAllTested, CounterSpellerRight,  CounterYandexRight int) {
	mut.Lock()
	defer mut.Unlock()
	fmt.Printf("\nResults:\n TotalTests: %d\n SpellerRate %.2f%%, YandexRate %.2f%%\n", 
	CounterAllTested, float64(CounterSpellerRight) / float64(CounterAllTested) * 100, 
	float64(CounterYandexRight) / float64(CounterAllTested) * 100 )
	os.Exit(0)
}

func main() {
	var mu sync.Mutex
	var CounterAllTested, CounterSpellerRight, CounterYandexRight int
	set := make(map[string]struct{})
	speller := speller.NewSpeller("config.yaml")
	done := make(chan struct{})
	
	// load model
	err := speller.LoadModel("models/bin-data.gz")
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		panic(err)
	}

	testFile, err := os.Open("sentences.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewScanner(testFile)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		select {
		case <-c:
			finish(&mu, CounterAllTested, CounterSpellerRight, CounterYandexRight)
		case <- done:
			finish(&mu, CounterAllTested, CounterSpellerRight, CounterYandexRight)
		}
	}()
	failLogs, err := os.Create("yaRigth_spellWrong_log.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(failLogs, "(word -> error) | yaSucced: *word* | spellerFail: *spellerSuggest*\n\n")
	for ok := reader.Scan(); ok; {
		msg := reader.Text()
		if _, ok := set[msg]; ok || !isCyrillic(msg) || len(msg) < 3{
			ok = reader.Scan()
			continue
		}
		set[msg] = struct{}{}
		msg = strings.Trim(msg, "\n")
		msg = speller.SpellCorrect(msg) //токенизация тестового слова
		myErrors := errorGenerator.GenerateTwoErrorNTimes(msg, testCasesPerWord)
		//не обработан случай, когда сгенерированная ошибка превращается в слово без орфографических ошибок.
		mu.Lock()
		for RightWord, generatedErrors := range myErrors {
			spelRight, yaRigth := 0, 0
			fmt.Printf("Tested word is %s:\n", RightWord)
			for _, generatedError := range generatedErrors {
				yandexResult := testYA(generatedError)
				spellerResult := speller.SpellCorrect(generatedError)
				CounterAllTested++
				if spellerResult == RightWord {
					CounterSpellerRight++
					spelRight++
				}
				if yandexResult == RightWord {
					CounterYandexRight++
					yaRigth++
				}
				if (yandexResult == RightWord && spellerResult != RightWord) {
					fmt.Fprintf(failLogs, "(%s -> %s) | yaSucced: %s | spellerFail: %s\n", RightWord, generatedError, yandexResult, spellerResult)
				}
				fmt.Printf("generated error is: %s; S: %s %v |", generatedError, spellerResult, spellerResult == RightWord)
				fmt.Printf(" Y: %s %v\n", yandexResult, yandexResult == RightWord)
			}
			fmt.Printf("spellerRight: %d, yaRight %d\n", spelRight, yaRigth)
			fmt.Println("------------------------------------------------------------------")
			if CounterAllTested > nTests {
				done <- struct{}{}
			}
		}
		mu.Unlock()
		ok = reader.Scan()
	}
}
