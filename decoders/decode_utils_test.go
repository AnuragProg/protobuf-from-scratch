package decoders

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLittleEndianVarint(t *testing.T) {
	testData := []struct {
		expected           uint64
		littleEndianStream io.Reader
	}{
		{ // 1 byte long barint
			expected: 0b0000_0100,
			littleEndianStream: bytes.NewBuffer([]byte{
				0b0000_0100,
			}),
		},
		{ // 2 byte long barint
			expected: 0b0000_0010_0001_0001,
			littleEndianStream: bytes.NewBuffer([]byte{
				0b1001_0001,
				0b0000_0100,
			}),
		},
		{ // 3 byte long barint
			expected: 0b0000_0100_0000_0100_0101_0000,
			littleEndianStream: bytes.NewReader([]byte{
				0b1101_0000, // least significant byte
				0b1000_1000,
				0b0001_0000, // most significant byte
			}),
		},
	}
	for _, data := range testData {
		expected := data.expected
		littleEndianStream := data.littleEndianStream
		actual, err := parseLittleEndianVarint(littleEndianStream)
		if err != nil {
			t.Fatal(err.Error())
		}

		assert.Equal(t, expected, actual)
	}
}

func TestParseString(t *testing.T) {
	testData := []struct {
		expected string
		stream   *bytes.Reader // little endian length followed by string
	}{

		{
			expected: "helloworld",
			stream: bytes.NewReader(
				append(
					[]byte{0b1010}, // 10
					[]byte("helloworld")...,
				),
			),
		},
		{
			expected: "aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oO4pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP1qQ2rR3sS4tT5uU6vV7wW8xX9yY0zZ1aA2bB3cC4dD5eE6fF7gG8hH9iI0jJ1kK2lL3mM4nN5oO6pP7qQ8rR9sS0tT1uU2vV3wW4xX5yY6zZ7aA8bB9cC0dD1eE2fF3gG4hH5iI6jJ7kK8lL9mM0nN1oO2pP3qQ4rR5sS6tT7uU8vV9wW0xX1yY2zZ3",
			stream: bytes.NewReader(
				append(
					[]byte{0b1011_1000, 0b10}, // 312
					[]byte("aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oO4pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP1qQ2rR3sS4tT5uU6vV7wW8xX9yY0zZ1aA2bB3cC4dD5eE6fF7gG8hH9iI0jJ1kK2lL3mM4nN5oO6pP7qQ8rR9sS0tT1uU2vV3wW4xX5yY6zZ7aA8bB9cC0dD1eE2fF3gG4hH5iI6jJ7kK8lL9mM0nN1oO2pP3qQ4rR5sS6tT7uU8vV9wW0xX1yY2zZ3")...,
				),
			),
		},
	}

	for _, data := range testData {
		expected := data.expected
		stream := data.stream
		actual, err := parseString(stream)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, expected, actual)
	}

}
