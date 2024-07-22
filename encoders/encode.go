package encoders

import (
	"bytes"
	"errors"
	"io"
	"protobuf-from-scratch/types"
)

func Encode(data interface{}) (io.Reader, error) {
	switch data := data.(type) {
	case types.ProjectType:
		return EncodeProjectType(data), nil
	default:
		return nil, errors.New("unknown data type")
	}
}

func EncodeProjectType(data types.ProjectType) io.Reader {

	serializedDataStreams := []io.Reader{}

	nameTagStream := serializeLittleEndianVarint((types.PROJECT_TYPE_NAME_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	nameLenAndValueStream := serializeString(data.Name)
	serializedDataStreams = append(serializedDataStreams, nameTagStream, nameLenAndValueStream)

	descTagStream := serializeLittleEndianVarint((types.PROJECT_TYPE_DESCRIPTION_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	descLenAndValueStream := serializeString(data.Description)
	serializedDataStreams = append(serializedDataStreams, descTagStream, descLenAndValueStream)

	timestampTagStream := serializeLittleEndianVarint((types.PROJECT_TYPE_TIMESTAMP_FIELD_NO << 3) | types.WIRE_TYPE_VARINT)
	timestampValueStream := serializeLittleEndianVarint(data.Timestamp)
	serializedDataStreams = append(serializedDataStreams, timestampTagStream, timestampValueStream)

	tagsTagBytes, _ := io.ReadAll(serializeLittleEndianVarint((types.PROJECT_TYPE_TAGS_FIELD_NO << 3) | types.WIRE_TYPE_STRING))
	for _, projectTag := range data.Tags {
		projectTagLenAndValueStream := serializeString(projectTag)
		serializedDataStreams = append(serializedDataStreams, bytes.NewReader(tagsTagBytes), projectTagLenAndValueStream)
	}

	return io.MultiReader(serializedDataStreams...)
}
