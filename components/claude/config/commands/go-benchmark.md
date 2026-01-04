---
description: Run and analyze Go benchmarks
argument-hint: [pattern]
allowed-tools: Read, Bash
---

## Task

Help the user run and analyze Go benchmarks for performance testing.

## Quick Benchmark

Using dotfiles function:
```bash
gobench
# Runs all benchmarks

gobench BenchmarkAdd
# Runs benchmarks matching pattern
```

## Manual Benchmarking

### Run Benchmarks
```bash
# All benchmarks
go test -bench=. ./...

# Specific benchmark
go test -bench=BenchmarkAdd ./...

# With memory allocation stats
go test -bench=. -benchmem ./...
```

### Benchmark Output
```
BenchmarkAdd-8     1000000000    0.5 ns/op    0 B/op    0 allocs/op
```
- `BenchmarkAdd-8`: Name and GOMAXPROCS
- `1000000000`: Iterations run
- `0.5 ns/op`: Nanoseconds per operation
- `0 B/op`: Bytes allocated per op
- `0 allocs/op`: Allocations per op

## Writing Benchmarks

### Basic Benchmark
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```

### With Setup
```go
func BenchmarkProcess(b *testing.B) {
    // Setup (not timed)
    data := loadTestData()

    b.ResetTimer() // Start timing

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

### Sub-benchmarks
```go
func BenchmarkSort(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            data := generateData(size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                sort.Ints(data)
            }
        })
    }
}
```

### Parallel Benchmarks
```go
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            DoWork()
        }
    })
}
```

## Comparing Benchmarks

### Save Baseline
```bash
go test -bench=. -benchmem ./... > baseline.txt
```

### Compare with benchstat
```bash
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Run new benchmark
go test -bench=. -benchmem ./... > new.txt

# Compare
benchstat baseline.txt new.txt
```

Output:
```
name    old time/op    new time/op    delta
Add-8   0.50ns ± 1%    0.48ns ± 1%   -4.00%
```

## Profiling During Benchmarks

### CPU Profile
```bash
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Memory Profile
```bash
go test -bench=. -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Block Profile
```bash
go test -bench=. -blockprofile=block.prof ./...
```

## pprof Commands

```bash
# Interactive mode
go tool pprof cpu.prof

# Common commands:
# top        - Show top functions
# list func  - Show function source
# web        - Open in browser

# Direct web view
go tool pprof -http=:8080 cpu.prof
```

## Benchmark Best Practices

1. **Isolate what you're measuring**
2. **Use b.ResetTimer() after setup**
3. **Avoid compiler optimizations removing code**
4. **Run multiple times for consistency**
5. **Use benchstat for comparisons**

### Prevent Optimization
```go
var result int

func BenchmarkAdd(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = Add(1, 2)
    }
    result = r // Prevent dead code elimination
}
```

## Common Options

```bash
# Run for minimum time
go test -bench=. -benchtime=5s ./...

# Run specific count
go test -bench=. -count=10 ./...

# Disable tests, only benchmarks
go test -bench=. -run=^$ ./...
```

## Dotfiles Integration

- `gobench [pattern]` - Run benchmarks with optional filter
- `gotest [pattern]` - Run tests (alias for comparison)
