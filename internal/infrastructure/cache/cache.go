package cache

import (
	"context"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var c = gocache.New(10*time.Minute, 15*time.Minute)

func Set(_ context.Context, k string, v any, d time.Duration) { c.Set(k, v, d) }
func Get(_ context.Context, k string) (any, bool)             { return c.Get(k) }

// ───── health helpers ─────
type Stats struct{ Items int }

func StatsCurrent() Stats { return Stats{Items: c.ItemCount()} }
