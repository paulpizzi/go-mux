package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var app App

func main() {
	app.Initialize("postgres", "", "postgres")

	args_print := flag.Bool("print", false, "set true to print products")
	args_insert := flag.Bool("insert", false, "set true to insert products")
	args_name := flag.String("name", "", "product name")
	args_price := flag.Float64("price", 0.0, "product price")

	flag.Parse()

	if *args_print && *args_insert {
		fmt.Println("Can only insert OR print!")
		return
	}
	if !*args_print && !*args_insert {
		fmt.Println("Must either insert (-insert) or print (-print)")
		return
	}

	if *args_print {
		printProducts()

	} else if *args_insert {
		var argsCorrect = true
		if *args_name == "" {
			fmt.Println("Must specify name with -name")
			argsCorrect = false
		}
		if *args_price <= 0.0 {
			fmt.Println("Must specify price with -price, which is >0")
			argsCorrect = false
		}

		if argsCorrect {
			insertProduct(*args_name, *args_price)
		}
	}
}

func insertProduct(name string, price float64) {
	str := `{"name":"` + name + `", "price": ` + fmt.Sprintf("%f", price) + `}`
	var jsonStr = []byte(str)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, req)

	if response.Code == http.StatusCreated {
		fmt.Println("Successfully inserted new Product: ", str)
	} else {
		fmt.Println("Error while inserting new Product. Error Code:", response.Code)
	}
}

func printProducts() {
	str := `{"start":"` + strconv.Itoa(0) + `", "count": ` + strconv.Itoa(10) + `}`
	var jsonStr = []byte(str)

	req, _ := http.NewRequest("GET", "/products", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, req)

	if response.Code == http.StatusOK {
		fmt.Println(response.Body.String())

	} else {
		fmt.Println("Error Code", response.Code)
	}
}
