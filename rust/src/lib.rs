#![feature(test)]

extern crate test;

use std::sync::{Arc, Mutex, Barrier};
use std::sync::atomic::{AtomicU64, Ordering};
use std::thread;



fn add_using_mutex(sum_result: Arc<Mutex<f64>>) {
    let partial_sum = 42.0;
    let mut sum_result = sum_result.lock().unwrap();
    *sum_result += partial_sum;
}

fn add_using_barrier_and_mutex(sum_result: Arc<Mutex<f64>>, barrier: Arc<Barrier>) {
    let partial_sum = 42.0;
    barrier.wait();
    let mut sum_result = sum_result.lock().unwrap();
    *sum_result += partial_sum;
}

fn add_using_cas(sum: Arc<AtomicU64>) {
    let partial_sum = 42.0;
    loop {
        // Rust doesn't have atomic f64 because hardware doesn't support it,
        // so we use a u64 as container for the compare-and-swap operation.
        let sum_old = sum.load(Ordering::SeqCst);
        let sum_new = f64::from_bits(sum_old) + partial_sum;
        if sum.compare_exchange(sum_old, sum_new.to_bits(), Ordering::SeqCst, Ordering::SeqCst).is_ok() {
            break;
        }
    }
}


#[cfg(test)]
mod tests {
    use super::*;

    const NUMBER_OF_THREADS: usize = 256;

    #[bench]
    fn bench_mutex(b: &mut test::Bencher) {        
        b.iter(|| {
            let sum = Arc::new(Mutex::new(0.0));
            thread::scope(|thread_scope| {
                for _ in 0..NUMBER_OF_THREADS {
                    let sum = Arc::clone(&sum);
                    thread_scope.spawn(|| {
                        add_using_mutex(sum);
                    });
                }
            })
        });
    }

    #[bench]
    fn bench_barrier_and_mutex(b: &mut test::Bencher) {       
        b.iter(|| {
            let sum = Arc::new(Mutex::new(0.0));
            let barrier = Arc::new(Barrier::new(NUMBER_OF_THREADS));
            thread::scope(|thread_scope| {
                for _ in 0..NUMBER_OF_THREADS {
                    let sum = Arc::clone(&sum);
                    let barrier = Arc::clone(&barrier);
                    thread_scope.spawn(|| {
                        add_using_barrier_and_mutex(sum, barrier);
                    });
                }
            })
        });
    }

    #[bench]
    fn bench_cas(b: &mut test::Bencher) {        
        b.iter(|| {
            // Rust doesn't have atomic f64 because hardware doesn't support it,
            // so we use a u64 as a container for the f64 bits instead.
            let sum = Arc::new(AtomicU64::new(0.0_f64.to_bits()));
            thread::scope(|thread_scope| {
                for _ in 0..NUMBER_OF_THREADS {
                    let sum = Arc::clone(&sum);
                    thread_scope.spawn(|| {
                        add_using_cas(sum);
                    });
                }
            })
        });
    }
}