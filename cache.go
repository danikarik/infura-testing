package main

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

const (
	// DefaultTTL for cache items.
	DefaultTTL = 2 * time.Minute
	// ETHPedingSetTTL shows how long to cache pending transaction sets.
	ETHPedingSetTTL = 5 * time.Minute
	// GasTTL show how long to cache gas responses.
	GasTTL = 30 * time.Second
	// GasKeyTemplate is cache key for gas price.
	GasKeyTemplate = "gas/price/network/%d"
	// ETHPendingSetTemplate is a cache key for pending transaction unqie sets.
	ETHPendingSetTemplate = "eth/pending/set/%s"
)

var (
	// The global cache instance
	httpCache = cache.New(DefaultTTL, DefaultTTL*2)
)
