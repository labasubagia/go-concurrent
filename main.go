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
	cn.UseChannel()
	// end here

	elapsed := time.Since(start)
	log.Println("took ", elapsed)
}
