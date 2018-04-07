package std

import (
	"bufio"
	"os"
)

type Input struct {
	*bufio.Scanner
	committed bool
}

func NewInput() *Input {
	return &Input{
		Scanner:   bufio.NewScanner(os.Stdin),
		committed: true,
	}
}

func (r *Input) Scan() bool {
	if !r.committed {
		return true
	}
	r.committed = false
	return r.Scanner.Scan()
}

func (r *Input) Commit() error {
	r.committed = true
	return nil
}
