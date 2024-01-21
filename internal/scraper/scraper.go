package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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

var timeFormats = []string{
	time.RFC1123Z,
	time.ANSIC,
	time.RFC822,
	time.RFC850,
	time.Stamp,
	time.RFC1123,
	time.RFC3339,
	time.RFC822Z,
	time.UnixDate,
	time.RubyDate,
	time.RFC3339Nano,
}

func parseTime(unknownTime string) (time.Time, error) {
	var parsedTime time.Time
	var err error
	for _, layout := range timeFormats {
		parsedTime, err = time.Parse(layout, unknownTime)
		if err == nil {
			return parsedTime, err
		}
	}

	return parsedTime, err
}

func FetchFeeds(DB *database.Queries, interval, batchSize int32) {
	var wg sync.WaitGroup
	for true {
		feeds, _ := DB.GetNextFeedsToFetch(context.Background(), batchSize)
		for _, feed := range feeds {
			wg.Add(1)
			go func(url string, feedId uuid.UUID) {
				defer wg.Done()
				rssData, _ := fetchXML(url)
				for _, post := range rssData.Channel.Item {
					now := time.Now().UTC()
					description := post.Description
					descriptionValid := false
					if post.Description != "" {
						descriptionValid = true
					}
					publishedAt, err := parseTime(post.PubDate)
					if err != nil {
						fmt.Printf("Could not parse %s as ANSIC date\n", post.PubDate)
					}
					params := database.CreatePostParams{
						ID:          uuid.New(),
						CreatedAt:   now,
						UpdatedAt:   now,
						Title:       post.Title,
						Description: sql.NullString{String: description, Valid: descriptionValid},
						Url:         post.Link,
						PublishedAt: publishedAt.UTC(),
						FeedID:      feedId,
					}
					_, err = DB.CreatePost(context.Background(), params)
					if err != nil && err != sql.ErrNoRows {
						fmt.Printf("Error creating post: %v\n", err)
					}
				}
			}(feed.Url, feed.ID)
		}

		wg.Wait()

		time.Sleep(time.Second * time.Duration(interval))
	}
}
