package btils

import "sync/atomic"

type ThreaderManager[T any] struct {
	channel chan T

	workers  int
	callback func(in T)

	counter int64
}

func NewThreadManager[T any](workers int, callback func(in T)) *ThreaderManager[T] {
	tm := &ThreaderManager[T]{
		channel: make(chan T, workers),

		workers:  workers,
		callback: callback,
	}

	return tm
}

func (tm *ThreaderManager[T]) Start() {
	for i := 0; i < tm.workers; i++ {
		go func() {
			for in := range tm.channel {
				tm.callback(in)
				atomic.AddInt64(&tm.counter, -1)
			}
		}()
	}
}

func (tm *ThreaderManager[T]) Feed(in T) {
	atomic.AddInt64(&tm.counter, 1)
	tm.channel <- in
}

func (tm *ThreaderManager[T]) IsDone() bool {
	return atomic.LoadInt64(&tm.counter) == 0
}

func (tm *ThreaderManager[T]) Stop() {
	close(tm.channel)
}
