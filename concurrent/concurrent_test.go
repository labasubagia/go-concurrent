package concurrent_test

import (
	"testing"
	"time"

	"github.com/labasubagia/go-concurrent/concurrent"
	"github.com/labasubagia/go-concurrent/util"
)

var data = util.GenNestedDuration(10, 10, time.Millisecond*500)
var cn = concurrent.NewConcurrent(data, 10, false)

func BenchmarkWaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cn.UseWaitGroup()
	}
}

func BenchmarkSemaphoreOuter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cn.UseSemaphoreFoOuter()
	}
}

func BenchmarkSemaphoreInner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cn.UseSemaphoreForInner()
	}
}

func BenchmarkSemaphoreNested1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cn.UseSemaphoreNested1()
	}
}

func BenchmarkSemaphoreNested2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cn.UseSemaphoreNested2()
	}
}
