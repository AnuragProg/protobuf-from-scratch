# Description

Protobuf from scratch: Developed with aim to implement protocol buffers from scratch and to make it's encoding/decoding faster than that of json's with reduced payload size through the use of binary

## Features

- Protobuf serialization and deserialization
- Performance benchmarking
- Optimization to reduce latency spikes
- Profiling with pprof
- Visual performance comparisons

## Installation

Clone the repository and install dependencies:

```sh
git clone https://github.com/AnuragProg/protobuf-from-scratch
cd protobuf-from-scratch
go mod download
```

## Usage

Run the client, which in turn uses builtin ```encoding/json``` go library to compare it with custom implemented encoders and decoders of protobuf

```sh
make client
```

Result of which is ```cpu.prof``` file which can be viewed through ```go tool pprof -http=:3000 cpu.prof```

## Optimizations (Mistakes that led to latency spikes in encoding/decoding and their fixes)

Context: We have a type name `ProjectType` with schema
```go
type ProjectType struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Timestamp   uint64   `json:"timestamp"`
	Tags        []string `json:"tags"`
}
```
that we are serializing/deserializing using custom encoders/decoders (```encoders.*```) and ```encoding/json``` encoders/decoders

*Note: methods of type ```encoders.*``` are to be focused on as they are the ones we are trying to optimize, e.g. encoders.EncodeProjectType, encoders.serializeString, etc*

### Optimization 1

Cause: Was returning ```io.Reader``` as result of serialization of fields, resulting in frequent reads from reader anytime we need to process the underlying bytes, resulting in large allocations for new slices

Fix: Returned ```[]byte``` directly, preventing overhead of frequent reads from the reader to access it

#### Before Optimization:

![Before Optimization](./assets/before-optimization-continuous-reading-writing-to-ioreader.png)

#### After Optimization:

![After Optimization](./assets/after-optimization-continuous-reading-writing-to-ioreader.png)


### Optimization 2

Cause: Multiple huge allocations causing latency spikes due to mallocs, slice resizing & copying of underlying data in newly resized slice

Fix: Use of single buffer, and enforcing all the subsequent serializers to use it directly to push the results to

#### Before Optimization (comparison with ```encoding/json```'s Marshal):

![Before Optimization](./assets/before-optimization-large-proto-encoding-latency-due-to-multiple-allocations.png)

#### After Optimization:

![After Optimization](./assets/after-optimization-large-proto-encoding-latency-due-to-multiple-allocations.png)


