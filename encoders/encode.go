package encoders

import (
	"bytes"
	"protobuf-from-scratch/types"
)

// func Encode(data interface{}) (io.Reader, error) {
// 	switch data := data.(type) {
// 	case types.ProjectType:
// 		return EncodeProjectType(data), nil
// 	default:
// 		return nil, errors.New("unknown data type")
// 	}
// }

// NOTE: Optimization 1: Returning bytes array directly instead of wrapping around in io reader, avoids unnecessary reads and writes to access the bytes
// NOTE: Optimization 2: Graphs showed large allocations(mallocs & growing of slices) & showed latency spikes when pre allocating, hence used single buffer to allow utils to directly push the result directly to buffer, removing multiple allocations in place of one
func EncodeProjectType(data types.ProjectType) []byte {
	estimatedSize := len(data.Name) + 
					 len(data.Description) + 
					 8 /* for timestamp 64 bits */ + 
					 10*len(data.Tags)

	serializedBuffer := bytes.NewBuffer(make([]byte, 0, estimatedSize))
	// serializedBuffer := []byte{}

	// nameTagBytes := serializeLittleEndianVarint((types.PROJECT_TYPE_NAME_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	// nameLenAndValueBytes := serializeString(data.Name)
	serializeLittleEndianVarint(serializedBuffer, (types.PROJECT_TYPE_NAME_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	serializeString(serializedBuffer, data.Name)
	// serializedBuffer.Write(nameTagBytes)
	// serializedBuffer.Write(nameLenAndValueBytes)
	// serializedBuffer = append(serializedBuffer, nameTagBytes...)
	// serializedBuffer = append(serializedBuffer, nameLenAndValueBytes...)

	// descTagBytes := serializeLittleEndianVarint((types.PROJECT_TYPE_DESCRIPTION_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	// descLenAndValueBytes := serializeString(data.Description)
	serializeLittleEndianVarint(serializedBuffer, (types.PROJECT_TYPE_DESCRIPTION_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	serializeString(serializedBuffer, data.Description)
	// serializedBuffer.Write(descTagBytes)
	// serializedBuffer.Write(descLenAndValueBytes)
	// serializedBuffer = append(serializedBuffer, descTagBytes...)
	// serializedBuffer = append(serializedBuffer, descLenAndValueBytes...)

	// timestampTagBytes := serializeLittleEndianVarint((types.PROJECT_TYPE_TIMESTAMP_FIELD_NO << 3) | types.WIRE_TYPE_VARINT)
	// timestampValueBytes := serializeLittleEndianVarint(data.Timestamp)
	serializeLittleEndianVarint(serializedBuffer, (types.PROJECT_TYPE_TIMESTAMP_FIELD_NO << 3) | types.WIRE_TYPE_VARINT)
	serializeLittleEndianVarint(serializedBuffer, data.Timestamp)
	// serializedBuffer.Write(timestampTagBytes)
	// serializedBuffer.Write(timestampValueBytes)
	// serializedBuffer = append(serializedBuffer, timestampTagBytes...)
	// serializedBuffer = append(serializedBuffer, timestampValueBytes...)

	tempTagsTagBuffer := bytes.NewBuffer([]byte{})
	// tagsTagBytes := serializeLittleEndianVarint((types.PROJECT_TYPE_TAGS_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	serializeLittleEndianVarint(tempTagsTagBuffer, (types.PROJECT_TYPE_TAGS_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	for _, projectTag := range data.Tags {
		serializedBuffer.Write(tempTagsTagBuffer.Bytes())
		serializeString(serializedBuffer, projectTag)
		// projectTagLenAndValueStream := serializeString(projectTag)
		// serializedBuffer.Write(tagsTagBytes)
		// serializedBuffer.Write(projectTagLenAndValueStream)
		// serializedBuffer = append(serializedBuffer, tagsTagBytes...)
		// serializedBuffer = append(serializedBuffer, projectTagLenAndValueStream...)
	}
	return serializedBuffer.Bytes()
}
