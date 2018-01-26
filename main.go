package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"gopkg.in/xmlpath.v2"
)

var isInStock = false
var channelAccessToken string

var APISecret string

func main() {
	http.HandleFunc("/", MainPage)
	fmt.Println("listening...")
	err := http.ListenAndServe(":4747", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//MainPage fuction for http response
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. This is our first Go web app on Heroku!")
}

//LineWebHook fuction for http response
func LineWebHook(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(channelAccessToken, APISecret)
	panicError(err)
	events, err := bot.ParseRequest(r)
	panicError(err)

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				var message = event.Message.(*linebot.TextMessage)
				fmt.Println(message.Text)
			}

		}
	}

}

func GetStockInfo(responseBody *io.ReadCloser) {
	root, err := xmlpath.ParseHTML(*responseBody)
	panicError(err)
	xpath := xmlpath.MustCompile("//div[@class = \"buyNowButton\"]")
	if stockString, ok := xpath.String(root); ok {
		stockString = strings.TrimSpace(stockString)
		stockString = strings.ToLower(stockString)
		fmt.Println(stockString)
		switch stockString {
		case "out of stock":
			fmt.Println("sorry :(")
		case "i want to buy this item":
			fmt.Println("yay in stock")
			// email user
		default:
			fmt.Println("invalid url inputed ")

		}
	}
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
