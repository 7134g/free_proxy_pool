package crawler

import (
	"free_proxy_pool/config"
	"math/rand"
	"sort"
	"sync"
)

type Store struct {
	lock sync.Mutex
	body map[string]*proxy

	slice []*proxy
}

func (s *Store) add(u string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.body[u] = &proxy{Link: u}
}

func (s *Store) get(u string) (*proxy, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.body[u]
	return data, ok
}

func (s *Store) inc(u string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.body[u] == nil {
		s.body[u] = &proxy{Link: u, Score: 10}
	}
	s.body[u].Score++
	s.body[u].sucCount++

	return true
}

func (s *Store) dnc(u string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.body[u] == nil {
		return false
	}
	s.body[u].Score--
	s.body[u].errCount++

	if s.body[u].Score < 0 {
		delete(s.body, u)
	}

	return true
}

func (s *Store) sort() {
	s.lock.Lock()
	list := make([]*proxy, 0)
	for _, d := range s.body {
		list = append(list, d)
	}
	s.lock.Unlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})

	if len(s.slice) > config.Cfg.PoolCap {
		s.slice = list[:config.Cfg.PoolCap]
	} else {
		s.slice = list
	}

}

func (s *Store) GetMaxList() []*proxy {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.slice) > 10 {
		return s.slice[:10]
	} else {
		return s.slice
	}
}

func (s *Store) GetOne(index int) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	if index < 0 {
		index = 0
	}

	if index >= len(s.slice) {
		index = len(s.slice) - 1
	}

	if index == 0 {
		ps := 0
		if len(s.slice) > 10 {
			ps = 10
		} else {
			ps = len(s.slice)
		}

		index = rand.Intn(ps)
	}

	return s.slice[index].Link
}

func (s *Store) GetCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.slice)
}

func (s *Store) Del(link string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.body, link)
}
