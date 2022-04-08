package yandexspeller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

)

// Config holds configuration of this component.
type Config struct {
	Address string
	Lang    string
}

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type yandexSpellerResponse []struct {
	Word      string   `json:"word"`
	Position  int      `json:"pos"`
	Len       int      `json:"len"`
	Corrected []string `json:"s"`
}

// YandexSpeller represents a wrapper for the Yandex speller's API.
type YandexSpeller struct {
	config     Config
	httpClient httpClient
}

// New returns a new YandexSpeller instance.
func New(
	 config Config, httpClient httpClient,
) *YandexSpeller {

	if config.Address == "" {
		config.Address = "https://speller.yandex.net/services/spellservice"
	}

	return &YandexSpeller{
		config:     config,
		httpClient: httpClient,
	}
}

func (ys *YandexSpeller) prepareURL(text string) (string, error) {
	address := ys.config.Address + ".json/checkText"

	u, err := url.Parse(address)
	if err != nil {
		return "",nil
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", nil
	}

	q.Add("lang", ys.config.Lang)
	q.Add("text", text)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// SpellCheck returns corrected string checked by Yandex's speller.
func (ys YandexSpeller) SpellCheck(query string) string {

	res, err := ys.makeRequestToYandexSpeller(query)
	if err != nil {
		return ""
	}

	corrected := ys.correctQuery(query, res)


	return corrected
}

func (ys YandexSpeller) makeRequestToYandexSpeller(
	query string,
) (yandexSpellerResponse, error) {
	url, err := ys.prepareURL(query)
	if err != nil {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil,nil
	}

	resp, err := ys.httpClient.Do(req)
	if err != nil {
		return nil, err
	}


	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}


	var res yandexSpellerResponse

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, nil
	}


	return res, nil
}

func (ys YandexSpeller) correctQuery(
	query string, res yandexSpellerResponse,
) string {
	if len(res) == 0 {
		return query
	}

	terms := strings.Split(query, " ")

	for i := 0; i < len(res); i++ {
		// extract a word
		word := query

		// sring from the left
		for j := 0; j < res[i].Position; j++ {
			_, size := utf8.DecodeRuneInString(word)

			word = word[size:]
		}

		// reduce on incoming Len
		// it's impossible to have a length more then a count of runes
		if utf8.RuneCountInString(word) < res[i].Len {
			res[i].Len = utf8.RuneCountInString(word)
		}

		// string from the right
		for utf8.RuneCountInString(word) != res[i].Len {
			_, size := utf8.DecodeLastRuneInString(word)

			word = word[:len(word)-size]
		}

		// find the word
		for j := 0; j < len(terms); j++ {
			if terms[j] == word {
				terms[j] = res[i].Corrected[0]
			}
		}
	}

	return strings.Join(terms, " ")
}
