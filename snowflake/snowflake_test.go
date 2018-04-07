package snowflake

import (
	"fmt"
	"testing"
	"time"
)

func TestMaxID(t *testing.T) {
	fmt.Println(-1 ^ (-1 << 8))
	fmt.Println(1 << 8)
}

func TestTimestamp(t *testing.T) {
	date := time.Now()
	fmt.Println(date.Unix())
	startDate, err := time.Parse("2006-01-02", date.Format("2006-01-02"))
	if err != nil {
		return
	}
	endDate := startDate.Add(time.Hour * 24)
	fmt.Println(startDate, endDate)
	fmt.Println(startDate.Unix(), endDate.Unix())
	start := startDate.Unix() * 1000
	end := start + 3600*24*1000

	fmt.Println(start, end)
}

func TestSnowflake(t *testing.T) {
	sf := New(3).SetDataCenterID(2).SetComputerID(2)
	fmt.Println(sf.Boundary(time.Now(), 0))
	for i := 0; i < 10; i++ {
		id, err := sf.NextID()
		fmt.Println("*", id, err)

		fmt.Println(sf.RestoreDate(id))
	}
}

//go test -bench="." -run BenchmarkSnowflake
func BenchmarkSnowflake(b *testing.B) {
	go func() {
		var sf = New(1)
		for i := 0; i < b.N; i++ {
			sf.NextID()
		}
	}()

	go func() {
		var sf = New(2)
		for i := 0; i < b.N; i++ {
			sf.NextID()
		}
	}()

	var sf = New(3)
	for i := 0; i < b.N; i++ {
		sf.NextID()
	}
}
