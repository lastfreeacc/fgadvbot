package teleapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type method string

const (
	apiURL string = "https://api.telegram.org/bot"
 	sendMessageMthd method = "sendMessage"
 	getUpdates method = "getUpdates"
)

type bot struct {
	token        string
	updateCh     chan *Update
	currenOffset int64
}

func (bot *bot) makeURL(m method) string {
	return fmt.Sprintf("%s%s/%s", apiURL, bot.token, m)
}

// Bot ...
type Bot interface {
	SendMessage(int64, string) error
	Listen() <-chan *Update
}

// NewBot ...
func NewBot(t string) Bot {
	bot := bot{
		token:        t,
		updateCh:     make(chan *Update, 100),
		currenOffset: 0,
	}
	return &bot
}

// SendMessage ...
func (bot *bot) SendMessage(chatID int64, text string) error {
	jsonStr := fmt.Sprintf(`{"chat_id":"%d","text":"%s"}`, chatID, text)
	json := []byte(jsonStr)
	endPnt := bot.makeURL(sendMessageMthd)
	req, err := http.NewRequest("POST", endPnt, bytes.NewBuffer(json))
	if err != nil {
		log.Printf("[Error] in build req: %s", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Error] in send req: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Warning] can not read api answer: {method: %s, data:%s}, err: %s", sendMessageMthd, json, err)
	}
	return nil
}

func (bot *bot) Listen() <-chan *Update {
	go doUpdates(bot)
	return bot.updateCh
}

func doUpdates(bot *bot) {
	endPnt := bot.makeURL(getUpdates)
	for {
		jsonStr := fmt.Sprintf(`{"offset":%d, "timeout": 60}`, bot.currenOffset+1)
		jsonBlob := []byte(jsonStr)
		req, err := http.NewRequest("POST", endPnt, bytes.NewBuffer(jsonBlob))
		if err != nil {
			log.Printf("[Warning] can not getUpdates: %s", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Warning] in send req: %s", err.Error())
			continue
		}
		defer resp.Body.Close()
		respBlob, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Warning] can not read api answer: {method: %s, data:%s}, err: %s\n", getUpdates, jsonBlob, err)
		}
		var result getUpdatesResp
		err = json.Unmarshal(respBlob, &result)
		if err != nil {
			log.Printf("[Warning] can not unmarshal resp: %s\n", err.Error())
			continue
		}
		if !result.Ok {
			log.Printf("[Warning] result not ok\n")
			continue
		}
		for _, update := range result.Result {
			bot.updateCh <- update
			if update.UpdateID > bot.currenOffset {
				bot.currenOffset = update.UpdateID
			}
		}

	}
}