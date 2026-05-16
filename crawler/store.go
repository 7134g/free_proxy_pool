package crawler

import (
	"free_proxy_pool/config"
	"math/rand"
	"sort"
	"sync"
)

type Store struct {
	lock sync.RWMutex
	body map[string]*proxy

	slice []*proxy
}

func (s *Store) add(u string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.body[u] = newProxy(u)
	redisSyncScore(u, s.body[u].Score)
}

func (s *Store) get(u string) (*proxy, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	data, ok := s.body[u]
	return data, ok
}

func (s *Store) inc(u string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.body[u] == nil {
		s.body[u] = newProxy(u)
	}
	s.body[u].Score++
	s.body[u].sucCount++
	redisSyncScore(u, s.body[u].Score)

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

	if s.body[u].Score <= 0 {
		delete(s.body, u)
		redisZRem(u)
	} else {
		redisSyncScore(u, s.body[u].Score)
	}

	return true
}

func (s *Store) sort() {
	s.lock.Lock()
	list := make([]*proxy, 0, len(s.body))
	for _, d := range s.body {
		list = append(list, d)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})

	if len(list) > config.Cfg.PoolCap {
		s.slice = list[:config.Cfg.PoolCap]
	} else {
		s.slice = list
	}
	s.lock.Unlock()
}

func (s *Store) GetMaxList() []*proxy {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if len(s.slice) == 0 {
		return nil
	}

	if len(s.slice) > 10 {
		return s.slice[:10]
	} else {
		return s.slice
	}
}

func (s *Store) GetOnce(index int) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.slice == nil || len(s.slice) == 0 {
		return ""
	}

	if index < 0 {
		index = 0
	}

	if index >= len(s.slice) {
		index = len(s.slice) - 1
	}

	if index == 0 {
		// max
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

func (s *Store) Random() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if len(s.slice) == 0 {
		return ""
	}

	index := rand.Intn(len(s.slice))
	return s.slice[index].Link
}

func (s *Store) GetCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.slice)
}

// Slice 返回当前排序后切片的一份拷贝，外部调用方无需加锁
func (s *Store) Slice() []*proxy {
	s.lock.RLock()
	defer s.lock.RUnlock()
	result := make([]*proxy, len(s.slice))
	copy(result, s.slice)
	return result
}

func (s *Store) Del(link string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.body, link)
	redisZRem(link)
}
