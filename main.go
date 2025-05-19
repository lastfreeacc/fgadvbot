package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"saska.me/fgadvbot/fgaapi"
	"saska.me/fgadvbot/teleapi"
)

type cmd string

const (
	confFilename     = "fgadvbot.conf.json"
	startCmd     cmd = "/start"
	advCmd       cmd = "/adv"
	herCmd       cmd = "/her"
	codeCmd      cmd = "/code"
)

func (c cmd) isMe(msg string) bool {
	return strings.HasPrefix(msg, string(c))
}

var (
	conf     = make(map[string]interface{})
	botToken string
	bot      teleapi.Bot
)

func main() {
	myInit()
	upCh := bot.Listen()
	for update := range upCh {
		cmd := update.Message.Text
		switch true {
		case startCmd.isMe(cmd):
			doStrart(update)
		case advCmd.isMe(cmd):
			doAdv(update)
		case herCmd.isMe(cmd):
			doHer(update)
		case codeCmd.isMe(cmd):
			doCode(update)
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
	data, err := os.ReadFile(filename)
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
	/adv for not her
	/code leave it for me`)
	bot.SendMessage(update.Message.Chat.ID, msg)
}

const fgaAPIFail string = `Не могу дать то, что ты хочешь.
Попробуй зайти на оригинал: `

func doAdv(update *teleapi.Update) {
	adv, err := fgaapi.GetRandomAdvice()
	if err != nil {
		log.Printf("[Warning] can not get random advice: '%s'\n", err)
		err = bot.SendMessage(update.Message.Chat.ID, fgaAPIFail+"https://fucking-great-advice.ru")
		if err != nil {
			log.Printf("[Warning] some troubles with send, err: %s", err)
		}
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
		log.Printf("[Warning] can not get advice for her: '%s'\n", err)
		err = bot.SendMessage(update.Message.Chat.ID, fgaAPIFail+"https://fucking-great-advice.ru/advice/for-her")
		if err != nil {
			log.Printf("[Warning] some troubles with send, err: %s", err)
		}
		return
	}
	err = bot.SendMessage(update.Message.Chat.ID, adv.Text)
	if err != nil {
		log.Printf("[Warning] some troubles with send, err: %s", err)
	}
}

func doCode(update *teleapi.Update) {
	adv, err := fgaapi.GetRandomCoderAdvice()
	if err != nil {
		log.Printf("[Warning] can not get advice for coder: '%s'\n", err)
		err = bot.SendMessage(update.Message.Chat.ID, fgaAPIFail+"https://fucking-great-advice.ru/advice/coding")
		if err != nil {
			log.Printf("[Warning] some troubles with send, err: %s", err)
		}
		return
	}
	err = bot.SendMessage(update.Message.Chat.ID, adv.Text)
	if err != nil {
		log.Printf("[Warning] some troubles with send, err: %s", err)
	}
}
