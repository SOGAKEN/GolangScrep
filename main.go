package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			// Open the file
			csvfile, err := os.Open(file.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer csvfile.Close()

			reader := csv.NewReader(csvfile)
			// reading all rows at once
			keywordLines, err := reader.ReadAll()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Create output file
			outfileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + "_results.csv"
			outfile, err := os.Create(outfileName)
			if err != nil {
				fmt.Println(err)
				continue
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
	}
}
