package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rjz/pdxdonuts/generate"
	"os"
	"strings"
)

var apiKey = os.Getenv("GOOGLE_API_KEY")
var mapboxAccessToken = os.Getenv("MAPBOX_ACCESS_TOKEN")

var (
	optDest     = flag.String("dest", "dist", "Output directory")
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

func cli() {
	if apiKey == "" {
		usageAndExit("Please specify GOOGLE_API_KEY")
	} else if mapboxAccessToken == "" {
		usageAndExit("Please specify MAPBOX_ACCESS_TOKEN")
	} else if optLocation == nil {
		usageAndExit("-location is required")
	}

	pageData, err := generate.LoadPageData("vars.json")
	if err != nil {
		panic("Failed reading vars.json")
	}
	pageData.MapboxAccessToken = mapboxAccessToken

	opts := MapOpts{
		GoogleApiKey: apiKey,
		Keyword:      *optKeyword,
		Location:     *optLocation,
		Types:        strings.Split(*optType, "|"),
		PageData:     *pageData,
	}

	ctx := context.Background()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := RenderMap(ctx, opts, dir); err != nil {
		panic(err)
	}
}
