package namesys

import (
	"time"

	path "gx/ipfs/QmNgXoHgXU1HzNb2HEZmRww9fDKE9NfDsvQwWLHiKHpvKM/go-path"
)

func (ns *mpns) cacheGet(name string) (path.Path, bool) {
	if ns.cache == nil {
		return "", false
	}

	ientry, ok := ns.cache.Get(name)
	if !ok {
		return "", false
	}

	entry, ok := ientry.(cacheEntry)
	if !ok {
		// should never happen, purely for sanity
		log.Panicf("unexpected type %T in cache for %q.", ientry, name)
	}

	if time.Now().Before(entry.eol) {
		return entry.val, true
	}

	ns.cache.Remove(name)

	return "", false
}

func (ns *mpns) cacheSet(name string, val path.Path, ttl time.Duration) {
	if ns.cache == nil || ttl <= 0 {
		return
	}
	ns.cache.Add(name, cacheEntry{
		val: val,
		eol: time.Now().Add(ttl),
	})
}

type cacheEntry struct {
	val path.Path
	eol time.Time
}
