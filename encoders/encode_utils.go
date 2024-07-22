package encoders

import (
	"bytes"
	"io"
	"protobuf-from-scratch/utils"
)

func serializeLittleEndianVarint(varint uint64) io.Reader {

	varintBytes := []byte{}

	for { // not going to put check here for varint != 0 as 0 could be the possible starting value for serializing

		// parse non continuation bits of current part of the varint
		currentByte := byte(varint & uint64(utils.NON_CONTINUATION_BITS))

		// 7 bits will be consumed and hence those needs to be removed
		varint >>= 7

		// check whether next part exists
		if varint != 0 { 

			// set continuation bit
			currentByte |= utils.CONTINUATION_BIT
		}

		varintBytes = append(varintBytes, currentByte)

		// break if next part doesn't exist
		if varint == 0 {
			break
		}
	}

	return bytes.NewReader(varintBytes)
}


/*
Serializes binary string data

Output binary format : <little-endian-string-len><actual-string>
*/
func serializeString(data string) io.Reader {
	length := len(data)
	lengthStream := serializeLittleEndianVarint(uint64(length))
	stringStream := bytes.NewReader([]byte(data))
	return io.MultiReader(
		lengthStream, 
		stringStream,
	)
}
