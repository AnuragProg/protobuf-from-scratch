package decoders

import (
	"errors"
	"fmt"
	"io"
	"protobuf-from-scratch/types"
	"protobuf-from-scratch/utils"
)

func Decode(stream io.ReadCloser, _type types.Type) (interface{}, error) {
	defer stream.Close()

	switch _type {
	case types.PROJECT_TYPE:
		return DecodeProjectType(stream)
	default:
		return nil, errors.New("invalid _type")
	}
}

func DecodeProjectType(stream io.Reader) (types.ProjectType, error) {
	projectType := types.ProjectType{
		Name: "",
		Description: "",
		Timestamp: 0,
		Tags: []string{},
	}

	for {
		// parse tag
		tag, err := parseLittleEndianVarint(stream)
		if err == io.EOF {
			break
		}
		if err != nil {
			return projectType, err
		}

		// separate field no and wire type
		wireType := tag & uint64(utils.WIRE_TYPE_BITS)
		fieldNo := tag >> 3

		// parse respective field
		switch wireType {
		case types.WIRE_TYPE_STRING:
			value, err := parseString(stream)
			if err != nil {
				return projectType, err
			}
			switch fieldNo {
			case types.PROJECT_TYPE_NAME_FIELD_NO:
				projectType.Name = value
			case types.PROJECT_TYPE_DESCRIPTION_FIELD_NO:
				projectType.Description = value
			case types.PROJECT_TYPE_TAGS_FIELD_NO:
				projectType.Tags = append(projectType.Tags, value)
			default:
				return projectType, errors.New(fmt.Sprintf("invalid field no: %v", fieldNo))
			}
		case types.WIRE_TYPE_VARINT:
			value, err := parseLittleEndianVarint(stream)
			if err != nil {
				return projectType, err
			}
			switch fieldNo{
			case types.PROJECT_TYPE_TIMESTAMP_FIELD_NO:
				projectType.Timestamp = value
			default:
				return projectType, errors.New(fmt.Sprintf("invalid field no: %v", fieldNo))
			}
		default:
			return projectType, errors.New(fmt.Sprintf("invalid wire type: %v", wireType))
		}
	}
	return projectType, nil
}
