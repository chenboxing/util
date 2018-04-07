package std

import (
	"os"
)

type Output struct{}

func NewOutput() *Output {
	return &Output{}
}

func (o *Output) Send(item []byte) error {
	item = append(item, '\n')
	_, err := os.Stdout.Write(item)
	return err
}
