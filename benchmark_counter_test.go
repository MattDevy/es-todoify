package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkMutexCounter benchmarks counter increment using sync.Mutex
func BenchmarkMutexCounter(b *testing.B) {
	var mu sync.Mutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// BenchmarkAtomicCounter benchmarks counter increment using atomic operations
func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}

// BenchmarkMutexCounterNoContention benchmarks mutex with no parallel contention
func BenchmarkMutexCounterNoContention(b *testing.B) {
	var mu sync.Mutex
	var counter int64

	for i := 0; i < b.N; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

// BenchmarkAtomicCounterNoContention benchmarks atomic with no parallel contention
func BenchmarkAtomicCounterNoContention(b *testing.B) {
	var counter atomic.Int64

	for i := 0; i < b.N; i++ {
		counter.Add(1)
	}
}

// BenchmarkRWMutexRead benchmarks reading counter with sync.RWMutex
func BenchmarkRWMutexRead(b *testing.B) {
	var mu sync.RWMutex
	var counter int64 = 42

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = counter
			mu.RUnlock()
		}
	})
}

// BenchmarkAtomicRead benchmarks reading counter with atomic operations
func BenchmarkAtomicRead(b *testing.B) {
	var counter atomic.Int64
	counter.Store(42)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = counter.Load()
		}
	})
}

// BenchmarkRWMutexReadNoContention benchmarks RWMutex read with no parallel contention
func BenchmarkRWMutexReadNoContention(b *testing.B) {
	var mu sync.RWMutex
	var counter int64 = 42

	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = counter
		mu.RUnlock()
	}
}

// BenchmarkAtomicReadNoContention benchmarks atomic read with no parallel contention
func BenchmarkAtomicReadNoContention(b *testing.B) {
	var counter atomic.Int64
	counter.Store(42)

	for i := 0; i < b.N; i++ {
		_ = counter.Load()
	}
}

// BenchmarkMixedRWMutex benchmarks mixed read/write with RWMutex (90% reads, 10% writes)
func BenchmarkMixedRWMutex(b *testing.B) {
	var mu sync.RWMutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if i%10 == 0 {
				// Write operation (10%)
				mu.Lock()
				counter++
				mu.Unlock()
			} else {
				// Read operation (90%)
				mu.RLock()
				_ = counter
				mu.RUnlock()
			}
		}
	})
}

// BenchmarkMixedAtomic benchmarks mixed read/write with atomics (90% reads, 10% writes)
func BenchmarkMixedAtomic(b *testing.B) {
	var counter atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if i%10 == 0 {
				// Write operation (10%)
				counter.Add(1)
			} else {
				// Read operation (90%)
				_ = counter.Load()
			}
		}
	})
}

// BenchmarkChannelWrite benchmarks writing to a buffered channel
func BenchmarkChannelWrite(b *testing.B) {
	ch := make(chan int64, 1000)

	// Goroutine to drain the channel
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ch <- 1
		}
	})
	b.StopTimer()

	close(ch)
	<-done
}

// BenchmarkChannelRead benchmarks reading from a buffered channel
func BenchmarkChannelRead(b *testing.B) {
	ch := make(chan int64, 1000)

	// Goroutine to fill the channel
	done := make(chan struct{})
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- 42
		}
		close(done)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			<-ch
		}
	})
	b.StopTimer()

	<-done
}

// BenchmarkChannelWriteNoContention benchmarks channel write with no parallel contention
func BenchmarkChannelWriteNoContention(b *testing.B) {
	ch := make(chan int64, 1000)

	// Goroutine to drain the channel
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- 1
	}
	b.StopTimer()

	close(ch)
	<-done
}

// BenchmarkChannelReadNoContention benchmarks channel read with no parallel contention
func BenchmarkChannelReadNoContention(b *testing.B) {
	ch := make(chan int64, 1000)

	// Goroutine to fill the channel
	done := make(chan struct{})
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- 42
		}
		close(done)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ch
	}
	b.StopTimer()

	<-done
}

// BenchmarkMixedChannel benchmarks mixed read/write with buffered channel (90% reads, 10% writes)
func BenchmarkMixedChannel(b *testing.B) {
	readCh := make(chan int64, 1000)
	writeCh := make(chan int64, 1000)

	// Goroutine to handle writes
	done := make(chan struct{})
	go func() {
		counter := int64(0)
		for range writeCh {
			counter++
			// Put value into read channel
			select {
			case readCh <- counter:
			default:
			}
		}
		close(done)
	}()

	// Pre-fill read channel
	for i := 0; i < 100; i++ {
		readCh <- 42
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if i%10 == 0 {
				// Write operation (10%)
				select {
				case writeCh <- 1:
				default:
				}
			} else {
				// Read operation (90%)
				select {
				case <-readCh:
				default:
				}
			}
		}
	})
	b.StopTimer()

	close(writeCh)
	<-done
}
