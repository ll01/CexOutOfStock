package main

import (
	"os"
	"testing"
)

func TestOpeningDatabase(t *testing.T) {
	database = OpenDatabase()
	defer database.Close()
	if _, err := os.Stat(databaseName); os.IsNotExist(err) == true {
		t.Error("database file is not created at %v", databaseName)
	}
}

func TestRejectUsersProductURL(t *testing.T) {
	var TestURL = "www.google.com"
	var TestID = "TestID"
	var output = InsertEntryIntoDatabase(TestURL, TestID)

	if output != "Sorry this isn't a valid CEX product page." {
		t.Error("this url should be rejected ", TestURL)
	}
}

func TestRejectUsersProductNotValidURL(t *testing.T) {
	var TestURL = "Hello Bot !"
	var TestID = "TestID"
	var output = InsertEntryIntoDatabase(TestURL, TestID)

	if output != "Sorry this isn't a valid CEX product page." {
		t.Error("this url should be rejected ", TestURL)
	}
}

func TestProductInStock(t *testing.T) {
	DeleteTestRecords()
	var TestURL = "https://uk.webuy.com/product-detail?id=smem9qaeb&categoryName=memory---desktop-ddr3&superCatName=computing&title=8-gb-pc12800-ddr3-1600mhz-240-pin-memory"
	var TestID = "TestID"
	database = OpenDatabase()
	defer database.Close()

	var output = InsertEntryIntoDatabase(TestURL, TestID)
	rows, err := database.Query("SELECT userid FROM products where url = \"%v\"", TestURL)
	panicError(err)
	if output != "This product is in stock" && rows.Next() == true {
		t.Error("this url should be rejected ", TestURL)
	}
}

func TestProuctOutOfStock(t *testing.T) {
	DeleteTestRecords()
	var TestURL = "https://uk.webuy.com/product-detail?id=sgranvigtx650ti1gb&categoryName=graphics-cards-pci-e&superCatName=computing&title=nvidia-geforce-gtx-650-ti-1gb-dx11"
	var TestID = "TestID"
	database = OpenDatabase()
	defer database.Close()

	var output = InsertEntryIntoDatabase(TestURL, TestID)
	prep, err := database.Prepare("SELECT userid FROM products where url = ?")
	panicError(err)
	defer prep.Close()
	rows, err := prep.Query(TestURL)
	panicError(err)
	if output != "sorry not in stock but will alert you when it is :)" {
		t.Error("this url should be inserted %v", TestURL)
	}
	if rows.Next() {
		var username = ""
		rows.Scan(&username)
		if username != TestID {
			t.Error("this url should be inserted %v ", TestURL)
		}
	} else {
		t.Error("this url should be inserted %v ", TestURL)
	}
}

func TestNotificationOfItemBackInStock(t *testing.T) {
	var TestURL = "https://uk.webuy.com/product.php?sku=SMEM16G21331#.Wm95oDfLdPY"
	// dataMap  := make(map[string]string)
	DeleteTestRecords()
	database = OpenDatabase()
	defer database.Close()
	_, err := database.Exec("INSERT INTO users(userid) values TestID")
	panicError(err)
	_, err = database.Exec("INSERT INTO products(userid, url,lastupdated)" +
		"VALUES(TestID," + TestURL + ",date('now', '-1 month')")
	panicError(err)

	rows, err := database.Query("Select * from users left join products on products.userid=users.userid wherelastupdated < date('now', '-7 days')  ")
	panicError(err)
	for rows.Next() {
		//  dataMap[]
	}
	

}

func DeleteTestRecords() {
	 database = OpenDatabase()
	defer database.Close()

	database.Exec("Delete from users where userid = TestID")
	database.Exec("Delete from products where userid = TestID")
}
