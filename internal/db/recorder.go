package db

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	defaultCacheExpiration = 6 * time.Hour
	defaultCacheCleanup    = 10 * time.Minute
	defaultCacheSize       = 50
)

type Recorder struct {
	pool *cache.Cache
}

func NewRecorder() *Recorder {
	c := cache.New(defaultCacheExpiration, defaultCacheCleanup)
	c.OnEvicted(func(key string, v any) {
		logrus.Debugf("[db][recorder] key:%s deleted", key)
		l, ok := v.(*lru.Cache[int, any])
		if ok {
			putlru(l)
		}
	})

	return &Recorder{
		pool: c,
	}
}

func (r *Recorder) Record(userid string, v any) error {
	lru_, ok := r.pool.Get(userid)
	if !ok {
		l := getlru()
		l.Add(0, v)
		r.pool.Set(userid, l, cache.DefaultExpiration)
		return nil
	} else {
		l := lru_.(*lru.Cache[int, any])
		key := l.Len()
		if key >= defaultCacheSize {
			key = l.Keys()[l.Len()-1]
			key++
		}
		l.Add(key, v)
		r.pool.Set(userid, l, cache.DefaultExpiration)
	}

	return nil
}

func (r *Recorder) Get(userid string) ([]any, error) {
	lru_, ok := r.pool.Get(userid)
	if !ok {
		return nil, nil
	}
	l := lru_.(*lru.Cache[int, any])
	if l.Len() == 0 {
		return nil, nil
	}

	res := make([]any, 0, l.Len())
	for _, key := range l.Keys() {
		v, ok := l.Get(key)
		if ok {
			res = append(res, v)
		}
	}
	return res, nil
}

var lruPool = sync.Pool{New: func() any {
	l, _ := lru.New[int, any](defaultCacheSize)
	return l
}}

func getlru() *lru.Cache[int, any] {
	return lruPool.Get().(*lru.Cache[int, any])
}

func putlru(l *lru.Cache[int, any]) {
	if l != nil {
		l.Purge()
		lruPool.Put(l)
	}
}
