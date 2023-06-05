package resp

import (
	"bufio"
	"io"
)

type Enconder struct {
	writer bufio.Writer
}

func NewEnconder(w io.Writer) *Enconder {
	return &Enconder{
		writer: *bufio.NewWriter(w),
	}
}
