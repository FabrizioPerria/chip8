build:
	go build .

run:
	go run .

runprof:
	go run . -cpuprofile=cpu.prof -memprofile=mem.prof

showcpuprof:
	go tool pprof -http=:8080 cpu.prof

showmemprof:
	go tool pprof -http=:8080 mem.prof

test:
	go test -v .


.PHONY: build run runprof test showcpuprof showmemprof
