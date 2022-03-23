package main

import (
	"encoding/xml"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	WebhookURL = "https://379a-91-201-246-66.ngrok.io"
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
	"Dou":  "https://dou.ua/feed/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Println("please provide bot token")
		os.Exit(1)
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	//
	//_, err = bot.Request(tgbotapi.DeleteWebhookConfig{})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//wh, err := tgbotapi.NewWebhook(WebhookURL)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = bot.Request(wh)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//info, err := bot.GetWebhookInfo()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if info.LastErrorDate != 0 {
	//	log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	//}
	//
	//// новини які приходять в БОТ
	//updates := bot.ListenForWebhook("/")
	//go http.ListenAndServe(":8080", nil)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("%+v\n", update)
		if update.Message.Text == "Привіт" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		} else if url, ok := rss[update.Message.Text]; ok {
			rss, err := getNews(url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"sorry, error happened",
				))
			}
			for _, item := range rss.Items {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					item.URL+"\n"+item.Title,
				))
			}
		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"there is only available Habr",
			))
		}
	}
}
