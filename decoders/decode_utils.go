package decoders

import (
	"errors"
	"fmt"
	"io"
	
	"protobuf-from-scratch/utils"
)


/*
Deserializes little endian varint
*/
func deserializeLittleEndianVarint(stream io.Reader) (uint64, error){
	var varint uint64
	temp := make([]byte, 1)
	position := 0

	for {
		_, err := stream.Read(temp)
		if err == io.EOF {
			return varint, io.EOF
		}
		if err != nil {
			return 0, err
		}

		// store the current byte
		currentByte := temp[0]

		// determine whether to do the next loop
		shallContinue := (currentByte & utils.CONTINUATION_BIT) > 0

		// remove continuation bit
		currentByte &= ^utils.CONTINUATION_BIT

		// add to varint
		varint |= uint64(currentByte) << position

		position += 7

		if !shallContinue {
			break
		}
	}

	return varint, nil
}


/*
Deserializes binary string data

Expected binary format : <little-endian-string-len><actual-string>
*/
func deserializeString(stream io.Reader) (string, error) {
	length, err := deserializeLittleEndianVarint(stream)
	if err != nil {
		return "", err
	}

	str := make([]byte, length)
	parsedLen, err := stream.Read(str)
	if err != nil {
		return "", err
	}
	if uint64(parsedLen) != length {
		return "", errors.New(fmt.Sprintf("invalid length encountered, expected: %v, got: %v", length, parsedLen))
	}

	return string(str), nil
}
