package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"reportSpeller/internal/errorGenerator"
	"reportSpeller/internal/yandexspeller"

	"git.wildberries.ru/oer/tokenizer/normalize"
	"github.com/Saimunyz/speller" //спеллер
)

var (
	nTests           = 10000
	testCasesPerWord = 5
)
var freqMap map[string]int

//FLAGS
var ShowMemory = flag.Bool("m", false, "Show memory usage")
var TestLogic = flag.Bool("logic", false, "Show memory usage")
var TwoError = flag.Bool("e2", false, "Generate two errors for tests")
var NoError = flag.Bool("e0", false, "Don't generate errors")
var SilentLogs = flag.Bool("s", false, "Dont show testCases")
var ShowSlow = flag.Bool("lags", false, "Show top 10 slowest queries and it's time")
var benchmarkMode = flag.Bool("b", false, "Benchmark mode")
var allBench_Quality = flag.Bool("all", false, "Show quality after benchmark")
var NWorkers = flag.Int("w", 0, "N workers for test. if 0 then syncroTest")
var ErrorsEveryNWords = flag.Int("errorFreq", 0, "How many errors generate per query (NWordsQuery / freqErr = totalErrors in query)")

func benchmark(twoError bool, speller func(string) string) (int, time.Duration) {
	var testCounter int
	testCases := make([]map[string][]string, 0, nTests)
	testFile, err := os.Open("CleanedUniqueRandomQueries.txt")
	if err != nil {
		panic(err)
	}
	txt, err := ioutil.ReadAll(testFile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(txt), "\n")
	log.Println("error generator is starting")
	if twoError {
		for _, v := range lines {
			testCases = append(testCases, errorGenerator.GenerateTwoErrorNTimes(v, testCasesPerWord))
		}
	} else {
		for _, v := range lines {
			testCases = append(testCases, errorGenerator.GenerateOneErrorNTimes(v, testCasesPerWord))
		}
	}
	log.Println("errors are generated", len(testCases)*testCasesPerWord)
	start := time.Now()
	log.Println("synchro test has been started")
	for _, test := range testCases {
		for _, errors := range test {
			for _, errorWord := range errors {
				testCounter++
				speller(errorWord)
			}
		}
	}
	log.Printf("tests %d Elapsed time %v", testCounter, time.Since(start))
	log.Printf("query/s %f query/ms %f", float64(testCounter)/float64(time.Since(start).Seconds()),
		float64(testCounter)/float64(time.Since(start).Milliseconds()))
	return testCounter, time.Since(start)
}

