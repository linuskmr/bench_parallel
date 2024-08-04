package main

import (
	"sync"
	"testing"
	"sync/atomic"
)

const N_THREADS = 256

func BenchmarkMutex(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		sum := 0.0

		sumMutex := sync.Mutex{}

		finish := sync.WaitGroup{}
		finish.Add(N_THREADS)

		for n := 0; n < N_THREADS; n++ {
			go AddUsingMutex(&sum, &sumMutex, &finish)
		}

		finish.Wait()
	}
}

func BenchmarkBarrierAndMutex(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		sum := 0.0

		sumMutex := sync.Mutex{}

		sumBarrier := sync.WaitGroup{}
		sumBarrier.Add(N_THREADS)

		finish := sync.WaitGroup{}
		finish.Add(N_THREADS)

		for n := 0; n < N_THREADS; n++ {
			go AddUsingBarrierAndMutex(&sum, &sumBarrier, &sumMutex, &finish)
		}

		finish.Wait()
	}
}

func BenchmarkCaS(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		sum := atomic.Value{}
		sum.Store(0.0)

		finish := sync.WaitGroup{}
		finish.Add(N_THREADS)

		for n := 0; n < N_THREADS; n++ {
			go AddUsingCaS(&sum, &finish)
		}

		finish.Wait()
	}
}


/// AddUsingMutex locks `sumMutex` before adding to `sum`.
func AddUsingMutex(sum *float64, sumMutex *sync.Mutex, finish *sync.WaitGroup) {
	partialSum := 42.0
	
	sumMutex.Lock()
	defer sumMutex.Unlock()
	*sum += partialSum

	finish.Done()
}

/// AddUsingBarrierAndMutex waits for all threads to reach `sumBarrier` before locking `sumMutex` and adding to `sum`.
func AddUsingBarrierAndMutex(sum *float64, sumBarrier *sync.WaitGroup, sumMutex *sync.Mutex, finish *sync.WaitGroup) {
	partialSum := 42.0

	// Signal that *this* goroutine is done
	sumBarrier.Done()

	// Wait for all other goroutines to finish
	sumBarrier.Wait()
	
	{
		sumMutex.Lock()
		defer sumMutex.Unlock()
		*sum += partialSum
	}

	finish.Done()
}

/// AddUsingCaS uses a compare-and-swap loop to add to `sum`.
func AddUsingCaS(sum *atomic.Value, finish *sync.WaitGroup) {
	partialSum := 42.0
	
	for {
		sumOld := sum.Load().(float64)
		sumNew := sumOld + partialSum
		if sum.CompareAndSwap(sumOld, sumNew) {
			break
		}
	}

	finish.Done()
}