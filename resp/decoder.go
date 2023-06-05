package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type RedisType byte

const (
	String RedisType = '+'
	Error  RedisType = '-'
	Int    RedisType = ':'
	Bulk   RedisType = '$'
	Array  RedisType = '*'
)

type RESP struct {
	Type   RedisType
	Values any
}

type Decoder struct {
	reader bufio.Reader
}

func (redisType RedisType) String() string {
	switch redisType {
	case '+':
		return "String"
	case '-':
		return "Error"
	case ':':
		return "Int"
	case '$':
		return "Bulk"
	case '*':
		return "Array"
	case 'n':
		return "nil"
	}

	return "unknown"
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: *bufio.NewReader(r),
	}
}

func (decoder *Decoder) Decode() (RESP, error) {
	for {
		line, err := decoder.readLine()
		if err != nil {
			break
		}

		redisType := RedisType(line[0])
		switch redisType {
		case String:
			return RESP{
				Type:   redisType,
				Values: string(line[1:]),
			}, nil
		case Error:
			return RESP{
				Type:   redisType,
				Values: errors.New(string(line[1:])),
			}, nil
		case Int:
			intFromResp, err := strconv.Atoi(string(line[1:]))
			if err != nil {
				return RESP{}, err
			}
			return RESP{
				Type:   redisType,
				Values: intFromResp,
			}, nil
		case Bulk:
			return RESP{
				Type:   redisType,
				Values: string(line[1:]),
			}, nil
		case Array:
			numElements, err := strconv.Atoi(string(line[1:]))
			if err != nil {
				return RESP{}, err
			}
			if numElements == -1 {
				return RESP{
					Type:   'n',
					Values: nil,
				}, nil
			}
			values := make([]RESP, numElements)
			for i := 0; i < numElements; i++ {
				values[i], err = decoder.Decode()
				if err != nil {
					return RESP{}, err
				}
			}
			return RESP{
				Type:   redisType,
				Values: values,
			}, nil
		}
	}

	return RESP{}, nil
}

func (decoder *Decoder) readLine() ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := decoder.reader.ReadBytes('\n')
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
