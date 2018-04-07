package cache_test

import (
	"fmt"
	"testing"

	"github.com/chenboxing/uline_clear/config"
	"github.com/chenboxing/util/cache"
)

func init() {
	cfg := config.Load()
	cache.Configure(
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Password,
		cfg.Redis.DBIndex,
		10)
}

type User struct {
	Name string
	Age  int
}

func TestCache(t *testing.T) {
	user := &User{
		Name: "qianlnk",
		Age:  27,
	}
	if err := cache.Set("lnkcache", user, 0); err != nil {
		fmt.Println(err)
		return
	}

	var res *User
	if err := cache.Get("lnkcache", &res); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func TestGet(t *testing.T) {
	var res string
	//cache.Set("abcd", 10, 0)
	err := cache.Get("withdraw_link_id", &res)
	fmt.Println(res, err)
}

func TestSet(t *testing.T) {
	cache.Set("can_withdraw", 1, -1)
}
