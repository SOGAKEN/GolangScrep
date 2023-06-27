package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Keyword string
	Title   string
	URL     string
}

func Scrape(keyword string) []Result {
	// Google search URL
	googleURL := "http://www.google.com/search?q=" + keyword

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, _ := http.NewRequest("GET", googleURL, nil)

	// Set user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	// Execute HTTP request
	resp, _ := client.Do(req)

	// Create goquery document from HTTP response
	document, _ := goquery.NewDocumentFromReader(resp.Body)

	// Slice to store results
	results := []Result{}

	// Find search results
	sel := document.Find("div.yuRUbf")

	// For each item, extract title and URL
	for i := range sel.Nodes {
		item := sel.Eq(i)
		title := item.Find("h3").Text()
		link, _ := item.Find("a").Attr("href")
		link = strings.TrimPrefix(link, "/url?q=")

		if i < 5 { // only top 5
			results = append(results, Result{Keyword: keyword, Title: title, URL: link})
		} else {
			break
		}
	}

	return results
}

func main() {
	// Open the file
	csvfile, err := os.Open("keywords.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	// reading all rows at once
	keywordLines, err := reader.ReadAll()

	// Create output file
	outfile, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()

	writer := csv.NewWriter(outfile)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Keyword", "Title", "URL"})

	for _, line := range keywordLines {
		keyword := line[0]
		results := Scrape(keyword)

		for _, result := range results {
			writer.Write([]string{result.Keyword, result.Title, result.URL})
		}
	}
}
