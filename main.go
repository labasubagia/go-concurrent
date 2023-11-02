package main

import (
	"log"
	"time"

	"github.com/labasubagia/go-concurrent/concurrent"
	"github.com/labasubagia/go-concurrent/util"
)

func main() {
	cn := concurrent.NewConcurrent(
		util.GenNestedDuration(10, 10, time.Second),
		10,
		true,
	)

	start := time.Now()

	// use here
	// cn.UseWaitGroup()
	cn.UseSemaphoreNested2()

	// end here

	elapsed := time.Since(start)
	log.Println("took ", elapsed)
}
