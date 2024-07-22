package decoders

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"protobuf-from-scratch/types"
)

func TestDecodeProjectType(t *testing.T) {
	expected := types.ProjectType{
		Name: "John Doe",
		Description: "lorum ipsum",
		Timestamp: 302939020390,
		Tags: []string{"tag1", "tag2"},
	}

	// TODO: add varint conversion to the project type field nos
	// since project type field no.s do not exceed 5bits hence not doing varint conversions
	nameTag := byte((types.PROJECT_TYPE_NAME_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	nameLenBin := byte(len(expected.Name)) // 8 
	nameBin := []byte(expected.Name)

	descTag := byte((types.PROJECT_TYPE_DESCRIPTION_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	descLenBin := byte(len(expected.Description)) // 11 
	descBin := []byte(expected.Description)

	timestampTag :=	byte((types.PROJECT_TYPE_TIMESTAMP_FIELD_NO << 3) | types.WIRE_TYPE_VARINT)
	timestampBin := []byte{ 
		// little endian varint for 302939020390
		// raw binary 0b_1000_110_1000_100_0100_100_1010_011_1000_110_0110
		0b1_110_0110, // {continuation_bit}_{3bit}_{4bit}
		0b1_011_1000, // -------------------actual data 
		0b1_100_1010,
		0b1_100_0100,
		0b1_110_1000,
		0b0_000_1000,
	}

	tagsTag := byte((types.PROJECT_TYPE_TAGS_FIELD_NO << 3) | types.WIRE_TYPE_STRING)

	dataBin := []byte{}

	// format for string tag(fieldno+wiretype)-len-value
	dataBin = append(dataBin, nameTag)
	dataBin = append(dataBin, nameLenBin)
	dataBin = append(dataBin, nameBin...)

	dataBin = append(dataBin, descTag)
	dataBin = append(dataBin, descLenBin)
	dataBin = append(dataBin, descBin...)


	// format for varint tag(fieldno+wiretype)-value (len is not given and value is calculated using continuation bit)
	dataBin = append(dataBin, timestampTag)
	dataBin = append(dataBin, timestampBin...)


	// format for list of string, array of (tag(fieldno+wiretype)-len-value
	for _, projectTag := range expected.Tags {
		projectTagLenBin := byte(len(projectTag)) // assuming length is not going to exceed 8 bits | 1 byte
		projectTagBin := []byte(projectTag)


		dataBin = append(dataBin, tagsTag)
		dataBin = append(dataBin, projectTagLenBin)
		dataBin = append(dataBin, projectTagBin...)
	}

	dataStream := bytes.NewReader(dataBin)
	actual, err := DecodeProjectType(dataStream)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("expected: %+v\n", expected)
	fmt.Printf("actual: %+v\n", actual)
	assert.Equal(t, expected, actual)
}
