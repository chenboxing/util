package main

import "github.com/chenboxing/util/log/forwarder/redis"

type ioConfig struct {
	Type  string       `yaml:"type"  toml:"type"  json:"type"`
	Redis redis.Config `yaml:"redis" toml:"redis" toml:"redis"`
}

type input interface {
	Scan() bool
	Err() error
	Bytes() []byte
	Commit() error
}

type output interface {
	Send([]byte) error
}

type forwarder struct {
	input  input
	output output
}

func newForwarder(conf *conf) (*forwarder, error) {
	input, err := newInput(&conf.Input)
	if err != nil {
		return nil, err
	}
	output, err := newOutput(&conf.Output)
	if err != nil {
		return nil, err
	}
	return &forwarder{
		input:  input,
		output: output,
	}, nil
}

func (f *forwarder) Forward() error {
	for f.input.Scan() {
		if err := f.output.Send(f.input.Bytes()); err != nil {
			return err
		}
		if err := f.input.Commit(); err != nil {
			return err
		}
	}
	return f.input.Err()
}
