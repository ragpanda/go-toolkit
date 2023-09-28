package util

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/log"
)

type SnowFlakeIDGenerator interface {
	GenerateID() int64
}

const (
	timestampBits  = 32
	machineIDBits  = 5
	sequenceBits   = 15
	maxMachineID   = (1 << machineIDBits) - 1
	maxSequence    = (1 << sequenceBits) - 1
	timeShift      = machineIDBits + sequenceBits
	machineIDShift = sequenceBits
)

type SmallSnowflake struct {
	machineID     int64
	lastTimestamp int64
	sequence      int64
	lock          sync.Mutex
	startTime     int64
}

/*
NewSmallSnowflake

Snowflake Structure:
--------------------
|  Timestamp  | Machine ID | Sequence |
---------------------------------------
|   32 bits   |   5 bits   | 15 bits  |
---------------------------------------

Some environments cannot use bd NTID, and a front-end friendly id is required.
*/
func NewSmallSnowflake(machineID int64) (*SmallSnowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("Invalid machine ID")
	}

	startTime, _ := time.Parse("2006-01-02", "2023-01-01")

	return &SmallSnowflake{
		machineID: machineID,
		startTime: startTime.Unix(),
	}, nil
}

func (s *SmallSnowflake) GenerateID() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	timestamp := time.Now().Unix() - s.startTime

	if timestamp < s.lastTimestamp {
		panic("Clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence

		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().Unix()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := (timestamp << timeShift) | (s.machineID << machineIDShift) | s.sequence
	if s.sequence+1 > maxSequence {
		log.Warn(context.Background(), "id generate reach limit, wait piece")
		time.Sleep(1 * time.Second)
	}
	return id
}
