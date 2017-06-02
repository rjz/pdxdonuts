package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"text/template"
	"time"
)

// Edit as needed
var title = "City of Donuts"
var socialUrl = "https://rjz.github.io/pdxdonuts"
var socialTitle = "Portland, City of Donuts"
var socialImage = "donut.svg"
var googleAnalyticsId = "UA-100043557-1"

// A datum corresponds to one entry on the map
type datum struct {
	Location maps.LatLng `json:"location,omitempty"`
	Name     string      `json:"name,omitempty"`
	Vicinity string      `json:"vicinity,omitempty"`
}

var results []datum

var apiKey = os.Getenv("GOOGLE_API_KEY")
var mapboxAccessToken = os.Getenv("MAPBOX_ACCESS_TOKEN")

var (
	optKeyword  = flag.String("keyword", "donut", "Keyword to search for")
	optType     = flag.String("type", "restaurant|bakery", "Types to search for (delimited|by|pipe")
	optLocation = flag.String("location", "Portland, OR", "Location")
)

var re = regexp.MustCompile("oodoo")

func loadAll(c *maps.Client, r maps.NearbySearchRequest, limit int) error {
	resp, err := c.NearbySearch(context.Background(), &r)
	if err != nil {
		return err
	}

	for _, r := range resp.Results {
		if re.FindString(r.Name) == "" {
			results = append(results, datum{
				Location: r.Geometry.Location,
				Name:     r.Name,
				Vicinity: r.Vicinity,
			})
		}
	}

	resultCount := len(resp.Results)
	if resultCount >= limit || resp.NextPageToken == "" {
		return nil
	}

	// Take a deep, rate-limited breath before carrying on
	time.Sleep(5 * time.Second)

	nextR := maps.NearbySearchRequest{PageToken: resp.NextPageToken}
	nextLimit := limit - resultCount
	return loadAll(c, nextR, nextLimit)
}

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

	ctx := context.Background()
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed creating client: %s", err)
	}

	// Look up location
	log.Printf("Finding '%s'...\n", *optLocation)
	loc, err := c.Geocode(ctx, &maps.GeocodingRequest{
		Address: *optLocation,
	})

	if err != nil {
		log.Fatalf("failed geocoding: %s", err)
	} else if len(loc) < 1 {
		log.Fatalf("no geocoding results for '%s'", *optLocation)
	} else if len(loc) < 1 {
		log.Fatalf("more than one geocoding result for '%s'. Narrow it down!", *optLocation)
	}

	latLng := loc[0].Geometry.Location
	log.Printf("Found '%s' at %3.4f, %3.4f\n", *optLocation, latLng.Lat, latLng.Lng)

	initialRequest := maps.NearbySearchRequest{
		Type:     maps.PlaceType(*optType),
		Radius:   10000,
		Keyword:  *optKeyword,
		Location: &latLng,
	}
	maxResults := 100

	log.Println("Searching places...")
	if err := loadAll(c, initialRequest, maxResults); err != nil {
		log.Fatalf("failed searching: %s", err)
	}

	log.Println("Serializing results...")
	serializedResults, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("failed serializing results: %s", err)
	}

	log.Println("Printing up your new map...")
	compactResults := new(bytes.Buffer)
	if err := json.Compact(compactResults, serializedResults); err != nil {
		log.Fatalf("failed compacting JSON results: %s", err)
	}

	log.Println("We're done! Find the goods in ./dist...")
	templatize(dir, latLng, compactResults.Bytes())
}
