package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kr/pretty"
	"github.com/rjz/pdxdonuts/generate"
	"github.com/rjz/pdxdonuts/search"
	"googlemaps.github.io/maps"
	"log"
	"strings"
)

type MapOpts struct {
	GoogleApiKey      string   `json:"googleApiKey" validate:"required"`
	Keyword           string   `json:"keyword" validate:"required"`
	Types             []string `json:"types" validate:"required"`
	Location          string   `json:"location" validate:"required"`
	generate.PageData `json:"pageData" validate:"required"`
}

func RenderMap(ctx context.Context, opts MapOpts, dest string) error {
	c, err := maps.NewClient(maps.WithAPIKey(opts.GoogleApiKey))
	if err != nil {
		log.Printf("failed creating client: %s", err)
		return err
	}

	log.Println("Finding the results...")
	searchOpts := search.Options{
		Address: opts.Location,
		Type:    strings.Join(opts.Types, "|"),
		Keyword: opts.Keyword,
		Limit:   100,
		Radius:  10000, // m
	}

	s, err := search.Do(&searchOpts, c)
	if err != nil {
		log.Printf("failed searching: %s", err)
		return err
	}

	log.Println("Serializing results...")
	serializedResults, err := json.Marshal(s.Places)
	if err != nil {
		log.Printf("failed serializing results: %s", err)
		return err
	}

	log.Println("Printing up your new map...")
	compactResults := new(bytes.Buffer)
	if err := json.Compact(compactResults, serializedResults); err != nil {
		log.Printf("failed compacting JSON results: %s", err)
		return err
	}

	pageData := opts.PageData
	pageData.Lat = s.LatLng.Lat
	pageData.Lng = s.LatLng.Lng

	pretty.Log(opts.PageData)

	pageData.Data = compactResults.String()

	if err := generate.Do(dest, &pageData); err != nil {
		log.Printf("failed building template: %s", err)
		return err
	}

	log.Printf("We're done! Find the goods in %s...\n", dest)
	return nil
}
