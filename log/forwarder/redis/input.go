package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Input struct {
	c       redis.Conn
	listKey string
	err     error
	buf     []byte
}

func NewInput(config *Config) (*Input, error) {
	conn, err := config.Connect()
	if err != nil {
		return nil, err
	}
	return &Input{
		c:       conn,
		listKey: config.ListKey,
	}, nil
}

func (r *Input) Bytes() []byte {
	return r.buf
}

func (r *Input) Scan() bool {
	for {
		buf, err := redis.Bytes(r.c.Do("LINDEX", r.listKey, 0))
		if err != nil {
			if err != redis.ErrNil {
				r.err = err
				return false
			}
			time.Sleep(time.Second)
			continue
		}
		r.buf = buf
		return true
	}
}

func (r *Input) Err() error {
	return r.err
}

func (r *Input) Commit() error {
	_, err := r.c.Do("LTRIM", r.listKey, 1, -1)
	if err != nil {
		r.err = err
	}
	return err
}
