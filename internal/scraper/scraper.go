package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jfosburgh/go-rss/internal/database"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func fetchXML(url string) (rss, error) {
	r, err := http.Get(url)
	rssData := rss{}
	if err != nil {
		return rssData, fmt.Errorf("Error fetching %s: %v", url, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return rssData, fmt.Errorf("Status Error: %v", r.StatusCode)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return rssData, fmt.Errorf("Error Reading Body: %v", err)
	}

	err = xml.Unmarshal(data, &rssData)
	if err != nil {
		return rssData, fmt.Errorf("Error unmarshalling xml: %v", err)
	}

	return rssData, nil
}

func FetchFeeds(DB *database.Queries, interval, batchSize int32) {
	var wg sync.WaitGroup
	for true {
		feeds, _ := DB.GetNextFeedsToFetch(context.Background(), batchSize)
		for _, feed := range feeds {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				fmt.Printf("Fetching data from %s\n", url)
				rssData, _ := fetchXML(url)
				fmt.Printf("Found %d articles\n", len(rssData.Channel.Item))
			}(feed.Url)
		}

		fmt.Printf("Waiting for data to be fetched\n")
		wg.Wait()

		fmt.Printf("Sleeping for %d seconds\n", interval)
		time.Sleep(time.Second * time.Duration(interval))
	}
}
