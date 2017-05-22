package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {
	// Make a get request
	rs, err := http.Get("https://google.com")
	// Process response
	if err != nil {
		color.Red("WiFI OFF")
		//panic("Not connected to the net") // More idiomatic way would be to print the error and die unless it's a serious error

		// Learn about exit status in Golang
		os.Exit(1)
	}

	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)

	}

	bodyString := string(bodyBytes)

	println(bodyString)
}
