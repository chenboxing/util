package redis

import (
	"bufio"
	"io"
	"os"

	"github.com/chenboxing/util/log/forwarder/redis"
)

type Config struct {
	Addr     string `yaml:"addr"     toml:"addr"     json:"addr"`
	Password string `yaml:"password" toml:"password" json:"password"`
	ListKey  string `yaml:"list_key" toml:"list_key" json:"list_key"`
}

func New(config *Config) (io.Writer, error) {
	output, err := redis.NewOutput(&redis.Config{
		Addr:     config.Addr,
		Password: config.Password,
		ListKey:  config.ListKey,
	})
	if err != nil {
		return nil, err
	}
	r, w := io.Pipe()
	writer := &redisWriter{
		reader: r,
		writer: w,
		output: output,
	}
	go writer.Forward()
	return writer, nil
}

type redisWriter struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	output *redis.Output
}

func (r *redisWriter) Close() error {
	if r.reader != nil {
		r.reader.Close()
		r.reader = nil
	}
	if r.writer != nil {
		r.writer.Close()
		r.writer = nil
	}
	if r.output != nil {
		r.output.Close()
		r.output = nil
	}
	return nil
}

func (r *redisWriter) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (w *redisWriter) Forward() {
	s := bufio.NewScanner(w.reader)
	for s.Scan() {
		//os.Stderr.Write(s.Bytes())
		//os.Stderr.Write([]byte{'\n'})
		if err := w.output.Send(s.Bytes()); err != nil {
			os.Stderr.Write([]byte("Err=>" + err.Error() + "log=>"))
			os.Stderr.Write(s.Bytes())
			os.Stderr.Write([]byte{'\n'})
		}
	}
}
