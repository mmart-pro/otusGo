#!/bin/bash
#go test -bench=. -count=7 -benchmem >bench_fast.txt
benchstat bench_slow.txt bench_fast.txt

#go test . -bench . -cpuprofile cpu.prof
#go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof .
#go tool pprof -http=":8090" cpu.prof
