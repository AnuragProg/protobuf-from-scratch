# Project Title

Protobuf from scratch

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

### Optimization 1

#### Cause: Was returning io.Reader as result of serialization of fields, resulting in continuous read/write from reader while concatenating or processing underlying bytes

#### Fix: Returned []byte directly, preventing overhead of read/write to access it

#### Before Optimization:

![Before Optimization](./assets/before-optimization-continuous-reading-writing-to-ioreader.png)

#### After Optimization:

![After Optimization](./assets/after-optimization-continuous-reading-writing-to-ioreader.png)


### Optimization 2

#### Cause: Multiple huge allocations causing latency spikes due to mallocs, slice resizing & copying of underlying data in newly resized slice

#### Fix: Use of single buffer, and enforcing all the subsequent serializers to use it directly to push results to

#### Before Optimization:

![Before Optimization](./assets/before-optimization-large-proto-encoding-latency-due-to-multiple-allocations.png)

#### After Optimization:

![After Optimization](./assets/after-optimization-large-proto-encoding-latency-due-to-multiple-allocations.png)


