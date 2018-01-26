package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"

	"gopkg.in/xmlpath.v2"
)

var isInStock = false
var channelAccessToken string

var APISecret string

/*Need to figure out is it going to check forever and the rate at which it checks */
func main() {
	// get api key and secret from io
	channelAccessToken = os.Getenv("channel")
	if channelAccessToken == "" {
		//log.Fatal("API key not given ")
	}
	APISecret = os.Getenv("secret")
	portAsString := ":" + os.Getenv("PORT")

	if APISecret == "" {
		//log.Fatal("API secret not given")
	}

	http.HandleFunc("/", MainPage)
	http.HandleFunc("/line", func(w http.ResponseWriter, req *http.Request) {
		bot, err := linebot.New(APISecret, channelAccessToken)
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	http.ListenAndServe(portAsString, nil)

}

//MainPage fuction for http response
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. This is our first Go web app on Heroku!")
}

var c = 0

//LineWebHook fuction for http response
func LineWebHook(w http.ResponseWriter, r *http.Request) {
	c++
	/*
		bot, err := linebot.New(APISecret, channelAccessToken)
		panicError(err)



			events, err := bot.ParseRequest(r)
			panicError(err)
			//fmt.Println(w, "hellow")
			for _, event := range events {
				if event.Type == linebot.EventTypeMessage {
					switch event.Message.(type) {
					case *linebot.TextMessage:
						var message = event.Message.(*linebot.TextMessage)
						fmt.Println(w, message.Text)
					}

				}
			}
	*/
	fmt.Fprintf(w, strconv.Itoa(c))
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
