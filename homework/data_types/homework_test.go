package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Integer interface {
	uint16 | uint32 | uint64
}

func ToLittleEndian[T Integer](number T) T {
	if number == 0 {
		return number
	}

	var result T

	size := int(unsafe.Sizeof(number))

	for i := 0; i < size; i++ {
		currentByte := number >> (8 * i) & 0xFF
		shift := 8 * (size - i - 1)
		result |= currentByte << shift
	}

	return result
}

func TestConversion(t *testing.T) {
	t.Run("uint16", func(t *testing.T) {
		tests := map[string]struct {
			number uint16
			result uint16
		}{
			"test case #1": {
				number: 0x0000,
				result: 0x0000,
			},
			"test case #2": {
				number: 0xFFFF,
				result: 0xFFFF,
			},
			"test case #3": {
				number: 0x00FF,
				result: 0xFF00,
			},
			"test case #4": {
				number: 0xFF00,
				result: 0x00FF,
			},
			"test case #5": {
				number: 0x0102,
				result: 0x0201,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := ToLittleEndian(test.number)

				assert.Equal(t, test.result, result)
			})
		}
	})

	t.Run("uint32", func(t *testing.T) {
		tests := map[string]struct {
			number uint32
			result uint32
		}{
			"test case #1": {
				number: 0x00000000,
				result: 0x00000000,
			},
			"test case #2": {
				number: 0xFFFFFFFF,
				result: 0xFFFFFFFF,
			},
			"test case #3": {
				number: 0x00FF00FF,
				result: 0xFF00FF00,
			},
			"test case #4": {
				number: 0x0000FFFF,
				result: 0xFFFF0000,
			},
			"test case #5": {
				number: 0x01020304,
				result: 0x04030201,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := ToLittleEndian(test.number)

				assert.Equal(t, test.result, result)
			})
		}
	})

	t.Run("uint64", func(t *testing.T) {
		tests := map[string]struct {
			number uint64
			result uint64
		}{
			"test case #1": {
				number: 0x0000000000000000,
				result: 0x0000000000000000,
			},
			"test case #2": {
				number: 0xFFFFFFFFFFFFFFFF,
				result: 0xFFFFFFFFFFFFFFFF,
			},
			"test case #3": {
				number: 0x00FF00FF00FF00FF,
				result: 0xFF00FF00FF00FF00,
			},
			"test case #4": {
				number: 0x00000000FFFFFFFF,
				result: 0xFFFFFFFF00000000,
			},
			"test case #5": {
				number: 0x0102030405060708,
				result: 0x0807060504030201,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := ToLittleEndian(test.number)

				assert.Equal(t, test.result, result)
			})
		}
	})

}