// func benchmarkMulti(nWorkers int, twoError bool, speller1 func(string, map[string]int) string) (int, time.Duration) {
func benchmarkMulti(nWorkers int, twoError bool, speller1 func(string) string) (int, time.Duration) {
	var testCounter int
	var wg sync.WaitGroup

	queue := make(chan string, nWorkers)
	testCases := make([]map[string][]string, 0, nTests)
	testFile, err := os.Open("CleanedUniqueRandomQueries.txt")
	if err != nil {
		panic(err)
	}
	defer testFile.Close()
	txt, err := ioutil.ReadAll(testFile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(txt), "\n")
	log.Println("eror generator is starting")
	if !*NoError {
		if twoError {
			for _, v := range lines {
				// testCases = append(testCases, errorGenerator.GenerateTwoErrorNTimes(v, testCasesPerWord))
				if v == "" {
					continue
				}
				testCases = append(testCases, errorGenerator.NErrorPerEveryNWords(v, *ErrorsEveryNWords, 2, testCasesPerWord))

			}
		} else {
			for _, v := range lines {
				testCases = append(testCases, errorGenerator.NErrorPerEveryNWords(v, *ErrorsEveryNWords, 1, testCasesPerWord))
			}
		}
		log.Println("errors are generated", len(testCases)*testCasesPerWord)
		go func() {
			for _, test := range testCases {
				for _, errors := range test {
					for _, errorWord := range errors {
						queue <- errorWord
						testCounter++
					}
				}
			}
			close(queue)
		}()
	} else {
		go func() {
			for _, test := range lines {
				queue <- test
				testCounter++
			}
			close(queue)
		}()
	}
	slowest := make([]time.Duration, 10)
	slowestQuery := make([]string, 10)
	for i := range slowest {
		slowest[i] = time.Nanosecond
	}
	start := time.Now()
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for msg := range queue {
				t1 := time.Now()
				speller1(msg)
				dur := time.Since(t1)
				for i := range slowest {
					if dur > slowest[i] {
						slowest[i] = dur
						slowestQuery[i] = msg
						break
					}
				}
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	log.Printf("workers: %d tests %d Elapsed time %v", nWorkers, testCounter, time.Since(start))
	log.Printf("query/s %f query/ms %f", float64(testCounter)/float64(time.Since(start).Seconds()),
		float64(testCounter)/float64(time.Since(start).Milliseconds()))
	if *ShowSlow {
		for i := range slowest {
			fmt.Println(slowestQuery[i], slowest[i])
		}
	}
	return testCounter, time.Since(start)
}

func main() {
	var mu sync.Mutex
	sentenceCounter := fullSentenceTestCounters{}
	wordsCounter := wordsTestCounters{}
	done := make(chan struct{})
	set := make(map[string]struct{})

	flag.Parse()
	if *ShowMemory {
		log.Println("mem usage at launching")
		PrintMemUsage()
	}
	tokenizer := normalize.NewNormalizer()
	err := tokenizer.LoadDictionariesLocal("./data/words.csv.gz", "./data/spellcheck1.csv") //Для токенайзера
	if err != nil {
		log.Fatal(err)
	}
	freqMapFile, err := os.Open("datasets/AllRu-freq-dict.txt") //FREQ лучше свежий закинуть
	freqMap = make(map[string]int)
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
	//filterWords("бэдВорд1 халат бэдворд2 халат бэдворд3 халат халат халат халат бэдворд4 бэдворд5")

	speller1 := speller.NewSpeller("config.yaml")

	yandexSpellerClient := yandexspeller.New(
		yandexspeller.Config{
			Lang: "RU",
		},
		&http.Client{Timeout: time.Second * 20},
	)
	yandexSpellerClient.SpellCheck("")
	if *ShowMemory {
		fmt.Println("mem usage before model loading")
		PrintMemUsage()
	}
	err = speller1.LoadModel("models/AllRu-model_new.gz") //MODEL
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		done <- struct{}{}
		panic(err)
	}
	speller1.SpellCorrect3("игрв карьочная для аечеринки для веседой компаеии еарт")
	if *benchmarkMode {
		if *NWorkers > 0 {
			nTest2, timeDur2 := benchmarkMulti(*NWorkers, *TwoError, speller1.SpellCorrect2) //Передача функции в бенчмарк
			fmt.Println(nTest2, float64(nTest2)/float64(timeDur2.Milliseconds()))
			if *ShowMemory {
				log.Println("mem usage when test ends")
				PrintMemUsage()
			}
			if !*allBench_Quality {
				os.Exit(1)
			}
		} else {
			nTest2, timeDur2 := benchmark(*TwoError, speller1.SpellCorrect2) //Передача функции в бенчмарк
			fmt.Println(nTest2, float64(nTest2)/float64(timeDur2.Milliseconds()))
			if *ShowMemory {
				log.Println("mem usage when test ends")
				PrintMemUsage()
			}
			if !*allBench_Quality {
				os.Exit(1)
			}
		}
	}
	testFile, err := os.Open("CleanedUniqueRandomQueries.txt") // QUERY
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

	//speller1.SpellCorrect2("сумуа жееская шопрер спортивнаф чернз плнчо на плечр дорржная теаневая", freqMap)
	fmt.Fprint(failLogs, "(word -> error) | yaSucced: *word* | spellerFail: *spellerSuggest*\n\n")
	for msg, _ := range set {
		var flagLine, flagLine2, flagLine3, flagLine4, flagLine5, flagLine6 bool
		if len([]rune(msg)) < 3 {
			continue
		}
		set[msg] = struct{}{}
		msg = strings.Trim(msg, "\n")
		if !isCyrillic(msg) {
			continue
		}
		var myErrors map[string][]string
		if *TwoError {
			myErrors = errorGenerator.NErrorPerEveryNWords(msg, *ErrorsEveryNWords, 2, testCasesPerWord)
		} else if *NoError {
			myErrors = make(map[string][]string)
			myErrors[msg] = []string{msg}
			nTests = 9900
			testCasesPerWord = 1
		} else {
			myErrors = errorGenerator.NErrorPerEveryNWords(msg, *ErrorsEveryNWords, 1, testCasesPerWord)
		}
		mu.Lock()
		for RightWord, generatedErrors := range myErrors {
			spelRight, yaRigth := 0, 0
			if !*SilentLogs {
				fmt.Printf("Tested word is | %s |\n", RightWord)
			}
			for _, generatedError := range generatedErrors {
				yandexResult := ""
				// yandexResult := yandexSpellerClient.SpellCheck(generatedError)
				// spellerResult := speller1.SpellCorrect3(generatedError, freqMap)
				spellerResult := speller1.SpellCorrect3(generatedError)
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
				if !*SilentLogs {
					fmt.Printf("generated error is: %s; | S: %s %v |", generatedError, spellerResult, spellerResult == RightWord)
					fmt.Printf(" Y: %s %v |\n", yandexResult, yandexResult == RightWord)
				}
			}
			if !*SilentLogs {
				fmt.Printf("spellerRight: %d, yaRight %d \n", spelRight, yaRigth)
				fmt.Println("------------------------------------------------------------------")
			}
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
				fmt.Println(*TestLogic)
				var fail int
				if *TestLogic{
					fillTestLogicCollection()
					for _, v := range testLogic {
						result := speller1.SpellCorrect3(v[1])
						// result := speller1.SpellCorrect2(v[1])
						if result == v[0] {
							continue
						} else {
							fail++
							fmt.Printf("%s -> %s (%s)\n", v[1],result ,v[0])
						} 
					}
					fmt.Printf("Logic errors %d/10 ", fail)
				}
				done <- struct{}{}
				time.Sleep(time.Second * 10)
			}
		}
		mu.Unlock()
	}
}
var testLogic [][]string
func fillTestLogicCollection() {
	testLogic = make([][]string, 0)
	testLogic = append(testLogic, []string{"томат дородный", "томат дорожный"})
	testLogic = append(testLogic, []string{"чемодан дорожный", "чемодан дородный"})
	testLogic = append(testLogic, []string{"женские толстовки", "жесткие толстовки"})
	testLogic = append(testLogic, []string{"костюм для тренировок", "костюм для тонировок"})
	testLogic = append(testLogic, []string{"летний плащ", "летний плач"})
	testLogic = append(testLogic, []string{"набор свечей столбик", "набор свечей столик"})
	testLogic = append(testLogic, []string{"жидкое мыло", "жирное мыло"})
	testLogic = append(testLogic, []string{"фотошторы для мальчика", "фотошторы для пальчика"})
	testLogic = append(testLogic, []string{"повседневное боди", "повседневное боги"})
	testLogic = append(testLogic, []string{"подушка розовая", "подушка разовая"})


}
func finish(mut *sync.Mutex, c fullSentenceTestCounters, w wordsTestCounters, logFile *os.File) {
	mut.Lock()
	defer mut.Unlock()
	nErrors := 1
	if *TwoError {
		nErrors = 2
	} else if *NoError {
		nErrors = 0
	}
	fmt.Printf("\nResults (nErrors %d):\n TotalTests: %d\n SpellerRate %.2f%% (Norm: %.2f%%),  YandexRate %.2f%% (Norm: %.2f%%)\n", nErrors,
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

func TokenNgramsNew(query string, size int) []string {
	spaceIndexes := make([]int, 0, size)
	for i, r := range query {
		if unicode.IsSpace(r) {
			spaceIndexes = append(spaceIndexes, i)
		}
	}
	spaceIndexes = append(spaceIndexes, len(query))
	outCap := len(spaceIndexes) - size + 1
	if outCap < 0 {
		outCap = 1
	}
	out := make([]string, 0, outCap)
	leftIndex := 0
	for i := 0; i < len(spaceIndexes)-size+1; i++ {
		step := i + size - 1
		if step >= len(spaceIndexes) {
			step = len(spaceIndexes) - 1
		}
		out = append(out, query[leftIndex:spaceIndexes[step]])
		leftIndex = spaceIndexes[i] + 1
	}
	return out
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}





/*
spellcorrect.go
func (o *SpellCorrector) CheckInFreqDict(query string) bool {
	return o.spell.CheckExistance(query)
}

spell.go
func (s *Spell) CheckExistance(input string) bool {
	lookupParams := s.defaultLookupParams()
	dict := lookupParams.dictOpts.name
	// Check for an exact match
	_, exists := s.library.load(dict, input)
	return exists
}

*/