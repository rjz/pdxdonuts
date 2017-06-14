package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
)

type OpenGraphData struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
	Image string `json:"image"`
}

type SocialData struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

type PageData struct {
	Title             string        `json:"title"`
	Data              string        `json:"data"`
	Icon              string        `json:"icon"`
	URL               string        `json:"url"`
	Lat               float64       `json:"lat"`
	Lng               float64       `json:"lng"`
	GoogleAnalyticsId string        `json:"googleAnalyticsId"`
	MapboxAccessToken string        `json:"mapboxAccessToken"`
	OpenGraphTags     OpenGraphData `json:"openGraph"`
	SocialLinks       SocialData    `json:"social"`
}

func Apply(dir string, data *PageData) {
	pattern := filepath.Join(dir, "templates", "*.tmpl")
	t := template.Must(template.ParseGlob(pattern))
	t.Execute(os.Stdout, data)
}

func LoadPageData(filename string) (*PageData, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var pd PageData
	if err := json.Unmarshal(data, &pd); err != nil {
		return nil, err
	}

	if pd.SocialLinks.Facebook == "" {
		pd.SocialLinks.Facebook = fmt.Sprintf("https://www.facebook.com/sharer/sharer.php?u=%s", url.QueryEscape(pd.URL))
	}

	if pd.SocialLinks.Twitter == "" {
		pd.SocialLinks.Twitter = fmt.Sprintf("https://twitter.com/home?status=%s", url.QueryEscape(pd.URL+" "+pd.Title))
	}

	if pd.OpenGraphTags.Title == "" {
		pd.OpenGraphTags.Title = pd.Title
	}

	if pd.OpenGraphTags.Image == "" {
		pd.OpenGraphTags.Image = fmt.Sprintf("%s/%s", pd.URL, pd.Icon)
	}

	if pd.OpenGraphTags.Type == "" {
		pd.OpenGraphTags.Type = "website"
	}

	if pd.OpenGraphTags.URL == "" {
		pd.OpenGraphTags.URL = pd.URL
	}

	return &pd, nil
}
