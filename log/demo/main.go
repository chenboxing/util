package main

import (
	"time"

	"github.com/chenboxing/util/log"
)

//cmd to change level
//echo -n "help" | nc -4t  -w1 127.0.0.1 8765

//run it as follow: nohub ./demo > a &
func main() {
	log.SetFormatter("text")
	log.StartDaemon()

	for {
		log.Info("test")
		log.Fields{
			"name": "qianlnk",
			"age":  "27",
		}.Info("test fields")
		log.Error("test error")
		time.Sleep(time.Second * 2)
		test := make(map[string]interface{})
		test["lala"] = "haha"
		test["number"] = 123
		log.Fields{
			"aaa": "aaa",
			"bbb": "bbb",
		}.Add(log.Fields{
			"ddd": "ddd",
			"eee": "eee",
		}).Info("ccc")

		log.Fields{}.Add(log.Fields(test)).Del("number", "lala").Info("bababaa")
		log.Fields(test).Purpose(log.PpsPerformance).Error("testErr")
	}
}
