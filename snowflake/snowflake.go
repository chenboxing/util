package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch = int64(1475251200000) //2016年10月1日

	defaultWorkID       = int64(1)
	defaultDataCenterID = int64(1)
	defaultComputerID   = int64(1)
	defaultSmall        = false
)

// 40bit		5bit			 	2bit				5bit 				4bit			8bit
//|  timestamp | worker_id_bits | data_center_id_bits | computer_id_bits  | retain_id_bits |  sequence_bits |

type Snowflake struct {
	small bool //用于结算主键，16位

	workID       int64
	dataCenterID int64
	computerID   int64
	retainID     int64
	sequence     int64

	workIDBits       uint
	dataCenterIDBits uint
	computerIDBits   uint
	retainIDBits     uint
	sequenceBits     uint

	maxWorkID       int64
	maxDataCenterID int64
	maxComputerID   int64
	maxRetainID     int64
	maxSequence     int64

	workIDShift        uint
	dataCenterIDShift  uint
	computerIDShift    uint
	retainIDShift      uint
	timestampLeftShift uint

	epoch         int64
	lastTimestamp int64

	mu sync.Mutex
}

func maxID(bits uint) int64 {
	if bits == 0 {
		return 0
	}

	return -1 ^ (-1 << bits)
}

func New(workID int64) *Snowflake {
	sf := &Snowflake{
		workID:           workID,
		small:            defaultSmall,
		dataCenterID:     defaultDataCenterID,
		computerID:       defaultComputerID,
		workIDBits:       5,
		dataCenterIDBits: 2,
		computerIDBits:   5,
		retainIDBits:     4,
		sequenceBits:     8,
		epoch:            epoch,
		lastTimestamp:    -1,
	}

	sf.setMaxAndShift()

	if 0 > sf.workID || sf.workID > sf.maxWorkID {
		panic("workID out of range")
	}
	return sf
}
func (s *Snowflake) SetWorkID(id int64) *Snowflake {
	if s.small {
		s.workID = 0
	}

	s.workID = id

	if 0 > s.workID || s.workID > s.maxWorkID {
		panic("workID out of range")
	}

	return s
}

func (s *Snowflake) SetDataCenterID(id int64) *Snowflake {
	if s.small {
		s.dataCenterID = 0
	}

	s.dataCenterID = id

	if 0 > s.dataCenterID || s.dataCenterID > s.maxDataCenterID {
		panic("dataCenterID out of range")
	}

	return s
}

func (s *Snowflake) SetComputerID(id int64) *Snowflake {
	if s.small {
		s.computerID = 0
	}

	s.computerID = id

	if 0 > s.computerID || s.computerID > s.maxComputerID {
		panic("computerID out of range")
	}

	return s
}

func (s *Snowflake) SetRetainID(id int64) *Snowflake {
	if s.small {
		s.retainID = 0
	}

	s.retainID = id

	if 0 > s.retainID || s.retainID > s.maxRetainID {
		panic("retainID out of range")
	}

	return s
}

func (s *Snowflake) Small() *Snowflake {
	s.small = true
	s.dataCenterID = 0
	s.computerID = 0
	s.retainID = 0

	s.workIDBits = 5
	s.dataCenterIDBits = 0
	s.computerIDBits = 0
	s.retainIDBits = 0
	s.sequenceBits = 7

	s.setMaxAndShift()

	return s
}

func (s *Snowflake) setMaxAndShift() {
	s.maxWorkID = maxID(s.workIDBits)
	s.maxDataCenterID = maxID(s.dataCenterIDBits)
	s.maxComputerID = maxID(s.computerIDBits)
	s.maxRetainID = maxID(s.retainIDBits)
	s.maxSequence = maxID(s.sequenceBits)

	s.retainIDShift = s.sequenceBits
	s.computerIDShift = s.retainIDShift + s.retainIDBits
	s.dataCenterIDShift = s.computerIDShift + s.computerIDBits
	s.workIDShift = s.dataCenterIDShift + s.dataCenterIDBits
	s.timestampLeftShift = s.workIDShift + s.workIDBits
}

func (s *Snowflake) RestoreDate(id int64) (timestamp int64, workID, dataCenterID, computerID int64) {
	timestamp = int64(id >> s.timestampLeftShift)
	timestamp += epoch
	//timestamp /= 1000

	workID = (id >> s.workIDShift) & maxID(s.workIDBits)
	dataCenterID = (id >> s.dataCenterIDShift) & maxID(s.dataCenterIDBits)
	computerID = (id >> s.computerIDShift) & maxID(s.computerIDBits)

	return
}

func (s *Snowflake) Boundary(date time.Time, offset int64) (start int64, end int64) {
	startDate, err := time.Parse("2006-01-02", date.Format("2006-01-02"))
	if err != nil {
		return
	}

	start = startDate.Unix() * 1000
	end = start + 3600*24*1000 - offset

	shift := func(t int64) int64 {
		return ((t - epoch) << s.timestampLeftShift) |
			maxID(s.workIDBits) |
			maxID(s.dataCenterIDBits) |
			maxID(s.computerIDBits) |
			maxID(s.retainIDBits) | 1
	}

	return shift(start), shift(end)
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

RENEXT:
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	if s.lastTimestamp > timestamp {
		return 0, errors.New("Clock moved backwards")
	}

	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & s.maxSequence
		if s.sequence == 0 { //1毫秒内生成超过128个id则等下一毫秒
			time.Sleep(time.Millisecond)
			goto RENEXT
			//s.sequence = -1 & s.maxSequence
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return ((timestamp - epoch) << s.timestampLeftShift) |
		(s.workID << s.workIDShift) |
		(s.dataCenterID << s.dataCenterIDShift) |
		(s.computerID << s.computerIDShift) |
		(s.retainID << s.retainIDShift) |
		s.sequence, nil
}
