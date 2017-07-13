package fgaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiURL              = "http://fucking-great-advice.ru/api/"
	randomMth    method = "random"
	randomHerMth method = "random_by_tag/для%20нее"
)

var (
	// ErrBadStatus ...
	ErrBadStatus = errors.New("Bad status")
	//
	ErrEmptyBody = errors.New("Empty body")
)

type method string

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
	adv := &Advice{}
	err = json.Unmarshal(advBytes, adv)
	if err != nil {
		log.Printf("[Warning] some problems when unmarshal json, err: %s\n", err)
		return nil, err
	}
	return adv, nil
}