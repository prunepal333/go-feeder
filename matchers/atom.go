package matchers

import (
	"encoding/xml"
	"log"
	"net/http"
	"regexp"

	"example.com/search"
)

type (
	link struct {
		XMLName xml.Name `xml:"link"`
		URIs    []string `xml:"href,attr"`
	}
	entry struct {
		XMLName xml.Name `xml:"entry"`
		Title   string   `xml:"title"`
		Link    link     `xml:"link"`
		ID      string   `xml:"id"`
		Updated string   `xml:"updated"`
		Summary string   `xml:"summary"`
		Content content  `xml:"content"`
		Author  author   `xml:"author"`
	}
	content struct {
		XMLName xml.Name `xml:"content"`
		Item    string   `xml:"div"`
	}
	author struct {
		XMLName xml.Name `xml:"author"`
		Name    string   `xml:"name"`
		Email   string   `xml:"email"`
	}
	atomFeed struct {
		XMLName  xml.Name `xml:"feed"`
		Entry    []entry  `xml:"entry"`
		Title    string   `xml:"title"`
		Subtitle string   `xml:"subtitle"`
		Link     link     `xml:"link"`
		ID       string   `xml:"id"`
		Updated  string   `xml:"updated"`
	}
)

type atomMatcher struct{}

func init() {
	var matcher atomMatcher
	search.Register("atom", matcher)
}
func (m atomMatcher) retrieve(feed *search.Feed) (*atomFeed, error) {

	var document atomFeed
	response, err := http.Get(feed.URI)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	err = xml.NewDecoder(response.Body).Decode(&document)
	if err != nil {
		return nil, err
	}

	return &document, err
}
func (m atomMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	document, err := m.retrieve(feed)
	var results []*search.Result
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	for _, feedEntry := range document.Entry {
		matched, err := regexp.MatchString(searchTerm, feedEntry.Title)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, &search.Result{
				Field:   "Title",
				Content: feedEntry.Title,
				Link:    feedEntry.Link.URIs[0],
			})
		}
		matched, err = regexp.MatchString(searchTerm, feedEntry.Summary)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, &search.Result{
				Field:   "Summary",
				Content: feedEntry.Summary,
				Link:    feedEntry.Link.URIs[0],
			})
		}
	}
	return results, err
}
