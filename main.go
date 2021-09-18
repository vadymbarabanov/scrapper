package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vadymbarabanov/scrapper/scrapper"
)

const FILE_NAME string = "jobs.csv"

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleScrape(c echo.Context) error {
	defer os.Remove(FILE_NAME)
	term := strings.ToLower(scrapper.ClearStr(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(FILE_NAME, FILE_NAME)
}

func handleHome(c echo.Context) error {
	return c.File("home.html")
}
