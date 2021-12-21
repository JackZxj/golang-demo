package base

import (
	"log"
	"sync"
	"time"
)

// memoryCache uses to cache users' name and id info
var memoryCache sync.Map

// TODO: 可配置超时时间

// cacheDuration means cache expires time
var cacheDuration = time.Second * 10

func init() {
	// TODO: 可配置什么时候做清理

	// 每天4点清除一下缓存的用户信息
	go func() {
		now := time.Now()
		next := now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), 4, 0, 0, 0, next.Location())
		clearUserCache(next.Sub(now), 24*time.Hour)
	}()
}

type CachedUserInfo struct {
	Name      string
	ID        int64
	CreatTime time.Time
}

type Base struct {
	Name string

	CreatorInfo CachedUserInfo
}

func (c *Base) BasePrepare() {
	var err error
	if c.Name != "" {
		c.CreatorInfo, err = loadUser(c.Name)
		if err != nil {
			log.Fatalf("can not get creator ID: %v", err)
		}
	}
}

func loadUser(name string) (u CachedUserInfo, err error) {
	v, ok := memoryCache.Load(name)
	u.Name = name
	if ok {
		if uu, okok := v.(CachedUserInfo); okok {
			if uu.CreatTime.Add(cacheDuration).Before(time.Now()) {
				log.Println("got the user but expired")
				goto GET_USER
			} else {
				return uu, nil
			}
		}
		log.Println("unknown values, get a new one")
	}

GET_USER:
	u.ID = time.Now().Unix()
	u.CreatTime = time.Now()
	memoryCache.Store(u.Name, u)
	return
}

// clearUserCache clears the user's info cache every other period of time.
// The period is at least one hour. You can use "delay" to delay the
// first cleanup.
func clearUserCache(delay, period time.Duration) {
	if period.Hours() < time.Hour.Hours() {
		period = time.Hour
	}
	if delay.Seconds() < 0 {
		delay = 0
	}
	for {
		now := time.Now()
		next := now.Add(period).Add(delay)
		delay = 0
		t := time.NewTimer(next.Sub(now))
		<-t.C
		memoryCache.Range(func(key, value interface{}) bool {
			if uu, okok := value.(CachedUserInfo); okok {
				if uu.CreatTime.Add(cacheDuration).Before(time.Now()) {
					memoryCache.Delete(key)
				}
			} else {
				memoryCache.Delete(key)
			}
			return true
		})
	}
}

func ReadCache() {
	log.Println("#### read cache:")
	memoryCache.Range(func(key, value interface{}) bool {
		if uu, okok := value.(CachedUserInfo); okok {
			log.Println(uu)
		} else {
			memoryCache.Delete(key)
		}
		return true
	})
	log.Println("#### read done")
}
