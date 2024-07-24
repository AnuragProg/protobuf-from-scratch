package encoders

import (
	"fmt"
	"protobuf-from-scratch/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeProjectType(t *testing.T) {

	testData := types.ProjectType{
		Name:        "John Doe",
		Description: "lorum ipsum",
		Timestamp:   302939020390,
		Tags:        []string{"tag1", "tag2"},
	}

	nameTag := byte((types.PROJECT_TYPE_NAME_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	nameLenBin := byte(len(testData.Name)) // 8
	nameBin := []byte(testData.Name)

	descTag := byte((types.PROJECT_TYPE_DESCRIPTION_FIELD_NO << 3) | types.WIRE_TYPE_STRING)
	descLenBin := byte(len(testData.Description)) // 11
	descBin := []byte(testData.Description)

	timestampTag := byte((types.PROJECT_TYPE_TIMESTAMP_FIELD_NO << 3) | types.WIRE_TYPE_VARINT)
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

	expected := []byte{}

	// format for string tag(fieldno+wiretype)-len-value
	expected = append(expected, nameTag)
	expected = append(expected, nameLenBin)
	expected = append(expected, nameBin...)

	expected = append(expected, descTag)
	expected = append(expected, descLenBin)
	expected = append(expected, descBin...)

	// format for varint tag(fieldno+wiretype)-value (len is not given and value is calculated using continuation bit)
	expected = append(expected, timestampTag)
	expected = append(expected, timestampBin...)
	fmt.Printf("timestamp = %v\t", testData.Timestamp)
	fmt.Printf("timestamp tag = %+v\t", timestampTag)
	fmt.Printf("timestamp bin = %+v\t\n", timestampBin)

	// format for list of string, array of (tag(fieldno+wiretype)-len-value
	for _, projectTag := range testData.Tags {
		projectTagLenBin := byte(len(projectTag)) // assuming length is not going to exceed 8 bits | 1 byte
		projectTagBin := []byte(projectTag)

		expected = append(expected, tagsTag)
		expected = append(expected, projectTagLenBin)
		expected = append(expected, projectTagBin...)

		fmt.Println(projectTag)
		fmt.Printf("tagsTag = %+v\n", tagsTag)
		fmt.Printf("projectTagLenBin = %+v\n", projectTagLenBin)
		fmt.Printf("projectTagBin = %+v\n", projectTagBin)
	}

	actual := EncodeProjectType(testData)
	fmt.Printf("expected: %+v\n", expected)
	fmt.Printf("actual: %+v\n", actual)

	assert.Equal(t, expected, actual)
}
