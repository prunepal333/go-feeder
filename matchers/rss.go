package matchers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"example.com/search"
)

type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"PubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		GUID        string   `xml:"guid"`
		GeoRssPoint string   `xml:"georss:point"`
	}
	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}
	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Link           string   `xml:"link"`
		Description    string   `xml:"description"`
		PubDate        string   `xml:"pubdate"`
		LastBuildDate  string   `xml:"lastbuilddate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"imagingeditor"`
		WebMaster      string   `xml:"webmaster"`
		Image          string   `xml:"image"`
		Item           []item   `xml:"item"`
	}
	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)
type rssMatcher struct{}

func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}
func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	resp, err := http.Get(feed.URI)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
	}
	var document rssDocument
	//Decoding the response using xml decoder
	err = xml.NewDecoder(resp.Body).Decode(&document)

	return &document, err
}

func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	var results []*search.Result
	document, err := m.retrieve(feed)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	for _, channelItem := range document.Channel.Item {
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, &search.Result{
				Field:   "Title",
				Content: channelItem.Title,
				Link:    channelItem.Link,
			})
		}
		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, &search.Result{
				Field:   "Description",
				Content: channelItem.Description,
				Link:    channelItem.Link,
			})
		}
	}
	return results, nil
}
