package fgaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

type method string

const (
	apiURL                = "http://fucking-great-advice.ru/api/"
	randomMth      method = "random"
	randomHerMth   method = "random_by_tag/для%20нее"
	randomCoderMth method = "random_by_tag/кодеру"
)

var (
	// ErrBadStatus ...
	ErrBadStatus = errors.New("Bad status")
	// ErrEmptyBody ...
	ErrEmptyBody = errors.New("Empty body")
)

// Advice ..
type Advice struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Sound string `json:"sound"`
}

// GetRandomAdvice ...
func GetRandomAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s", apiURL, randomMth)
	return getAdvice(endPnt)
}

// GetRandomHerAdvice ...
func GetRandomHerAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s", apiURL, randomHerMth)
	return getAdvice(endPnt)
}

// GetRandomCoderAdvice ...
func GetRandomCoderAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s", apiURL, randomCoderMth)
	return getAdvice(endPnt)
}

func getAdvice(endPnt string) (*Advice, error) {
	r, err := http.Get(endPnt)
	if err != nil {
		log.Printf("[Warning] can not get random advice, err: %s\n", err)
		return nil, err
	}
	if r.StatusCode >= 400 {
		log.Printf("[Warning] bad status: %d\n", r.StatusCode)
		return nil, ErrBadStatus
	}
	body := r.Body
	if body == nil {
		log.Printf("[Warning] nil body\n")
		return nil, ErrEmptyBody
	}
	defer body.Close()
	advBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[Warning] some problems when read body, err: %s\n", err)
		return nil, err
	}
	var adv Advice
	err = json.Unmarshal(advBytes, &adv)
	if err != nil {
		log.Printf("[Warning] some problems when unmarshal json, err: %s\n", err)
		return nil, err
	}
	adv.Text = html.UnescapeString(adv.Text)
	return &adv, nil
}
