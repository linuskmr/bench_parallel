# Benchmarking mutex, barrier and atomic compare-and-swap in Rust


## Execution

```bash
$ cargo +nightly bench --jobs 1
```

Explanation:

- `+nightly` selects the nighty branch, because benchmarking without external crates [is currently unstable](https://doc.rust-lang.org/cargo/commands/cargo-bench.html)
- `--jobs 1` runs only one benchmark at a time. Running benchmarks in parallel makes no sense here, because the threads itself already use multiple threads.