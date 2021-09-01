package search

import (
	"fmt"
	"log"
)

type Result struct {
	Field   string
	Content string
	Link    string
}
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	searchResults, err := matcher.Search(feed, searchTerm)

	if err != nil {
		log.Println(err)
		return
	}
	for _, result := range searchResults {
		results <- result
	}
}
func Display(results chan *Result) {
	for result := range results {
		fmt.Printf("%s:\n%s\nFollow this link for more info %s\n", result.Field, result.Content, result.Link)
	}
}
