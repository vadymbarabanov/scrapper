package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type job struct {
	id       string
	title    string
	location string
	salary   string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	totalPages := getPagesCount()

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(pageNumber int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(pageNumber*50)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".result").Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		title := card.Find(".jobTitle").Text()
		salary := card.Find(".salary-snippet").Text()
		summary := card.Find(".job-snippet").Text()
	})

}

func getPagesCount() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	pagesCount := 0

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pagesCount = s.Find("a").Length()
	})

	return pagesCount
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Response status code is", res.StatusCode)
	}
}
