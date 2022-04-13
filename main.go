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

// slowTest := []string{
// 	"детчкий кпем зажмвай ка от чсадин царааок и ущибов с масласи обоепихи сяты и шалыея мое солнышкл",
// 	"полртенце сахровое ьанное кухрнное для рук для ног для дица подарое сужчине пвпе мужц андоей",
// 	"каотина картмна на холсие еартина на холсье для игтерьера олееь геометпический арт х",
// 	"пюкзак женскиц для левочки для щколы для рабрты для офмса городсклй для прогулрк",
// 	"значрк кокаода на рилотку шппку краснач щвезда иеталл эсаль снрия совктская чимволика",
// 	"швабпа с отжииом и ведоом для мыття полоы оког сьен для уборкт ведрл со гваброй вкдро и шваюра",
// 	"швпбра с отдимом и ведррм для мвтья поллв оклн мтен для убррки вкдро со шааброй аедро и швабоа",
// 	"чай череый клубникп со сдивками гр чай чепный с жхинацеей и лмпой гр",
// 	"картинв по номепам ван гог во мне бабочкв на розоыом букетк х см холмт на подрамниуе",
// 	"картинп по номнрам живопмсь по номкрам кафк на берегц хллст на родрамнике х см",
// }

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
	for i := 0; i < len(spaceIndexes) - size + 1; i++ {
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
	if twoError{
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
	if !*NoError{
		if twoError{
			for _, v := range lines {
				testCases = append(testCases, errorGenerator.GenerateTwoErrorNTimes(v, testCasesPerWord))
			}
		} else {
			for _, v := range lines {
				testCases = append(testCases, errorGenerator.GenerateOneErrorNTimes(v, testCasesPerWord))
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
	for i := range slowest {
		fmt.Println(slowestQuery[i], slowest[i])
	}
	return testCounter, time.Since(start)
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
var TwoError = flag.Bool("e2", false, "Generate two errors for tests")
var NoError = flag.Bool("e0", false, "Don't generate errors")

var b = flag.Bool("b", false, "Bench mode")
var NWorkers = flag.Int("w", 0, "N workers for test. if 0 then syncroTest")
func main() {
	// var NWorkers *int
	// var TwoError, OneError *bool
	var mu sync.Mutex
	flag.Parse()

	fmt.Println(*TwoError, *NWorkers)

	log.Println("mem usage at launching")
	PrintMemUsage()
	tokenizer := normalize.NewNormalizer()
	err := tokenizer.LoadDictionariesLocal("./data/words.csv.gz", "./data/spellcheck1.csv") //Для токенайзера
	if err != nil {
		log.Fatal(err)
	}

	sentenceCounter := fullSentenceTestCounters{}
	wordsCounter := wordsTestCounters{}
	done := make(chan struct{})
	set := make(map[string]struct{})
	freqMapFile, err := os.Open("datasets/freq.txt") //FREQ лучше свежий закинуть
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
	speller1 := speller.NewSpeller("config.yaml")

	yandexSpellerClient := yandexspeller.New(
		yandexspeller.Config{
			Lang: "RU",
		},
		&http.Client{Timeout: time.Second * 20},
	)
	// speller2 := speller.NewSpeller("config.yaml")

	// nTest2, timeDur2 := benchmarkMulti(12, yandexSpellerClient.SpellCheck)
	// fmt.Println(nTest2, float64(nTest2)/float64(timeDur2.Milliseconds()))

	// nTest2, timeDur2 := benchmarkMulti(12,yandexSpellerClient.SpellCheck)
	// log.Println("mem usage when speller_1error test ends")
	// PrintMemUsage()
	// // // load model
	fmt.Println("mem usage before model loading")
	PrintMemUsage()
	// os.Exit(1)

	err = speller1.LoadModel("models/AllRu-model_tree.gz") //MODEL
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		done <- struct{}{}
		panic(err)
	}
	// err = speller2.LoadModel("models/AllRu-model.gz") //MODEL
	// if err != nil {
	// 	fmt.Printf("No such file: %v\n", err)
	// 	done <- struct{}{}
	// 	panic(err)
	// }

	speller1.SpellCorrect2("один и дваач и триич для четыре пятый и шестидесятый")
	speller1.SpellCorrect2("наклейка к мая это наша победа х см")


	if *b{
		if *NWorkers > 0 {
			nTest2, timeDur2 := benchmarkMulti(*NWorkers, *TwoError,speller1.SpellCorrect2)
			fmt.Println(nTest2, float64(nTest2)/float64(timeDur2.Milliseconds()))
			log.Println("mem usage when yandex test ends")
			PrintMemUsage()
			os.Exit(1)
		} else {
			nTest2, timeDur2 := benchmark(*TwoError, speller1.SpellCorrect2)
			fmt.Println(nTest2, float64(nTest2)/float64(timeDur2.Milliseconds()))
			log.Println("mem usage when yandex test ends")
			PrintMemUsage()
			os.Exit(1)
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


	yandexSpellerClient.SpellCheck("generatedError")
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
		var myErrors map[string][]string
		if *TwoError{
			myErrors = errorGenerator.GenerateTwoErrorNTimes(msg, testCasesPerWord)
		} else if *NoError {
			myErrors = make(map[string][]string)
			myErrors[msg] = []string{msg}
			nTests           = 9900
			testCasesPerWord = 1
		} else {
			myErrors = errorGenerator.GenerateOneErrorNTimes(msg, testCasesPerWord)
		}

		mu.Lock()
		for RightWord, generatedErrors := range myErrors {
			spelRight, yaRigth := 0, 0
			fmt.Printf("Tested word is | %s |\n", RightWord)
			for _, generatedError := range generatedErrors {
				yandexResult := ""
				// yandexResult := yandexSpellerClient.SpellCheck(generatedError)
				spellerResult := speller1.SpellCorrect2(generatedError)
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
	nErrors := 1
	if *TwoError {
		nErrors = 2
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
