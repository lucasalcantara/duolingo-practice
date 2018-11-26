package duolingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	duolingoUrl          = "https://www.duolingo.com/users/"
	duolingoTranslateUrl = "https://d2.duolingo.com/api/1/dictionary/hints/%s/%s?tokens=[%s]"
)

type Duolingo struct {
	value map[string]*json.RawMessage
}

func New(username, lng string) (Duolingo, error) {
	var (
		d   Duolingo
		err error
	)

	d.value, err = getData(username, lng)
	return d, err
}

func getData(username, lng string) (map[string]*json.RawMessage, error) {
	var jsonObj map[string]*json.RawMessage

	url := duolingoUrl + username
	body, err := doRequest(url)
	if err != nil {
		return jsonObj, err
	}

	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		return jsonObj, errors.New("Error unmarshal json: " + err.Error())
	}

	err = json.Unmarshal(*jsonObj["language_data"], &jsonObj)
	if err != nil {
		return jsonObj, errors.New("Error unmarshal json: " + err.Error())
	}

	err = json.Unmarshal(*jsonObj[lng], &jsonObj)
	if err != nil {
		return jsonObj, errors.New("Error unmarshal json: " + err.Error())
	}

	return jsonObj, nil
}

func doRequest(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.New("Error http new request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		err = errors.New("Error doing request: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Error reading body: " + err.Error())
		return
	}

	return
}

func (d Duolingo) Words() []string {
	var (
		words  []string
		skills []*json.RawMessage
	)

	json.Unmarshal(*d.value["skills"], &skills)
	for _, skill := range skills {
		var value map[string]*json.RawMessage
		json.Unmarshal(*skill, &value)
		l, _ := strconv.ParseBool(string(*value["learned"]))
		if l {
			var learnedWords []string
			json.Unmarshal(*value["words"], &learnedWords)
			words = append(words, learnedWords...)
		}
	}

	return words
}

func (d Duolingo) TranslateWords(source, target string, words []string) (map[string][]string, error) {
	urlWords := ""
	var translated map[string][]string

	for _, w := range words {
		urlWords += fmt.Sprintf("\"%s\",", string(w))
	}

	url := fmt.Sprintf(duolingoTranslateUrl, source, target, urlWords[:len(urlWords)-1])
	body, err := doRequest(url)
	if err != nil {
		return translated, err
	}

	err = json.Unmarshal(body, &translated)
	if err != nil {
		return translated, errors.New("Error unmarshal json: " + err.Error())
	}

	return translated, nil
}
