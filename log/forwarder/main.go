package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chenboxing/util/config"
	"github.com/chenboxing/util/log/forwarder/redis"
)

type conf struct {
	Input  ioConfig `yaml:"input"  toml:"input"  json:"input"`
	Output ioConfig `yaml:"output" toml:"output" json:"output"`
}

var defaultConf = conf{
	Input: ioConfig{
		Type: "redis",
		Redis: redis.Config{
			Addr:     "127.0.0.1:6379",
			Password: "",
			ListKey:  "log_key",
		},
	},
	Output: ioConfig{
		Type: "std",
	},
}

func main() {
	var conf = defaultConf
	checkError(config.Parse(&conf, "forwarder.yaml"))
	forwarder, err := newForwarder(&conf)
	checkError(err)
	checkError(forwarder.Forward())
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		time.Sleep(time.Second)
		os.Exit(-1)
	}
}
