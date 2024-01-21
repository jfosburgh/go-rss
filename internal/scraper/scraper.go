package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

func FetchXML(url string) (rss, error) {
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
