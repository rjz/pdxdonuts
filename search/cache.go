package search

import (
	"encoding/json"
	"fmt"
	"github.com/gitbookio/diskache"
	"time"
)

var CacheDir = ".cache"

type cachedSearch struct {
	RetrievedAt  time.Time `json:"retrievedAt"`
	Options      `json:"options"`
	SearchResult `json:"search"`
}

func cacheKey(o *Options) string {
	return fmt.Sprintf("%s!%s!%s!%s!%s", o.Address, o.Type, o.Keyword, o.Radius, o.Limit)
}

type Cache struct {
	*diskache.Diskache
}

func (sc *Cache) Get(o *Options) (*SearchResult, error) {
	if data, isCached := sc.Diskache.Get(cacheKey(o)); isCached {
		var sc cachedSearch
		if err := json.Unmarshal(data, &sc); err != nil {
			return nil, err
		}
		return &sc.SearchResult, nil
	}
	return nil, nil
}

func (sc *Cache) Set(o *Options, s *SearchResult) error {
	key := cacheKey(o)
	data, err := json.Marshal(cachedSearch{
		RetrievedAt:  time.Now(),
		Options:      *o,
		SearchResult: *s,
	})
	if err != nil {
		return err
	}

	return sc.Diskache.Set(key, data)
}

func NewCache() *Cache {
	dk, err := diskache.New(&diskache.Opts{
		Directory: CacheDir,
	})
	if err != nil {
		return nil
	}
	return &Cache{dk}
}
