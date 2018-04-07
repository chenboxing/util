//Package config for setting log config
//level: panic fatal error warn info debug
//formatter: logstash json text
package config

import (
	"strings"

	"github.com/chenboxing/util/log"
	"github.com/chenboxing/util/log/writer/redis"
)

type wirterType int

const (
	redisWriter wirterType = iota
	fileWriter
	stdWriter
)

func parseWriterType(tp string) (wirterType, error) {
	switch strings.ToLower(tp) {
	case "redis":
		return redisWriter, nil
	case "file":
		return fileWriter, nil
	default:
		return stdWriter, nil
	}
}

type config struct {
	Log struct {
		Mode      string `yaml:"mode"       toml:"mode"         json:"mode"`
		Level     string `yaml:"level"      toml:"level"        json:"level"`
		Formatter string `yaml:"formatter"  toml:"formatter"    json:"formatter"`
		Release   string `yaml:"release"    toml:"release"      json:"release"`
		Port      int    `yaml:"port"       toml:"port"         json:"port"`
		Writer    struct {
			Type  string       `yaml:"type"       toml:"type"         json:"type"`
			Redis redis.Config `yaml:"redis"      toml:"redis"        json:"redis"`
			Path  string       `yaml:"path"       toml:"path"         json:"path"`
		} `yaml:"writer"       toml:"writer"         json:"writer"`
	} `yaml:"log" toml:"log" json:"log"`
}

func setLogOptions(file string) error {
	var cfg config
	if err := parse(&cfg, file); err != nil {
		return err
	}
	log.SetFormatter(cfg.Log.Formatter)
	log.SetMode(cfg.Log.Mode)
	log.SetLevel(cfg.Log.Level)
	log.SetRelease(cfg.Log.Release)
	log.SetPort(cfg.Log.Port)

	wt, _ := parseWriterType(cfg.Log.Writer.Type)
	switch wt {
	case redisWriter:
		w, err := redis.New(&cfg.Log.Writer.Redis)
		if err != nil {
			return err
		}

		log.SetOutput(w)
	case fileWriter:
		if err := log.SetOutputPath(cfg.Log.Writer.Path); err != nil {
			return err
		}
	default:
		log.Debug("default to stdout writer")
	}

	log.Fields{
		"logConfig": cfg.Log,
	}.Debug("set log config")
	log.StartDaemon()
	return nil
}
