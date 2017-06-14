package search

import (
	"fmt"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
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
	Type    string
	Keyword string
	Radius  uint
	Limit   int
}

// Search wraps search results
type Search struct {
	Places []*Place
	LatLng maps.LatLng
	client *maps.Client
	ctx    context.Context
}

func (s *Search) loadAll(r maps.NearbySearchRequest, limit int) error {
	resp, err := s.client.NearbySearch(s.ctx, &r)
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
	return s.loadAll(nextR, limit-resultCount)
}

// Do finds all results near an address
func (s *Search) Do(address string, opts *Options) error {
	// Look up location
	loc, err := s.client.Geocode(s.ctx, &maps.GeocodingRequest{
		Address: address,
	})

	if err != nil {
		return fmt.Errorf("failed geocoding: %s", err)
	} else if len(loc) < 1 {
		return fmt.Errorf("no geocoding results for '%s'", address)
	} else if len(loc) < 1 {
		return fmt.Errorf("more than one geocoding result for '%s'. Narrow it down!", address)
	}

	s.LatLng = loc[0].Geometry.Location

	initialRequest := maps.NearbySearchRequest{
		Type:     maps.PlaceType(opts.Type),
		Radius:   opts.Radius,
		Keyword:  opts.Keyword,
		Location: &s.LatLng,
	}
	if err := s.loadAll(initialRequest, opts.Limit); err != nil {
		return fmt.Errorf("failed searching: %s", err)
	}

	return nil
}

func NewSearch(client *maps.Client) *Search {
	return &Search{
		client: client,
		ctx:    context.Background(),
	}
}
