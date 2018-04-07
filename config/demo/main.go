package main

import (
	"fmt"
	"time"

	"github.com/chenboxing/util/config"
	"github.com/chenboxing/util/log"
)

type conf struct {
	X  string `yaml:"x"`
	Y  int    `yaml:"y"`
	AB string `yaml:"a_b"`
}

func main() {
	var conf conf

	fmt.Println(config.GetConfigAbsolutePath("config.yml"))
	config.Parse(&conf, config.GetConfigAbsolutePath("config.yml"))

	fmt.Printf("%#v\n", conf)
	fmt.Println("formatter:", log.GetFormatter())
	fmt.Println("release:", log.GetRelease())
	fmt.Println("mode:", log.GetMode())
	fmt.Println("level:", log.GetLevel())
	for {
		log.Fields{
			"test": "test message",
		}.Info("haha")
		log.Error("error")
		log.Warn("warn")
		log.Debug("debug")
		time.Sleep(time.Second * 2)
	}
}
