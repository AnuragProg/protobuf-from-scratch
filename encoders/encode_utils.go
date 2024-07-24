package encoders

import (
	"bytes"
	"protobuf-from-scratch/utils"
)

func serializeLittleEndianVarint(buffer *bytes.Buffer, varint uint64) {
	// estimatedSize := 9 /* 8 bytes + 1*8 continuation bits */
	// varintBytes := make([]byte, 0, estimatedSize)

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

		// varintBytes = append(varintBytes, currentByte)
		buffer.WriteByte(currentByte)

		// break if next part doesn't exist
		if varint == 0 {
			break
		}
	}

}


/*
Serializes binary string data

Output binary format : <little-endian-string-len><actual-string>
*/
func serializeString(buffer *bytes.Buffer, data string) {
	serializeLittleEndianVarint(buffer, uint64(len(data)))
	buffer.WriteString(data)
}
