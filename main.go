package main

import (
	"fmt"
	"net/http"
)

type urlResponse struct {
	url    string
	status string
}

func main() {
	c := make(chan urlResponse)
	results := map[string]string{}
	urls := []string{
		"https://www.google.com",
		"https://www.airbnb.com",
		"https://www.amazon.com",
		"https://www.reddit.com",
		"https://soundcloud.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}

	for i := 0; i < len(urls); i++ {
		res := <-c
		results[res.url] = res.status
	}

	for url, status := range results {
		fmt.Println(url, ":", status)
	}
}

func hitURL(url string, c chan<- urlResponse) {
	res, err := http.Get(url)

	if err != nil || res.StatusCode >= 400 {
		c <- urlResponse{status: "FAILED", url: url}
	} else {
		c <- urlResponse{status: "OK", url: url}
	}
}
