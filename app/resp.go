package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}

	return []Value{}
}

func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}

	return ""
}

func DeserializeRESP(buffStream *bufio.Reader) (Value, error) {
	typeByte, err := buffStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(typeByte) {
	case "+":
		return DeserializeSimpleString(buffStream)
	case "$":
		return DeserializeBulkString(buffStream)
	case "*":
		return DeserializeArray(buffStream)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(typeByte))
}

func DeserializeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func DeserializeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func DeserializeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	var array []Value

	for i := 1; i <= count; i++ {
		value, err := DeserializeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	var readBytes []byte

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}

	}

	return readBytes[:len(readBytes)-2], nil
}
