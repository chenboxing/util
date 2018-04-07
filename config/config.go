//Package config support json yaml(yml) toml
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"

)

//Helper is an obj of error
type Helper struct {
	Message string
}

//Error is the method of error
func (h *Helper) Error() string {
	return h.Message
}

func fileExists(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}

	f.Close()
	return true
}

//GetConfigAbsolutePath find config absolute path, for go test
func GetConfigAbsolutePath(file string) string {
	app := path.Dir(os.Args[0])
	if strings.HasPrefix(app, os.TempDir()) {
		return getConfigAbsolutePathForTest(file)
	}

	return getConfigAbsolutePathForBase(file)
}

func getConfigAbsolutePathForBase(file string) string {
	app := path.Base(os.Args[0])
	for _, dir := range []string{
		"",
		"config",
		"/etc/" + app,
		path.Join(os.Getenv("HOME"), "."+app),
	} {
		cf := path.Join(dir, file)
		if fileExists(cf) {
			return cf
		}
	}

	return ""
}

func getConfigAbsolutePathForTest(file string) string {
	_, filename, _, _ := runtime.Caller(2)
	dir := path.Dir(filename)
	for {
		for _, d := range []string{"", "config"} {
			cf := path.Join(dir, d, file)
			if fileExists(cf) {
				return cf
			}
		}
		dir = path.Dir(strings.TrimRight(dir, "/"))
		if dir == "/" {
			break
		}
	}

	return file
}

//Parse load config file and parse config
func Parse(cfg interface{}, file string) error {
	// set log option before any logging
	if err := setLogOptions(file); err != nil {
		return err
	}
	return parse(cfg, file)
}

func parse(cfg interface{}, file string) error {
	err := load(cfg, file)
	if h, ok := err.(*Helper); ok {
		fmt.Println(h.Error())
		os.Exit(1)
	}
	return err
}

func load(cfg interface{}, file string) error {
	err := parseFile(cfg, file)
	if err != nil {
		return err
	}



	parser := flags.NewParser(cfg, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return &Helper{e.Message}
		}
		return err
	}
	return nil
}

func parseFile(cfg interface{}, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	switch path.Ext(file) {
	case ".json":
		return json.NewDecoder(f).Decode(cfg)
	case ".yaml", ".yml":
		in, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(in, cfg)
	case ".toml":
		_, err := toml.DecodeReader(f, &cfg)
		return err
	default:
		return fmt.Errorf("unsupported config file format: %s", file)
	}
}
