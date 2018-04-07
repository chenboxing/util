package log_test

import (
	"testing"

	"github.com/chenboxing/util/log"
)

func TestLog(t *testing.T) {
	log.SetFormatter("logstash")
	log.SetLevel(log.DebugLevel)
	log.SetMode(log.Production)
	log.SetRelease("1.2.3")
	log.Fields{
		"dfdfd": 123.45,
		"key2":  []string{"value2", "sdff"},
		"key5":  "value2",
		"key1":  "value2",
	}.Infof("hello world %d", 123)

	log.Debugf("%d - %d - %s", 1, 2, "dfdf")
}

func TestRedis(t *testing.T) {
	//log.SetOutput()
}
