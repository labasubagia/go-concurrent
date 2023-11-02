package concurrent

import (
	"context"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func NewConcurrent(data [][]time.Duration, maxWorker int, isShowLog bool) *Concurrent {
	return &Concurrent{
		data:      data,
		isShowLog: isShowLog,
		maxWorker: maxWorker,
	}
}

type Concurrent struct {
	data      [][]time.Duration
	isShowLog bool
	maxWorker int
}

func (c Concurrent) UseWaitGroup() {
	var wg sync.WaitGroup
	for i, inner := range c.data {
		for j, duration := range inner {
			wg.Add(1)
			go func(i, j int, d time.Duration) {
				defer wg.Done()
				time.Sleep(d)
				if c.isShowLog {
					log.Println("process", i, "item", j)
				}
			}(i, j, duration)
		}
	}
	wg.Wait()
}

func (c Concurrent) UseSemaphoreFoOuter() {
	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(c.maxWorker))
	for i, inner := range c.data {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatal("failed acquire")
		}
		go func(i int, inner []time.Duration) {
			defer sem.Release(1)
			for j, duration := range inner {
				time.Sleep(duration)
				if c.isShowLog {
					log.Println("process", i, "item", j)
				}
			}
		}(i, inner)
	}
	if err := sem.Acquire(ctx, int64(c.maxWorker)); err != nil {
		log.Fatal("failed acquire")
	}
}

func (c Concurrent) UseSemaphoreForInner() {
	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(c.maxWorker))
	for i, inner := range c.data {
		for j, duration := range inner {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Fatal("failed acquire")
			}
			go func(i, j int, d time.Duration) {
				defer sem.Release(1)
				time.Sleep(d)
				if c.isShowLog {
					log.Println("process", i, "item", j)
				}
			}(i, j, duration)
		}
	}
	if err := sem.Acquire(ctx, int64(c.maxWorker)); err != nil {
		log.Fatal("failed acquire")
	}
}

// UseSemaphoreNested1
// this total semaphore is innerSize * outerSize in nested array
func (c Concurrent) UseSemaphoreNested1() {
	ctx := context.Background()

	semOuter := semaphore.NewWeighted(int64(c.maxWorker))
	for i, inner := range c.data {

		if err := semOuter.Acquire(ctx, 1); err != nil {
			log.Fatal("failed acquire for semaphore outer")
		}
		go func(i int, inner []time.Duration) {
			defer semOuter.Release(1)

			semInner := semaphore.NewWeighted(int64(c.maxWorker))
			for j, duration := range inner {
				if err := semInner.Acquire(ctx, 1); err != nil {
					log.Fatal("failed acquire for semaphore inner")
				}
				go func(i, j int, d time.Duration) {
					defer semInner.Release(1)
					time.Sleep(d)
					if c.isShowLog {
						log.Println("process", i, "item", j)
					}
				}(i, j, duration)
			}
			if err := semInner.Acquire(ctx, int64(c.maxWorker)); err != nil {
				log.Fatal("failed wait semaphore inner")
			}

		}(i, inner)
	}
	if err := semOuter.Acquire(ctx, int64(c.maxWorker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}

// UseSemaphoreNested2
// this total semaphore is innerSize + outerSize in nested array
func (c Concurrent) UseSemaphoreNested2() {
	ctx := context.Background()

	semOuter := semaphore.NewWeighted(int64(c.maxWorker))
	semInner := semaphore.NewWeighted(int64(c.maxWorker))

	for i, inner := range c.data {

		if err := semOuter.Acquire(ctx, 1); err != nil {
			log.Fatal("failed acquire for semaphore outer")
		}
		go func(i int, inner []time.Duration) {
			defer semOuter.Release(1)

			for j, duration := range inner {
				if err := semInner.Acquire(ctx, 1); err != nil {
					log.Fatal("failed acquire for semaphore inner")
				}
				go func(i, j int, d time.Duration) {
					defer semInner.Release(1)
					time.Sleep(d)
					if c.isShowLog {
						log.Println("process", i, "item", j)
					}
				}(i, j, duration)
			}

		}(i, inner)
	}

	// use the outermost semaphore
	if err := semOuter.Acquire(ctx, int64(c.maxWorker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}

// UseSemaphoreNested3
// ! DANGER: Need to fix, do not use
func (c Concurrent) UseSemaphoreNested3() {
	ctx := context.Background()

	sem := semaphore.NewWeighted(int64(c.maxWorker))
	for i, inner := range c.data {

		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatal("failed acquire for semaphore outer")
		}
		go func(i int, inner []time.Duration) {
			defer sem.Release(1)

			for j, duration := range inner {
				if err := sem.Acquire(ctx, 1); err != nil {
					log.Fatal("failed acquire for semaphore inner")
				}
				go func(i, j int, d time.Duration) {
					defer sem.Release(1)
					time.Sleep(d)
					if c.isShowLog {
						log.Println("process", i, "item", j)
					}
				}(i, j, duration)
			}

		}(i, inner)
	}

	if err := sem.Acquire(ctx, int64(c.maxWorker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}
