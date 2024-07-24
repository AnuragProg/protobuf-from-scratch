package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"protobuf-from-scratch/decoders"
	"protobuf-from-scratch/encoders"
	"protobuf-from-scratch/types"
	"runtime"
	"runtime/pprof"
	"slices"
	"strings"
	"sync"
	"time"
)

type Stats struct {
	Encodingp90, Encodingp95, Encodingp99 int64
	Decodingp90, Decodingp95, Decodingp99 int64
	Size                                  int
}

func Benchmark(fun func()) int64 {
	start := time.Now()
	fun()
	return time.Since(start).Nanoseconds()
}

func GetStats(data types.ProjectType, statsType string, concurrency uint64) (Stats, error) {

	var encodeFunc, decodeFunc func()
	var dataBytes []byte

	switch statsType {
	case "json":
		dataBytes, _ = json.Marshal(data)
		encodeFunc = func() {
			_, _ = json.Marshal(data)
		}
		decodeFunc = func() {
			var jsonDecodedData types.ProjectType
			json.Unmarshal(dataBytes, &jsonDecodedData)
		}
	case "proto":
		dataBytes = encoders.EncodeProjectType(data)
		encodeFunc = func() {
			_ = encoders.EncodeProjectType(data)
		}
		dataStream := bytes.NewReader(dataBytes)
		decodeFunc = func() {
			_, _ = decoders.DecodeProjectType(dataStream)
		}
	default:
		return Stats{}, errors.New("invalid stats type")
	}

	benchmark := struct {
		mut          sync.Mutex
		encLatencies []int64
		decLatencies []int64
	}{
		mut:          sync.Mutex{},
		encLatencies: make([]int64, 0, concurrency),
		decLatencies: make([]int64, 0, concurrency),
	}

	var wg sync.WaitGroup

	// benchmarking encoding
	wg.Add(int(concurrency))
	for range concurrency {
		go func() {
			defer wg.Done()
			defer benchmark.mut.Unlock()
			latency := Benchmark(encodeFunc)
			benchmark.mut.Lock()
			benchmark.encLatencies = append(benchmark.encLatencies, latency)
		}()
	}

	// benchmarking decoding
	wg.Add(int(concurrency))
	for range concurrency {
		go func() {
			defer wg.Done()
			defer benchmark.mut.Unlock()
			latency := Benchmark(decodeFunc)
			benchmark.mut.Lock()
			benchmark.decLatencies = append(benchmark.decLatencies, latency)
		}()
	}

	wg.Wait()

	// sorting the data
	slices.Sort(benchmark.encLatencies)
	slices.Sort(benchmark.decLatencies)

	if len(benchmark.encLatencies) != len(benchmark.decLatencies) {
		panic("latencies count not equal")
	}

	p90Idx := int((float32(concurrency-1) * 90) / 100)
	p95Idx := int((float32(concurrency-1) * 95) / 100)
	p99Idx := int((float32(concurrency-1) * 99) / 100)

	// fmt.Printf("Enc latencies: %+v\n", benchmark.encLatencies[p90Idx:])
	// fmt.Printf("Dec latencies: %+v\n", benchmark.decLatencies[p90Idx:])

	// fmt.Println("p90Idx = ", p90Idx)
	// fmt.Println("p95Idx = ", p95Idx)
	// fmt.Println("p99Idx = ", p99Idx)

	return Stats{
		Encodingp90: benchmark.encLatencies[p90Idx],
		Encodingp95: benchmark.encLatencies[p95Idx],
		Encodingp99: benchmark.encLatencies[p99Idx],

		Decodingp90: benchmark.decLatencies[p90Idx],
		Decodingp95: benchmark.decLatencies[p95Idx],
		Decodingp99: benchmark.decLatencies[p99Idx],
		Size:        len(dataBytes),
	}, nil

}

func main() {

	// for profiling
	profFile, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(profFile)
	defer pprof.StopCPUProfile()

	// small data
	data := types.ProjectType{
		Name:        "John Doe",
		Description: "Lorum ipsum",
		Timestamp:   uint64(time.Now().Unix()),
		Tags:        []string{"tag1", "tag2", "tag3"},
	}


	var conc uint64 = 10_000
	fmt.Println("Stats for small payload:")
	fmt.Println("==================================================")
	jsonStats, _ := GetStats(data, "json", conc)
	fmt.Printf("Json %v concurrency: %+v\n", conc, jsonStats)
	runtime.GC()
	protoStats, _ := GetStats(data, "proto", conc)
	fmt.Printf("Proto %v concurrency: %+v\n", conc, protoStats)
	runtime.GC()
	fmt.Println("==================================================")


	// large data
	singleTag := strings.Repeat("tag", 100)
	tagCount := 1000
	tags := make([]string, tagCount)
	for i := range tags {
		tags[i] = singleTag
	}
	data = types.ProjectType{
		Name:        strings.Repeat("John Doe", 100),        // Large name
		Description: strings.Repeat("Lorum ipsum", 100),    // Large description
		Timestamp:   uint64(time.Now().Unix()),                // Current timestamp
		Tags:        tags,            // Large tags
	}
	fmt.Println("Stats for large payload:")
	fmt.Println("==================================================")
	jsonStats, _ = GetStats(data, "json", conc)
	fmt.Printf("Json %v concurrency: %+v\n", conc, jsonStats)
	runtime.GC()
	protoStats, _ = GetStats(data, "proto", conc)
	fmt.Printf("Proto %v concurrency: %+v\n", conc, protoStats)
	runtime.GC()
	fmt.Println("==================================================")

}
