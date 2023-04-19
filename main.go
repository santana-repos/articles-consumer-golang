package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var MAX_PAGE int

func revieveAPIjsonData(page int) map[string]interface{} {
	var result map[string]interface{}

	// Make an HTTP GET request to the API endpoint
	resp, err := http.Get(fmt.Sprintf(`%v%v`, "https://jsonmock.hackerrank.com/api/articles?page=", page))
	if err != nil {
		log.Printf("Failed to fetch articles from API:: Error:\n%v", err)
		return result
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var articlesData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&articlesData)

	if err != nil {
		log.Printf("Failed to decode articles:: Error:\n%v", err)
		return result
	}

	return articlesData
}

func retrieveMaxPage() int {
	if MAX_PAGE == 0 {
		articlesData := revieveAPIjsonData(1)
		totalPages, ok := articlesData["total_pages"]
		if !ok {
			MAX_PAGE = 1
		} else {
			MAX_PAGE = int(totalPages.(float64))
		}
	}
	return MAX_PAGE
}

func TopArticlesPerPage(page int) []string {
	result := []string{}

	articlesData := revieveAPIjsonData(page)

	// Extract the articles data from the response
	data, ok := articlesData["data"].([]interface{})
	if !ok {
		log.Printf("Invalid articles data")
		return result
	}

	// Convert the string data to the response payload format
	var articlesResults []string

	for _, articles := range data {
		if article, ok := articles.(map[string]interface{}); ok {
			if articleTitle := article["title"]; (articleTitle != nil) && (len(strings.TrimSpace(articleTitle.(string))) > 0) {

				articlesResults = append(articlesResults, articleTitle.(string))
			} else if articleTitleFallback := article["story_title"]; (articleTitleFallback != nil) && (len(strings.TrimSpace(articleTitleFallback.(string))) > 0) {

				articlesResults = append(articlesResults, articleTitleFallback.(string))
			}
		}
	}

	return append(result, articlesResults...)
}

func TopArticles(limit int) []string {
	result := []string{}

	if limit < 1 {
		return result
	}

	maxPage := retrieveMaxPage()
	var pagesToRetrieve int
	if limit > maxPage {
		pagesToRetrieve = maxPage
	} else {
		pagesToRetrieve = limit
	}

	for page := 1; page <= pagesToRetrieve; page++ {
		articles := TopArticlesPerPage(page)
		result = append(result, articles...)
	}

	return result
}

func main() {
	log.Println("*************************************************")
	var articleList = TopArticles(0)

	if len(articleList) < 1 {
		log.Println("No articles was retrieved!")
	}

	for index, article := range articleList {
		log.Printf("Article[%v] :: [%v]\n", index, article)
	}

	log.Println("*************************************************")

	articleList = TopArticles(1)
	for index, article := range articleList {
		log.Printf("Article[%v] :: [%v]\n", (index + 1), article)
	}

	log.Println("*************************************************")

	articleList = TopArticles(100)
	for index, article := range articleList {
		log.Printf("Article[%v] :: [%v]\n", index, article)
	}

	log.Println("*************************************************")
}
