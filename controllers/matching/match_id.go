package matching

import (
	"sync"
	"time"
)

type Snowflake struct {
	mu       sync.Mutex
	lastTs   int64
	sequence int64
	workerID int64
}

const (
	epoch        int64 = 1672531200000 // 2023-01-01
	workerBits         = 5
	sequenceBits       = 12

	maxWorkerID = -1 ^ (-1 << workerBits)
	maxSequence = -1 ^ (-1 << sequenceBits)

	workerShift = sequenceBits
	timeShift   = sequenceBits + workerBits
)

func NewSnowflake(workerID int64) *Snowflake {
	if workerID < 0 || workerID > maxWorkerID {
		panic("invalid workerID")
	}
	return &Snowflake{
		lastTs:   -1,
		workerID: workerID,
	}
}

func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	if now == s.lastTs {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTs {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTs = now

	return ((now - epoch) << timeShift) |
		(s.workerID << workerShift) |
		s.sequence
}
