package main

import (
	"CexOutOfStock/LineMessagingSettings"
	"CexOutOfStock/crash"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/line/line-bot-sdk-go/linebot"
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName     = "./database.db"
	settingsFilePath = "./settings.json"
	//OneWeekInHours Hours there are in a week
	OneWeekInHours = 168
	sqliteFormat   = "2006-01-02 15:04:05"
)

var isInStock = false
var database *sql.DB

/*Need to figure out is it going to check forever and the rate at which it checks */
func main() {

	RunServer()
}

//RunServer opens the database and configustockString all server settings
func RunServer() {
	database = OpenDatabase()
	defer database.Close()
	var settings = LineMessagingSettings.GetLineMessagingSettings()
	go startUpdatePolling()
	http.HandleFunc("/", MainPage)
	http.HandleFunc("/line", LineWebHook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), nil))
	// log.Fatal(http.ListenAndServeTLS(
	// 	":"+strconv.Itoa(settings.Port),
	// 	settings.CertFile, settings.KeyFile, nil))

}

//MainPage fuction for http response
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. This is our first Go web app on Heroku!")
}

//LineWebHook fuction for http response main function for communicating with the line api
func LineWebHook(w http.ResponseWriter, r *http.Request) {

	var bot = LineMessagingSettings.GetLineMessagingSettings().Bot
	events, err := bot.ParseRequest(r)
	if err == nil {
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					var replyToken = event.ReplyToken
					var targetsID = event.Source.UserID
					messageToSendToUser := InsertEntryIntoDatabase(message.Text, targetsID)
					SendRepy(replyToken, messageToSendToUser, bot)
				}

			}
		}
	} else {
		fmt.Fprintf(w, "This page is the line hook")
	}
	w.Header().Set("Server", "A Go Web Server for line messaging")
	w.WriteHeader(200)
}

func SendRepy(replyToken, message string, bot *linebot.Client) {
	_, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do()
	if err != nil {
		fmt.Printf("error sending reply. Token: %v\n", replyToken)
	}
}

func InsertEntryIntoDatabase(URLOfProductPage, ID string) string {
	isInStock, isValidProductPage := GetStockInfoFromUrl(URLOfProductPage)
	var messageToOutput = ""
	if isValidProductPage == true {
		if isInStock == false {

			var prep, err = database.Prepare("INSERT OR IGNORE INTO users (userid) VALUES(?)")
			defer prep.Close()
			crash.PanicError(err)
			prep.Exec(ID)

			prep, err = database.Prepare("INSERT OR IGNORE INTO products(userid,url,lastupdated) VALUES(?, ?, ?)")
			defer prep.Close()
			crash.PanicError(err)
			prep.Exec(ID, URLOfProductPage, time.Now().UTC().Format(sqliteFormat))

			messageToOutput = "sorry not in stock but will alert you when it is :)"
		} else {
			messageToOutput = "This product is in stock"
		}

	} else {
		messageToOutput = "Sorry this isn't a valid CEX product page. url" + (URLOfProductPage)
	}
	return messageToOutput
}

func SendPushNotification(UserID, messageToSend string, bot *linebot.Client) {
	_, err := bot.PushMessage(UserID, linebot.NewTextMessage(messageToSend)).Do()
	crash.PanicError(err)
}

func startUpdatePolling() {
	for {
		time.Sleep(OneWeekInHours * time.Hour)
		var bot = LineMessagingSettings.GetLineMessagingSettings().Bot
		var productsToUpdate = GetProductsToUpdate()
		SendPushNotificationToMultipleUsers(productsToUpdate, bot)
	}
}

func CheckProductsToUpdate(productData map[string]string, bot *linebot.Client) {
	for userid, productURL := range productData {
		isInStock, _ := GetStockInfoFromUrl(productURL)
		if isInStock == true {
			go SendPushNotification(userid, "item "+productURL+" is in stock", bot)
			prep, err := database.Prepare("DELETE FROM products WHERE userid =? AND url = ?")
			defer prep.Close()
			crash.PanicError(err)
			prep.Exec(userid, productURL)
		} else {
			prep, err := database.Prepare("UPDATE products SET lastupdated = date('now') WHERE" +
				"userid = ? AND url = ?")
			defer prep.Close()
			crash.PanicError(err)
			prep.Exec(userid, productURL)

		}
	}
}

func GetProductsToUpdate() map[string]string {
	var productData = make(map[string]string)
	rows, err := database.Query("Select * from users left join products on products.userid=users.userid where lastupdated < date('now', '-7 days')  ")
	crash.PanicError(err)
	for rows.Next() {
		var userID = ""
		var productURL = ""

		rows.Scan(&userID, &productURL)
		productData[userID] = productURL

	}
	return productData
}

func SendPushNotificationToMultipleUsers(messageData map[string]string, bot *linebot.Client) {
	for userid, message := range messageData {
		go SendPushNotification(userid, message, bot)
	}
}

func GetStockInfoFromUrl(url string) (isInStock, isValidProductPage bool) {

	// https://github.com/chromedp/examples/blob/master/text/main.go
	isInStock, isValidProductPage = GetStockInfofFromReqestBody(url)

	// if err == nil {

	// } else {
	// 	isInStock = false
	// 	isValidProductPage = false
	// 	fmt.Printf("invalid URL %v serched\n", url)
	// }
	return isInStock, isValidProductPage
}

func GetStockInfofFromReqestBody(url string) (isInStock, isValidProductPage bool) {
	baseContext, baseCancel := context.WithDeadline(
		context.Background(), time.Now().Add(time.Second*10))
	ctx, cancel := chromedp.NewContext(baseContext)
	defer cancel()
	defer baseCancel()

	// run task list
	var stockString string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#__nuxt", chromedp.ByID),
		chromedp.Text(`div.buyNowButton`, &stockString, chromedp.BySearch),
	)
	if err == nil && stockString != "" {
		stockString = strings.TrimSpace(stockString)
		fmt.Println(stockString)
		stockString = strings.ToLower(stockString)
		switch stockString {
		case "out of stock":
			isInStock = false
			isValidProductPage = true
		case "i want to buy this item":
			isInStock = true
			isValidProductPage = true
		default:
			isValidProductPage = false
			fmt.Println("invalid url inputted ")

		}
	}
	return isInStock, isValidProductPage
}

func OpenDatabase() *sql.DB {
	var db *sql.DB
	if _, err := os.Stat(databaseName); os.IsNotExist(err) {
		os.Create(databaseName)
		db, err = sql.Open("sqlite3", databaseName)
		crash.PanicError(err)
		var schema, err = ioutil.ReadFile("./databaseschema.txt")
		crash.PanicError(err)

		db.Exec(string(schema))

	} else {
		db, err = sql.Open("sqlite3", databaseName)
		crash.PanicError(err)
	}

	return db

}

// func GetPort() string {
// 	var output = ""
// 	var defaultPortNumber = "8000"
// 	portAsString := os.Getenv("CEXBOTPORT")
// 	if portAsString == "" {
// 		fmt.Println("port enviroment variable not set setting server to listen to port {0}", defaultPortNumber)
// 		output = defaultPortNumber
// 	} else {
// 		output = portAsString
// 	}
// 	return ":" + output
// }

// func getSSLkeys() (certFile, keyFile string) {
// 	certFile = strings.TrimSpace(os.Getenv("certFile"))
// 	keyFile = strings.TrimSpace(os.Getenv("KeyFile"))
// 	if certFile == "" || keyFile == "" {
// 		log.Fatal("no ssl key or certificat given")
// 	}
// 	return certFile, keyFile
// }
