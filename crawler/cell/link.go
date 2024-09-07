package cell

import (
	"free_proxy_pool/util/cas"
	"sync"
)

// LinkMap 记录所有爬取代理的链接，大于3则不再爬取
type LinkMap struct {
	body map[string]int

	lock sync.Locker
}

func NewLinkMap() *LinkMap {
	return &LinkMap{
		body: map[string]int{},
		lock: cas.NewSpinLock(),
	}
}

func (l *LinkMap) Add(link string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.body[link]; !ok {
		l.body[link] = 0
	}
	l.body[link]++
}

func (l *LinkMap) Check(link string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.body[link] <= 3
}
