package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/lastfreeacc/fgadvbot/fgaapi"
	"github.com/lastfreeacc/fgadvbot/teleapi"
)

type cmd string

const (
	confFilename     = "fgadvbot.conf.json"
	startCmd     cmd = "/start"
	advCmd       cmd = "/adv"
	herCmd       cmd = "/her"
)

func (c cmd) isMe(msg string) bool {
	return strings.HasPrefix(msg, string(c))
}

var (
	conf     = make(map[string]interface{})
	botToken string
	bot      teleapi.Bot
	// nextAdv  = "http://fucking-great-advice.ru/"
)

func main() {
	myInit()
	upCh := bot.Listen()
	for update := range upCh {
		switch true {
		case startCmd.isMe(update.Message.Text):
			doStrart(update)
		case advCmd.isMe(update.Message.Text):
			doAdv(update)
		case herCmd.isMe(update.Message.Text):
			doHer(update)
		// default:
		// 	bot.SendMessage(update.Message.Chat.ID, update.Message.Text)
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
	msg := fmt.Sprint(
		`Hello, i am an advice bot!
	Fucking Great Advice!
	Usage:
	/her for her
	/adv for not her`)
	bot.SendMessage(update.Message.Chat.ID, msg)
}

func doAdv(update *teleapi.Update) {
	adv, err := fgaapi.GetRandomAdvice()
	if err != nil {
		log.Printf("[Warning] can not get random advice: '%s'\n", err)
		return
	}
	err = bot.SendMessage(update.Message.Chat.ID, adv.Text)
	if err != nil {
		log.Printf("[Warning] some troubles with send, err: %s", err)
	}
}

func doHer(update *teleapi.Update) {
	adv, err := fgaapi.GetRandomHerAdvice()
	if err != nil {
		log.Printf("[Warning] can not get random her advice: '%s'\n", err)
		return
	}
	err = bot.SendMessage(update.Message.Chat.ID, adv.Text)
	if err != nil {
		log.Printf("[Warning] some troubles with send, err: %s", err)
	}
}

// func doAdv(update *teleapi.Update) {
// 	r, err := http.Get(nextAdv)
// 	if err != nil {
// 		log.Printf("[Warning] can not get advice, err: %s\n", err)
// 		return
// 	}
// 	if r.StatusCode >= 400 {
// 		log.Printf("[Warning] bad status: %d\n", r.StatusCode)
// 		return
// 	}
// 	body := r.Body
// 	if body == nil {
// 		log.Printf("[Warning] nil body: %s", body)
// 		return
// 	}
// 	defer body.Close()

// 	root, err := html.Parse(body)
// 	if err != nil {
// 		log.Printf("[Warning] can not parse, err: %s", err)
// 		return
// 	}
// 	next, err := parse.GetElementByID(root, "next")
// 	if err != nil {
// 		log.Printf("[Warning] can not find next, err: %s", err)
// 	}
// 	if next != nil {
// 		nextHref := parse.GetAttr(next, "href")
// 		if nextHref != "" {
// 			nextAdv = nextHref
// 		}
// 	}
// 	adv, err := parse.GetElementByID(root, "advice")
// 	if err != nil {
// 		log.Printf("[Warning] can not find advice, err: %s", err)
// 		return
// 	}
// 	msg := parse.GetTextFromTag(adv)

// 	err = bot.SendMessage(update.Message.Chat.ID, msg)
// 	if err != nil {
// 		log.Printf("[Warning] some troubles with send, err: %s", err)
// 	}
// }
