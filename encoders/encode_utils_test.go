package encoders

import (
	"bytes"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeLittleEndianVarint(t *testing.T) {
	testData := []struct {
		varint   uint64
		expected []byte
	}{
		{ // 1 byte long varint
			varint: 0b0010_1000,
			expected: []byte{
				0b0_010_1000, // {continuation_bit}_{3bit}_{4bit}
				// -------------------actual data--
			},
		},
		{ // 2 byte long varint
			varint: 0b1111_1010_1011_1100,
			expected: []byte{
				0b1_011_1100, // {continuation_bit}_{3bit}_{4bit}
				0b1_111_0101, // -------------------actual data--
				0b0_000_0011,
			},
		},
		{ // 3 byte long varint
			varint: 0b0111_1010_1111_1010_1011_1100,
			expected: []byte{
				0b1_011_1100, // {continuation_bit}_{3bit}_{4bit}
				0b1_111_0101, // -------------------actual data--
				0b1_110_1011,
				0b0_000_0011,
			},
		},
	}

	for _, data := range testData {
		actual := bytes.NewBuffer([]byte{})
		serializeLittleEndianVarint(actual, data.varint)
		assert.ElementsMatch(t, data.expected, actual.Bytes())
	}
}

func TestSerializeString(t *testing.T) {
	testData := []struct {
		stringData string
		expected   []byte // little endian length followed by string
	}{
		{
			stringData: "helloworld",
			expected: slices.Concat(
				[]byte{0b0_000_1010}, // 10 
				[]byte("helloworld"),
			),
		},
		{
			stringData: "aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oO4pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP1qQ2rR3sS4tT5uU6vV7wW8xX9yY0zZ1aA2bB3cC4dD5eE6fF7gG8hH9iI0jJ1kK2lL3mM4nN5oO6pP7qQ8rR9sS0tT1uU2vV3wW4xX5yY6zZ7aA8bB9cC0dD1eE2fF3gG4hH5iI6jJ7kK8lL9mM0nN1oO2pP3qQ4rR5sS6tT7uU8vV9wW0xX1yY2zZ3",
			expected: slices.Concat(
				[]byte{0b1_011_1000, 0b0_000_0010}, // 312 in little endian
				[]byte("aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oO4pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP1qQ2rR3sS4tT5uU6vV7wW8xX9yY0zZ1aA2bB3cC4dD5eE6fF7gG8hH9iI0jJ1kK2lL3mM4nN5oO6pP7qQ8rR9sS0tT1uU2vV3wW4xX5yY6zZ7aA8bB9cC0dD1eE2fF3gG4hH5iI6jJ7kK8lL9mM0nN1oO2pP3qQ4rR5sS6tT7uU8vV9wW0xX1yY2zZ3"),
			),
		},
	}

	for _, data := range testData{
		actual := bytes.NewBuffer([]byte{})
		serializeString(actual, data.stringData)
		assert.ElementsMatch(t, data.expected, actual.Bytes())
	}
}
