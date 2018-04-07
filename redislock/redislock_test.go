package redislock

import (
	"fmt"
	"testing"
	"time"

	redis "gopkg.in/redis.v5"
)

func TestLock(t *testing.T) {
	rd := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "127.0.0.1:6379",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  3 * time.Second,
		PoolSize:     10,
	})
	defer rd.Close()

	go func() {
		Alock := NewLock("xxxxx", 5*time.Second)
		err := Alock.Lock(rd) //5 秒后自动删除Alock
		fmt.Println("111", err)
		time.Sleep(7 * time.Second) //等待7秒
		err = Alock.Unlock(rd)      //想删除的是Alock锁，但是Alock 已经被自动删除 ,Block由于value 不一样，所以也不会删除
		fmt.Println(err)
	}()

	time.Sleep(6 * time.Second) //此时Alock 已经被删除
	Block := NewLock("xxxxx", 5*time.Second)
	err := Block.Lock(rd) //此时 会获取新的lock Block
	fmt.Println("222", err)

	time.Sleep(2 * time.Second)
	Clock := RedisLock{lockKey: "xxxxx"}
	err = Clock.Lock(rd) //想获取新的lock Clock，但由于 Block还存在，返回错误
	fmt.Println("333", err)

	time.Sleep(10 * time.Second)

}
