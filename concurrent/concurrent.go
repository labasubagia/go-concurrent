package concurrent

import (
	"context"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func NewConcurrent(data [][]time.Duration, worker int, isShowLog bool) *Concurrent {
	return &Concurrent{
		data:      data,
		isShowLog: isShowLog,
		worker:    worker,
	}
}

type Concurrent struct {
	data      [][]time.Duration
	isShowLog bool
	worker    int
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

func (c Concurrent) UseChannel() {
	channels := make(chan any, c.worker)
	for i, inner := range c.data {
		for j, duration := range inner {
			channels <- 0
			go func(i, j int, d time.Duration) {
				defer func() {
					<-channels
				}()
				time.Sleep(d)
				if c.isShowLog {
					log.Println("process", i, "item", j)
				}
			}(i, j, duration)
		}
	}
	close(channels)
}

func (c Concurrent) UseWaiGroupAndChannelNested() {
	var wg sync.WaitGroup
	chInner := make(chan any, c.worker)
	chOuter := make(chan any, c.worker)

	for i, inner := range c.data {
		chOuter <- 0
		wg.Add(1)
		go func(i int, inner []time.Duration) {
			defer func() {
				defer wg.Done()
				<-chOuter
			}()

			for j, duration := range inner {
				chInner <- 0
				wg.Add(1)
				go func(i, j int, d time.Duration) {
					defer func() {
						defer wg.Done()
						<-chInner
					}()
					time.Sleep(d)
					if c.isShowLog {
						log.Println("process", i, "item", j)
					}
				}(i, j, duration)
			}
		}(i, inner)
	}

	wg.Wait()
	close(chInner)
	close(chOuter)
}

func (c Concurrent) UseSemaphoreFoOuter() {
	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(c.worker))
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
	if err := sem.Acquire(ctx, int64(c.worker)); err != nil {
		log.Fatal("failed acquire")
	}
}

func (c Concurrent) UseSemaphoreForInner() {
	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(c.worker))
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
	if err := sem.Acquire(ctx, int64(c.worker)); err != nil {
		log.Fatal("failed acquire")
	}
}

// UseSemaphoreNested1
// this total semaphore is innerSize * outerSize in nested array
func (c Concurrent) UseSemaphoreNested1() {
	ctx := context.Background()

	semOuter := semaphore.NewWeighted(int64(c.worker))
	for i, inner := range c.data {

		if err := semOuter.Acquire(ctx, 1); err != nil {
			log.Fatal("failed acquire for semaphore outer")
		}
		go func(i int, inner []time.Duration) {
			defer semOuter.Release(1)

			semInner := semaphore.NewWeighted(int64(c.worker))
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
			if err := semInner.Acquire(ctx, int64(c.worker)); err != nil {
				log.Fatal("failed wait semaphore inner")
			}

		}(i, inner)
	}
	if err := semOuter.Acquire(ctx, int64(c.worker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}

// UseSemaphoreNested2
// this total semaphore is innerSize + outerSize in nested array
func (c Concurrent) UseSemaphoreNested2() {
	ctx := context.Background()

	semOuter := semaphore.NewWeighted(int64(c.worker))
	semInner := semaphore.NewWeighted(int64(c.worker))

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
	if err := semOuter.Acquire(ctx, int64(c.worker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}

// UseSemaphoreNested3
// ! DANGER: Need to fix, do not use
func (c Concurrent) UseSemaphoreNested3() {
	ctx := context.Background()

	sem := semaphore.NewWeighted(int64(c.worker))
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

	if err := sem.Acquire(ctx, int64(c.worker)); err != nil {
		log.Fatal("failed acquire semaphore outer")
	}
}
