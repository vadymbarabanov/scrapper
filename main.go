package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type job struct {
	id              string
	title           string
	salary          string
	summary         string
	companyLocation string
	companyName     string
}

var prefixJobLink string = "https://kr.indeed.com/viewjob?jk="
var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	jobs := []job{}
	c := make(chan []job)
	totalPages := getPagesCount()

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done! Extracted:", len(jobs), "jobs")
}

func getPage(pageNumber int, mainC chan<- []job) {
	pageURL := baseURL + "&start=" + strconv.Itoa(pageNumber*50)
	fmt.Println("Requesting:", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	c := make(chan job)
	result := []job{}

	pages := doc.Find(".result")
	pages.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < pages.Length(); i++ {
		result = append(result, <-c)
	}

	mainC <- result
}

func writeJobs(jobs []job) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"JOB LINK", "TITLE", "SALARY", "SUMMARY", "COMPANY NAME", "COMPANY LOCATION"}
	checkErr(w.Write(headers))

	for _, job := range jobs {
		jobSlice := []string{
			prefixJobLink + job.id,
			job.title,
			job.salary,
			job.summary,
			job.companyName,
			job.companyLocation,
		}
		checkErr(w.Write(jobSlice))
	}
}

func extractJob(card *goquery.Selection, c chan<- job) {
	id, _ := card.Attr("data-jk")
	title := clearStr(card.Find(".jobTitle").Text())
	salary := clearStr(card.Find(".salary-snippet").Text())
	summary := clearStr(card.Find(".job-snippet").Text())
	companyLocation := clearStr(card.Find(".companyLocation").Text())
	companyName := clearStr(card.Find(".companyName").Text())

	c <- job{
		id:              id,
		title:           title,
		salary:          salary,
		summary:         summary,
		companyLocation: companyLocation,
		companyName:     companyName,
	}
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

func clearStr(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), "")
}
