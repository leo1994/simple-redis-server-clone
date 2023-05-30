package main

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	stringType = '+'
	errorType  = '-'
	intType    = ':'
	bulkType   = '$'
	arrayType  = '*'
)

type RedisType string

const (
	String RedisType = "string"
	Error  RedisType = "error"
	Int    RedisType = "int"
	Bulk   RedisType = "bulk"
	Array  RedisType = "array"
)

type RedisValue struct {
	Type  RedisType
	Value any
}

func parseRedisProtocol(reader *bufio.Reader) (RedisValue, error) {
	typ, err := reader.ReadByte()
	if err != nil {
		return RedisValue{}, err
	}
	switch typ {
	case bulkType:
		return parseBulkString(reader)
	case arrayType:
		return parseArray(reader)
	}
	return RedisValue{}, fmt.Errorf("type %s is not defined", string(typ))
}

func parseArray(reader *bufio.Reader) (RedisValue, error) {
	delim, _ := readUntilCRLF(reader)
	arrLen, _ := strconv.Atoi(delim)

	commands := RedisValue{
		Type:  Array,
		Value: make([]RedisValue, arrLen),
	}
	for i := 0; i < arrLen; i++ {
		command, _ := parseRedisProtocol(reader)
		commands.Value.([]RedisValue)[i] = command
	}
	return commands, nil
}

func parseBulkString(reader *bufio.Reader) (RedisValue, error) {
	delim, _ := readUntilCRLF(reader)
	strLen, _ := strconv.Atoi(delim)

	buf := make([]byte, strLen+2)
	if _, err := reader.Read(buf); err != nil {
		return RedisValue{}, err
	}

	return RedisValue{
		Type:  Bulk,
		Value: removeCRLF(string(buf)),
	}, nil
}

func removeCRLF(s string) string {
	return s[:len(s)-2]
}

func readUntilCRLF(reader *bufio.Reader) (string, error) {
	buf, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return removeCRLF(string(buf)), nil
}
