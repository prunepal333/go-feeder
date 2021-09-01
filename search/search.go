package search

import (
	"log"
	"sync"
)

var matchers = make(map[string]Matcher)

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
func Run(searchTerm string) {

	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	results := make(chan *Result)

	var wg sync.WaitGroup

	//individual feeds are processed independently
	wg.Add(len(feeds))

	for _, feed := range feeds {
		matcher, exists := matchers[feed.Type]
		// fmt.Println(reflect.TypeOf(matcher))
		if !exists {
			matcher = matchers["default"]
		}
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			wg.Done()
		}(matcher, feed)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	Display(results)
}
