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
	var TestURL = "https://uk.webuy.com/phones/product.php?mastersku=SAPPI8P64GGR&sku=SAPPI8P64GGRUNLB#.Wm675jfLdPY"
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
	var TestURL = "https://uk.webuy.com/product.php?sku=SMEM16G21331#.Wm95oDfLdPY"
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

func DeleteTestRecords() {
	database = OpenDatabase()
	defer database.Close()

	database.Exec("Delete from users where userid = TestID")
	database.Exec("Delete from products where userid = TestID")
}
