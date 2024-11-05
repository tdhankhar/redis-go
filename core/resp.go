package core

import (
	"errors"
	"fmt"
)

func readInt64(data []byte) (int64, int, error) {
	i := 1
	var num int64 = 0
	for ; data[i] != '\r'; i++ {
		num = num * 10 + int64(data[i] - '0')
	}
	return num, i + 2, nil
}

func readSimpleString(data []byte) (string, int, error) {
	i := 1
	for ; data[i] != '\r'; i++ {}
	return string(data[1:i]), i + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	bytes, i, _ := readInt64(data)
	return string(data[i:i + int(bytes)]), i + int(bytes) + 2, nil
}

func readArray(data []byte) ([]interface{}, int, error) {
	length, i, _ := readInt64(data)
	elems := make([]interface{}, length)
	for j := 0; j < int(length); j++ {
		elem, k, err := DecodeOne(data[i:])
		if err != nil {
			return nil, 0, err
		}
		elems[j] = elem
		i += k
	}
	return elems, i, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case ':':
		return readInt64(data)
	case '+':
		return readSimpleString(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	case '-':
		return readError(data)
	}
	return nil, 0, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	elem, err := Decode(data)
	if err != nil {
		return nil, err
	}
	elems := elem.([]interface{})
	tokens := make([]string, len(elems))
	for i := range tokens {
		tokens[i] = elems[i].(string)
	}
	return tokens, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	elem, _, err := DecodeOne(data)
	return elem, err
}

func Encode(data interface{}) []byte {
	switch value := data.(type) {
	case string:
		return []byte(fmt.Sprintf("+%s\r\n", value))
	case BulkString:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", value))
	default:
		return []byte{}
	}
}