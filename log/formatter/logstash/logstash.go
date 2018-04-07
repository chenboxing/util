package logstash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

// Formatter generates json in logstash format.
// Logstash site: http://logstash.net/
type LogstashFormatter struct {
	Type string // if not empty use for logstash type field.

	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

func (f *LogstashFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["@version"] = 1

	if f.TimestampFormat == "" {
		f.TimestampFormat = time.RFC3339Nano
	}

	entry.Data["@timestamp"] = entry.Time.Format(f.TimestampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		entry.Data["fields.message"] = v
	}
	entry.Data["message"] = entry.Message

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		entry.Data["fields.level"] = v
	}
	entry.Data["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			entry.Data["fields.type"] = v
		}
		entry.Data["type"] = f.Type
	}

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal %#v to JSON (%v)", entry.Data, err)
	}
	dst := recoverSpecialChar(serialized)
	return append(dst, '\n'), nil
}

func recoverSpecialChar(src []byte) []byte {
	dst := bytes.Replace(src, []byte("\\u003c"), []byte("<"), -1)
	dst = bytes.Replace(dst, []byte("\\u003e"), []byte(">"), -1)
	dst = bytes.Replace(dst, []byte("\\u0026"), []byte("&"), -1)
	return dst
}
