package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/net/html"
	"github.com/lastfreeacc/fgadvbot/teleapi"
	"github.com/lastfreeacc/fgadvbot/parse"
	"fmt"
)

type cmd string

const (
	confFilename = "fgadvbot.conf.json"
	startCmd  = "/start"
	advCmd  = "/adv"
	herCmd  = "/her"
)

var (
	conf = make(map[string]interface{})
	botToken string
	bot teleapi.Bot
)

func main() {
	myInit()
	upCh := bot.Listen()
	for update := range upCh {
		switch update.Message.Text {
		case startCmd:
			doStrart(update)
		case advCmd:
			doAdv(update)
		default:
			bot.SendMessage(update.Message.Chat.ID, update.Message.Text)
		} 
	}
}

func myInit() {
	readMapFromJSON(confFilename, &conf)
	botToken, ok := conf["botToken"]
	if !ok || botToken == "" {
		log.Fatalf("[Error] can not find botToken in config file: %s\n", confFilename)
	}
	bot = teleapi.NewBot(botToken.(string))
}

func readMapFromJSON(filename string, mapVar *map[string]interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("[Warning] can not read file '%s'\n", filename)
	}
	if err := json.Unmarshal(data, mapVar); err != nil {
		log.Fatalf("[Warning] can not unmarshal json from file '%s'\n", filename)
	}
	log.Printf("[Info] read data from file: %s:\n%v\n", filename, mapVar)
}

func doStrart(update *teleapi.Update) {
	msg := fmt.Sprint("Hello, usage: \n/her for her\n/adv for not her")
	bot.SendMessage(update.Message.Chat.ID, msg)
}

func doAdv(update *teleapi.Update) {
	r, err := http.Get("http://fucking-great-advice.ru/")
	if err != nil {
		log.Printf("[Warning] can not get advice, err: %s\n", err)
		return
	}
	if r.StatusCode >= 400 {
		log.Printf("[Warning] bad status: %d\n", r.StatusCode)
		return
	}
	body := r.Body
	if body == nil {
		log.Printf("[Warning] nil body: %s", body)
		return
	}
	defer body.Close()

	root, err := html.Parse(body)
	if err != nil {
		log.Printf("[Warning] can not parse, err: %s", err)
		return
	}
	adv, err := parse.GetElementByID(root, "advice")
	if err != nil {
		log.Printf("[Warning] can not find advice, err: %s", err)
		return
	}
	// msg := adv.Data
	bot.SendMessage(update.Message.Chat.ID, adv.Data)
}