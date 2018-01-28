package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/xmlpath.v2"
)

const (
	databaseName = "./database.db"
)

var isInStock = false
var channelAccessToken string
var APISecret string
var database *sql.DB

/*Need to figure out is it going to check forever and the rate at which it checks */
func main() {

	database = OpenDatabase()
	// get api key and secret from io
	channelAccessToken = os.Getenv("channel")
	if channelAccessToken == "" {
		//log.Fatal("API key not given ")
	}
	APISecret = os.Getenv("secret")

	if APISecret == "" {
		//log.Fatal("API secret not given")
	}

	http.HandleFunc("/", MainPage)
	http.HandleFunc("/line", LineWebHook)
	var portNumber = GetPort()
	http.ListenAndServe(portNumber, nil)

}

//MainPage fuction for http response
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. This is our first Go web app on Heroku!")
}

//LineWebHook fuction for http response
func LineWebHook(w http.ResponseWriter, r *http.Request) {

	bot, err := linebot.New(APISecret, channelAccessToken)
	panicError(err)

	events, err := bot.ParseRequest(r)
	panicError(err)
	//fmt.Println(w, "hellow")
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var replyToken = event.ReplyToken
				var targetsID = event.Source.UserID
				InsertEntryIntoDatabase(message.Text, replyToken, targetsID, bot)
			}

		}
	}
}

func InsertEntryIntoDatabase(URLOfProductPage, replyToken, ID string, bot *linebot.Client) {
	isInStock, isValidProductPage := GetStockInfoFromUrl(URLOfProductPage)
	var messageToOutput = ""
	if isValidProductPage == true {
		if isInStock == false {
			database.Exec("INSERT OR IGNORE INTO users (userid) VALUES({0})", ID)
			database.Exec("INSERT OR IGNORE INTO products(userid,url,lastupdated)"+
				"VALUES({0}, {1}, {2})", ID, URLOfProductPage, time.Now().UTC())
			messageToOutput = "sorry not in stock but will alert you when it is :)"
		}
		messageToOutput = "This product is in stock"
	} else {
		messageToOutput = "Sorry this isn't a valid CEX product page."
	}

	bot.ReplyMessage(replyToken, linebot.NewTextMessage(messageToOutput))
}

func SendProductUpdates() {

}

func GetStockInfoFromUrl(url string) (isInStock, isValidProductPage bool) {
	responce, err := http.Get(url)
	if err != nil {
		fmt.Println("invalid URL {0} serched ", url)
	}
	defer responce.Body.Close()
	isInStock, isValidProductPage = GetStockInfofFromReqestBody(&responce.Body)
	return isInStock, isValidProductPage
}

func GetStockInfofFromReqestBody(responseBody *io.ReadCloser) (isInStock, isValidProductPage bool) {
	root, err := xmlpath.ParseHTML(*responseBody)
	panicError(err)
	xpath := xmlpath.MustCompile("//div[@class = \"buyNowButton\"]")
	if stockString, ok := xpath.String(root); ok {
		stockString = strings.TrimSpace(stockString)
		stockString = strings.ToLower(stockString)
		fmt.Println(stockString)
		switch stockString {
		case "out of stock":
			isInStock = false
			isValidProductPage = true
		case "i want to buy this item":
			isInStock = true
			isValidProductPage = true
		default:
			isValidProductPage = false
			fmt.Println("invalid url inputed ")

		}
	}
	return isInStock, isValidProductPage
}

func OpenDatabase() *sql.DB {
	var db *sql.DB
	if _, err := os.Stat(databaseName); os.IsNotExist(err) {
		os.Create(databaseName)
		db, err = sql.Open("sqlite3", databaseName)
		panicError(err)
		var schema, err = ioutil.ReadFile("./databaseschema.txt")
		panicError(err)

		db.Exec(string(schema))

	} else {
		db, err = sql.Open("sqlite3", databaseName)
		panicError(err)
	}

	return db

}

func GetPort() string {
	var output = ""
	var defaultPortNumber = "8000"
	portAsString := os.Getenv("PORT")
	if portAsString == "" {
		fmt.Println("port enviroment variable not set setting server to listen to port {0}", defaultPortNumber)
		output = defaultPortNumber
	} else {
		output = portAsString
	}
	return ":" + output
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
