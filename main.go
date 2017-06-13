package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rjz/pdxdonuts/search"
	"googlemaps.github.io/maps"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"text/template"
)

// Edit as needed
var title = "City of Donuts"
var socialUrl = "https://rjz.github.io/pdxdonuts"
var socialTitle = "Portland, City of Donuts"
var socialImage = "donut.svg"
var googleAnalyticsId = "UA-100043557-1"

var apiKey = os.Getenv("GOOGLE_API_KEY")
var mapboxAccessToken = os.Getenv("MAPBOX_ACCESS_TOKEN")

var (
	optKeyword  = flag.String("keyword", "donut", "Keyword to search for")
	optType     = flag.String("type", "restaurant|bakery", "Types to search for (delimited|by|pipe")
	optLocation = flag.String("location", "Portland, OR", "Location")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func templatize(dir string, latLng maps.LatLng, data []byte) {
	pattern := filepath.Join(dir, "templates", "*.tmpl")
	t := template.Must(template.ParseGlob(pattern))
	t.Execute(os.Stdout, map[string]interface{}{
		"Title":             title,
		"Data":              string(data),
		"Lat":               latLng.Lat,
		"Lng":               latLng.Lng,
		"GoogleAnalyticsId": googleAnalyticsId,
		"MapboxAccessToken": mapboxAccessToken,
		"OpenGraphTags": map[string]interface{}{
			"Title": socialTitle,
			"Type":  "website",
			"URL":   socialUrl,
			"Image": fmt.Sprintf("%s/%s", socialUrl, socialImage),
		},
		"SocialLinks": map[string]interface{}{
			"Facebook": fmt.Sprintf("https://www.facebook.com/sharer/sharer.php?u=%s", url.QueryEscape(socialUrl)),
			"Twitter":  fmt.Sprintf("https://twitter.com/home?status=%s", url.QueryEscape(socialUrl+" "+socialTitle)),
		},
	})
}

func main() {
	flag.Parse()

	if apiKey == "" {
		usageAndExit("Please specify GOOGLE_API_KEY")
	} else if mapboxAccessToken == "" {
		usageAndExit("Please specify MAPBOX_ACCESS_TOKEN")
	} else if optLocation == nil {
		usageAndExit("-location is required")
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed creating client: %s", err)
	}

	s := search.NewSearch(c)
	if err := s.Do(*optLocation, &search.Options{
		Type:    *optType,
		Keyword: *optKeyword,
		Limit:   100,
		Radius:  10000, // m
	}); err != nil {
		log.Fatalf("Search failed '%s'", err)
	}

	log.Println("Serializing results...")
	serializedResults, err := json.Marshal(s.Results)
	if err != nil {
		log.Fatalf("failed serializing results: %s", err)
	}

	log.Println("Printing up your new map...")
	compactResults := new(bytes.Buffer)
	if err := json.Compact(compactResults, serializedResults); err != nil {
		log.Fatalf("failed compacting JSON results: %s", err)
	}

	log.Println("We're done! Find the goods in ./dist...")
	templatize(dir, s.LatLng, compactResults.Bytes())
}
