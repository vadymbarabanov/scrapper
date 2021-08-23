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
	totalPages := getPagesCount()

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, getPage(i)...)
	}

	writeJobs(jobs)
	fmt.Println("Done! Extracted:", len(jobs), "jobs")
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

func getPage(pageNumber int) []job {
	pageURL := baseURL + "&start=" + strconv.Itoa(pageNumber*50)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	result := []job{}
	doc.Find(".result").Each(func(i int, card *goquery.Selection) {
		result = append(result, extractJob(card))
	})

	return result
}

func extractJob(card *goquery.Selection) job {
	id, _ := card.Attr("data-jk")
	title := clearStr(card.Find(".jobTitle").Text())
	salary := clearStr(card.Find(".salary-snippet").Text())
	summary := clearStr(card.Find(".job-snippet").Text())
	companyLocation := clearStr(card.Find(".companyLocation").Text())
	companyName := clearStr(card.Find(".companyName").Text())

	return job{
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
