package main

import (
	"github.com/chenboxing/util/log/forwarder/redis"
	"github.com/chenboxing/util/log/forwarder/std"
)

func newInput(config *ioConfig) (input, error) {
	switch config.Type {
	case "redis":
		return redis.NewInput(&config.Redis)
	}
	return std.NewInput(), nil
}
