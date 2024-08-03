package main

import (
	"math/rand"
	"sync"
	"testing"
	"sync/atomic"
)

const N_THREADS = 4

func BenchmarkMutex(bench *testing.B) {
	arr := randomFloatArray(-9999, 9999, bench.N)
	sumResult := 0.0
	sumResultMutex := sync.Mutex{}
	for n := 0; n < N_THREADS; n++ {
		go AddUsingMutex(arr, &sumResult, &sumResultMutex)
	}
}

func BenchmarkBarrierMutex(bench *testing.B) {
	arr := randomFloatArray(-9999, 9999, bench.N)
	sumResult := 0.0
	sumResultMutex := sync.Mutex{}
	sumResultBarrier := sync.WaitGroup{}
	sumResultBarrier.Add(N_THREADS)
	for n := 0; n < N_THREADS; n++ {
		go AddUsingBarrierAndMutex(arr, &sumResult, &sumResultBarrier, &sumResultMutex)
	}
}

func BenchmarkCaS(bench *testing.B) {
	arr := randomFloatArray(-9999, 9999, bench.N)
	sumResult := atomic.Value{}
	sumResult.Store(0.0)
	for n := 0; n < N_THREADS; n++ {
		go AddUsingCaS(arr, &sumResult)
	}
}

func AddUsingMutex(arr []float64, sumResult *float64, sumResultMutex *sync.Mutex) {
	sum := 0.0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}
	
	sumResultMutex.Lock()
	defer sumResultMutex.Unlock()
	*sumResult += sum
}

func AddUsingBarrierAndMutex(arr []float64, sumResult *float64, sumResultBarrier *sync.WaitGroup, sumResultMutex *sync.Mutex) {
	sum := 0.0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}
	// Signal that *this* goroutine is done
	sumResultBarrier.Done()

	// Wait for all other goroutines to finish
	sumResultBarrier.Wait()
	
	sumResultMutex.Lock()
	defer sumResultMutex.Unlock()
	*sumResult += sum
}

func AddUsingCaS(arr []float64, sumResult *atomic.Value) {
	sum := 0.0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}
	
	for {
		sumResultOld := sumResult.Load().(float64)
		sumResultNew := sumResultOld + sum
		if sumResult.CompareAndSwap(sumResultOld, sumResultNew) {
			break
		}
	}
}


func randomFloatArray(min, max float64, n int) []float64 {
    res := make([]float64, n)
    for i := range res {
        res[i] = min + rand.Float64() * (max - min)
    }
    return res
}