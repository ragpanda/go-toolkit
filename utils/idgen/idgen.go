package util

import (
	"context"
	"fmt"
	"hash/fnv"
	"net"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/log"
)

type SnowFlakeIDGenerator interface {
	GenerateID() int64
}

type SnowFlakeOption struct {
	StartTime    time.Time
	TimestampBit int64
	MachineBit   int64
	SequenceBit  int64
	MachineID    int64
	TimeUnit     time.Duration
}

func NewSnowflake(ctx context.Context, option *SnowFlakeOption) (*Snowflake, error) {
	if option.TimestampBit < 32 {
		err := fmt.Errorf("timestamp bit must be greater than 32")
		log.Error(ctx, err.Error())
		return nil, err
	}

	if option.MachineBit < 4 {
		err := fmt.Errorf("machine bit must be greater than 4")
		log.Error(ctx, err.Error())
		return nil, err
	}
	if option.MachineID < 0 || option.MachineID > (1<<option.MachineBit)-1 {
		err := fmt.Errorf("machine id out of range")
		log.Error(ctx, err.Error())
		return nil, err
	}

	if option.SequenceBit < 4 {
		err := fmt.Errorf("sequence bit must be greater than 4")
		log.Error(ctx, err.Error())
		return nil, err
	}
	if option.TimestampBit+option.MachineBit+option.SequenceBit > 63 {
		err := fmt.Errorf("total bit must be less than 63")
		log.Error(ctx, err.Error())
		return nil, err
	}

	self := &Snowflake{}
	self.machineID = option.MachineID
	self.maxMachineID = (1 << option.MachineBit) - 1
	self.maxSequence = (1 << option.SequenceBit) - 1
	self.timeShift = option.MachineBit + option.SequenceBit
	self.machineIDShift = option.SequenceBit
	self.startTime = option.StartTime
	if option.TimeUnit > time.Second || option.TimeUnit < time.Millisecond {
		return nil, fmt.Errorf("time unit not support")
	} else {
		unit := int64(option.TimeUnit / time.Millisecond)
		self.getTime = func(t time.Time) int64 {

			return t.UnixMilli() / unit
		}
	}
	self.op = option

	return self, nil
}

type Snowflake struct {
	machineID     int64
	lastTimestamp int64
	sequence      int64

	lock      sync.Mutex
	startTime time.Time
	getTime   func(time.Time) int64

	maxMachineID   int64
	maxSequence    int64
	timeShift      int64
	machineIDShift int64

	op *SnowFlakeOption
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
func NewSmallSnowflake(machineID int64) (*Snowflake, error) {
	startTime, _ := time.Parse("2006-01-02", "2024-01-01")
	return NewSnowflake(context.Background(), &SnowFlakeOption{
		StartTime:    startTime,
		TimestampBit: 32,
		MachineBit:   5,
		SequenceBit:  15,
		MachineID:    machineID,
		TimeUnit:     time.Second,
	})
}

func NewStandardSnowflake(machineID int64) (*Snowflake, error) {
	startTime, _ := time.Parse("2006-01-02", "2024-01-01")
	return NewSnowflake(context.Background(), &SnowFlakeOption{
		StartTime:    startTime,
		TimestampBit: 33,
		MachineBit:   14,
		SequenceBit:  16,
		MachineID:    machineID,
		TimeUnit:     time.Second,
	})
}

func GenerateMachineIDByMac(machineBitLen int64) int64 {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(context.Background(), "Error getting network interfaces: %s", err.Error())
		return 0
	}

	// Create a new FNV-1a hash
	h := fnv.New64a()

	// Write the Mac bytes to the hash
	for _, iface := range interfaces {
		// 获取并打印每个网络接口的硬件地址（MAC 地址）
		mac := iface.HardwareAddr
		if mac == nil || len(mac) == 0 {
			log.Debug(context.Background(), "Interface Name: %s, MAC Address: %s\n", iface.Name, mac)
		} else {
			h.Write(mac)
		}
	}

	// Get the hash sum
	hashSum := h.Sum64()

	// Compress the hash sum to 10 digits
	machineID := int64(hashSum % (1<<machineBitLen - 1))

	log.Info(context.Background(), "Generated Machine ID: %d", machineID)
	return machineID
}

func (s *Snowflake) GenerateID() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	for {
		timestamp := s.getTime(time.Now()) - s.getTime(s.startTime)

		if timestamp < s.lastTimestamp {
			panic(fmt.Sprintf("clock moved backwards %d <- %d", timestamp, s.lastTimestamp))
		}

		if timestamp == s.lastTimestamp {
			s.sequence = (s.sequence + 1) & s.maxSequence
			if s.sequence == 0 {
				for timestamp <= s.lastTimestamp {
					timestamp = time.Now().Unix()
				}
			}
		} else {
			s.sequence = 0
		}

		s.lastTimestamp = timestamp

		id := (timestamp << s.timeShift) | (s.machineID << s.machineIDShift) | s.sequence
		if s.sequence+1 > s.maxSequence {
			log.Warn(context.Background(), "id generate reach limit %d, wait piece", s.maxSequence)
			time.Sleep(s.op.TimeUnit)
			continue
		}
		return id
	}

}
