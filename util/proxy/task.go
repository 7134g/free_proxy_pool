package proxy

import "sync"

var (
	lock       sync.Locker
	errorCount int
)

func taskIncError() {
	lock.Lock()
	defer lock.Unlock()

	errorCount++
}

func TaskCheckError(max int) bool {
	lock.Lock()
	defer lock.Unlock()

	flag := errorCount >= max
	if flag {
		errorCount = 0
	}
	return flag
}
