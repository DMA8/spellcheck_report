package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	freq := make(map[string]int)
	file, err := os.Open("cleanedSentences.txt")
	if err != nil {
		panic(err)
	}
	file2, err := os.Create("without_singleWords.txt")
	if err != nil {
		panic(err)
	}
	reader, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(reader), "\n")
	for _, v := range(lines) {
		wrds := strings.Split(v, " ")
		for _, wrd := range wrds {
			freq[wrd]++
		}
	}
	for _, v := range lines {
		flag := true
		spltWrds := strings.Split(v, " ")
		for _, v := range spltWrds{
			if freq[v] == 1 {
				flag = false
				break
			}
		}
		if flag {
			fmt.Fprintf(file2, "%s\n", v)
		}
	}
}

func in(line, word string) bool {
	wrds := strings.Split(line, " ")
	for _, v := range wrds {
		if v == word {
			return true
		}
	}
	return false
}