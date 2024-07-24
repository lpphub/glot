package helper

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var LocalCache *cache.Cache

func initCache() {
	LocalCache = cache.New(5*time.Minute, 10*time.Minute)
}
