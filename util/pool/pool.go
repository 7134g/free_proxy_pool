package pool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultCleanIntervalTime = 2 * time.Second

	Waiting int = 0
	Working int = 1
	Idling  int = 2
)

type Task struct {
	TaskFunc func([]interface{})
	Param    []interface{}
}

type worker struct {
	task          chan *Task
	startWaitTime time.Time
	status        int
	mtx           sync.RWMutex
}

func (w *worker) changeStatus(status int) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.status = status
}

func (w *worker) currentStatus() int {
	w.mtx.RLock()
	defer w.mtx.RUnlock()
	return w.status
}

type Pool struct {
	capacity       int32
	workerCount    int32
	pruneDuration  time.Duration
	expiryDuration time.Duration
	workers        []*worker
	release        chan struct{}
	lock           sync.Locker
	once           sync.Once
	closed         bool
	wg             sync.WaitGroup
	Ctx            context.Context
}

func (p *Pool) Running() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *Pool) Cap() int32 {
	return atomic.LoadInt32(&p.capacity)
}

func NewPool(size int32, prune bool, expiry_time time.Duration) (*Pool, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		capacity:       size,
		release:        make(chan struct{}, 1),
		pruneDuration:  2 * time.Millisecond,
		expiryDuration: expiry_time,
		lock:           newSpinLock(),
		Ctx:            ctx,
	}
	if prune && size > 200 {
		go p.periodicallyPurge()
	}
	return p, cancel
}

func (p *Pool) periodicallyPurge() {
	heartbeat := time.NewTicker(p.pruneDuration)
	defer heartbeat.Stop()
	var (
		currentTime time.Time
		unInitTime  time.Time
	)
	for !p.closed {
		select {
		case <-heartbeat.C:
			if p.Running() <= 1 {
				continue
			}
			var pruned = 0
			for {
				pruned = 0
				currentTime = time.Now()
				p.lock.Lock()
				for i, w := range p.workers {
					if w.currentStatus() == Waiting {
						if w.startWaitTime == unInitTime {
							continue
						} else if currentTime.Sub(w.startWaitTime) >= p.expiryDuration {
							close(w.task)
							w.changeStatus(Idling)
							p.workers = append(p.workers[:i], p.workers[i+1:]...)
							atomic.AddInt32(&p.workerCount, -1)
							pruned++
							break
						}
					}
				}
				p.lock.Unlock()
				if p.Running() <= 1 || pruned == 0 {
					break
				}
			}
		case <-p.release:
			return
		}
	}
}

func (p *Pool) Submit(task *Task) error {
	if p == nil {
		return errors.New("this pool is nil")
	}
	if p.closed {
		return errors.New("this pool has been closed")
	}
	w := p.getWorker()
	w.changeStatus(Working)
	p.wg.Add(1)
	w.task <- task
	return nil
}

func (p *Pool) getWorker() *worker {
	var wk *worker

	// 如果超过上限，先清理空闲 worker
	if p.Running() > p.Cap() {
		for !p.closed && p.Running() > p.Cap() {
			trimmed := false
			p.lock.Lock()
			for i, w := range p.workers {
				if w.currentStatus() == Waiting {
					close(w.task)
					p.workers = append(p.workers[:i], p.workers[i+1:]...)
					atomic.AddInt32(&p.workerCount, -1)
					trimmed = true
					break
				}
			}
			p.lock.Unlock()
			if !trimmed {
				break
			}
		}
	}

	// 找空闲 worker 或创建新的
	for wk == nil && !p.closed {
		if p.Running() == p.Cap() {
			p.lock.Lock()
			for _, w := range p.workers {
				if w.currentStatus() == Waiting {
					wk = w
					break
				}
			}
			p.lock.Unlock()
			if wk == nil {
				time.Sleep(time.Duration(100) * time.Microsecond)
			}
		} else {
			p.lock.Lock()
			if p.Running() < p.Cap() {
				wk = &worker{
					task:   make(chan *Task, 1),
					status: Waiting,
					mtx:    sync.RWMutex{},
				}
				p.workers = append(p.workers, wk)
				atomic.AddInt32(&p.workerCount, 1)
				p.startWork(wk)
			}
			p.lock.Unlock()
		}
	}

	return wk
}

func (p *Pool) startWork(worker *worker) {
	go func() {
		for {
			select {
			case f, ok := <-worker.task:
				if !ok {
					return
				}
				f.TaskFunc(f.Param)
				p.wg.Done()
				worker.changeStatus(Waiting)
				worker.startWaitTime = time.Now()
			case <-p.Ctx.Done():
				worker.changeStatus(Waiting)
			}
		}
	}()
}

func (p *Pool) ReSize(size int32) {
	if size == p.Cap() {
		return
	}
	if size <= 0 {
		size = 1
	}
	atomic.StoreInt32(&p.capacity, size)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	p.once.Do(func() {
		close(p.release)
		p.lock.Lock()
		for _, worker := range p.workers {
			close(worker.task)
		}
		p.lock.Unlock()
		p.workers = p.workers[:0]
		atomic.StoreInt32(&p.workerCount, 0)
		p.closed = true
	})
}
