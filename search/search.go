package search

import (
	"fmt"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"log"
	"time"
)

// Place represents a single entry for the map
type Place struct {
	Id       string      `json:"id,omitempty"`
	Location maps.LatLng `json:"location,omitempty"`
	Name     string      `json:"name,omitempty"`
	Vicinity string      `json:"vicinity,omitempty"`
}

// Options expose search options
type Options struct {
	Address string
	Type    string
	Keyword string
	Radius  uint
	Limit   int
}

// SearchResult wraps search results
type SearchResult struct {
	Places []*Place
	LatLng maps.LatLng
}

func (s *SearchResult) loadAll(ctx context.Context, client *maps.Client, r maps.NearbySearchRequest, limit int) error {
	resp, err := client.NearbySearch(ctx, &r)
	if err != nil {
		return err
	}

	for _, r := range resp.Results {
		s.Places = append(s.Places, &Place{
			Id:       fmt.Sprintf("g!%s", r.PlaceID),
			Location: r.Geometry.Location,
			Name:     r.Name,
			Vicinity: r.Vicinity,
		})
	}

	resultCount := len(resp.Results)
	if resultCount >= limit || resp.NextPageToken == "" {
		return nil
	}

	// Take a deep, rate-limited breath before carrying on
	time.Sleep(5 * time.Second)

	nextR := maps.NearbySearchRequest{PageToken: resp.NextPageToken}
	return s.loadAll(ctx, client, nextR, limit-resultCount)
}

// do finds all results near an address
func (s *SearchResult) do(opts *Options, client *maps.Client) error {

	ctx := context.Background()

	// Look up location
	loc, err := client.Geocode(ctx, &maps.GeocodingRequest{
		Address: opts.Address,
	})

	if err != nil {
		return fmt.Errorf("failed geocoding: %s", err)
	} else if len(loc) < 1 {
		return fmt.Errorf("no geocoding results for '%s'", opts.Address)
	} else if len(loc) < 1 {
		return fmt.Errorf("more than one geocoding result for '%s'. Narrow it down!", opts.Address)
	}

	s.LatLng = loc[0].Geometry.Location

	initialRequest := maps.NearbySearchRequest{
		Type:     maps.PlaceType(opts.Type),
		Radius:   opts.Radius,
		Keyword:  opts.Keyword,
		Location: &s.LatLng,
	}
	if err := s.loadAll(ctx, client, initialRequest, opts.Limit); err != nil {
		return fmt.Errorf("failed searching: %s", err)
	}

	return nil
}

// Do performs a cacheable search
func Do(o *Options, c *maps.Client) (s *SearchResult, err error) {
	cache := NewCache()
	if cache != nil {
		if s, err = cache.Get(o); err != nil {
			log.Printf("cache error: %s\n", err)
		} else if s != nil {
			log.Println("cache hit, using cached results")
			return
		}
	}

	log.Println("cache miss, asking google")
	s = new(SearchResult)
	if err = s.do(o, c); err != nil {
		return
	}

	if cache != nil {
		if err := cache.Set(o, s); err != nil {
			log.Printf("failed caching results: %s\n", err)
		}
	}

	return
}
