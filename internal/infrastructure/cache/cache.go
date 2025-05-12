package cache

import (
	"context"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var c = gocache.New(10*time.Minute, 15*time.Minute)

func Set(_ context.Context, k string, v any, d time.Duration) { c.Set(k, v, d) }
func Get(_ context.Context, k string) (any, bool)             { return c.Get(k) }

func ItemCount() int { return c.ItemCount() }

// маленькая структурка для /health
type Stats struct{ Items int }

func StatHealth() Stats { return Stats{Items: ItemCount()} }
func Stat() Stats {
	return Stats{Items: len(c.Items())}
}
