package fgaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type method string

const (
	apiURL              = "http://fucking-great-advice.ru/api/"
	randomMth    method = "random"
	apiV2URL            = "http://fucking-great-advice.ru/api/v2/"
	randomTagMth method = "random-advices-by-tag?tag="
	forHerTag           = "for-her"
	codingTag           = "coding"
)

var (
	// ErrBadStatus ...
	ErrBadStatus = errors.New("Bad status")
	// ErrEmptyBody ...
	ErrEmptyBody = errors.New("Empty body")
	// ErrAPIError ...
	ErrAPIError = errors.New("Api error")
)

// GetRandomAdvice ...
func GetRandomAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s", apiURL, randomMth)

	return getAdvice(endPnt)
}

// GetRandomHerAdvice ...
func GetRandomHerAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s%s", apiV2URL, randomTagMth, forHerTag)

	return getAdviceV2(endPnt)
}

// GetRandomCoderAdvice ...
func GetRandomCoderAdvice() (*Advice, error) {
	endPnt := fmt.Sprintf("%s%s%s", apiV2URL, randomTagMth, codingTag)

	return getAdviceV2(endPnt)
}

func getAdvice(endPnt string) (*Advice, error) {
	body, err := getBody(endPnt)
	if err != nil {
		// all errors was processed in getBody
		return nil, err
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

func getAdviceV2(endPnt string) (*Advice, error) {
	body, err := getBody(endPnt)
	if err != nil {
		// all errors was processed in getBody
		return nil, err
	}

	defer body.Close()

	advBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[Warning] some problems when read body, err: %s\n", err)

		return nil, err
	}
	var advResp AdviceResponse
	err = json.Unmarshal(advBytes, &advResp)
	if err != nil {
		log.Printf("[Warning] some problems when unmarshal json, err: %s\n", err)

		return nil, err
	}
	if len(advResp.Errors) != 0 {
		log.Printf("[Warning] request has errors: %v\n", advResp.Errors)

		return nil, ErrAPIError
	}
	if len(advResp.Data) == 0 {
		log.Printf("[Warning] request has no data\n")

		return nil, ErrAPIError
	}
	adv := Advice{}
	adv.Text = html.UnescapeString(advResp.Data[0].Text)

	return &adv, nil
}

func getBody(endPnt string) (io.ReadCloser, error) {
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

	return body, nil
}
