package data

import (
	"strconv"
)

func SCR_UnpackInt(data []byte) string {
	var result string

	for _, b := range data {
		// Remove the +32 offset
		val := int(b) - 32

		// Format the number back to a string.
		strVal := strconv.Itoa(val)

		// If it's not the very first chunk, we MUST pad it with a leading zero
		// if the integer was less than 10 (e.g., converting 5 back to "05").
		if len(result) > 0 && val < 10 {
			strVal = "0" + strVal
		}

		result += strVal
	}

	return result
}

func SCR_PackInt(num int) []byte {
	value := strconv.Itoa(num)

	length := len(value)
	if length == 0 {
		return []byte{}
	}

	var result []byte

	// v3 = 2 - ((v9 & 1) != 0)
	chunkSize := 2
	if length%2 != 0 {
		chunkSize = 1
	}

	for i := 0; i < length; {
		// strncpy
		chunk := value[i : i+chunkSize]

		// atol
		val, _ := strconv.Atoi(chunk)

		// Add 32 and append to buffer
		result = append(result, byte(val+32))

		i += chunkSize
		chunkSize = 2 // All subsequent chunks are size 2.
	}

	return result
}
