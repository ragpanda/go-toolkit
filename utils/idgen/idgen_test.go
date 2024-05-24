package util

import (
	"sync"
	"testing"
	"time"
)

func TestSnowflakeGenerateID(t *testing.T) {
	snowflake, err := NewSmallSnowflake(1)
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	id := snowflake.GenerateID()
	t.Logf("gen id %d", id)
	// 检查生成的ID是否符合要求
	if id < 0 || id >= (1<<52) {
		t.Errorf("Generated ID is out of range")
	}
}

func TestSnowflakeIDIncrease(t *testing.T) {
	snowflake, err := NewSmallSnowflake(1)
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	lastID := int64(0)
	for i := 0; i < 100000; i++ {
		id := snowflake.GenerateID()
		t.Logf("gen id %d", id)
		// 检查生成的ID是否符合要求
		if id < 0 || id >= (1<<52) {
			t.Errorf("Generated ID is out of range")
			return
		}

		if id <= lastID {
			t.Errorf("Generated ID is not increasing: %d <= %d", id, lastID)
			return
		}
		lastID = id
	}

}

func TestSnowflakeGenerateStdID(t *testing.T) {
	snowflake, err := NewStandardSnowflake(GenerateMachineIDByMac(14))
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	id := snowflake.GenerateID()
	t.Logf("gen id %d", id)
	// 检查生成的ID是否符合要求
	if id < 0 || id >= (1<<63-1) {
		t.Errorf("Generated ID is out of range")
	}
}

func TestSnowflakeIDStdIncrease(t *testing.T) {
	snowflake, err := NewStandardSnowflake(GenerateMachineIDByMac(14))
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	lastID := int64(0)
	for i := 0; i < 300000; i++ {
		id := snowflake.GenerateID()
		//t.Logf("gen id %d", id)
		// 检查生成的ID是否符合要求
		if id < 0 || id >= (1<<63-1) {
			t.Errorf("Generated ID is out of range")
			return
		}

		if id <= lastID {
			t.Errorf("Generated ID is not increasing: %d <= %d", id, lastID)
			return
		}
		lastID = id
	}

}

func TestSnowflakeIDIncreaseWaitTime(t *testing.T) {
	snowflake, err := NewSmallSnowflake(1)
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	lastID := int64(0)
	for i := 0; i < 1000; i++ {
		id := snowflake.GenerateID()
		t.Logf("gen id %d", id)
		// 检查生成的ID是否符合要求
		if id < 0 || id >= (1<<52) {
			t.Errorf("Generated ID is out of range")
			return
		}

		if id <= lastID {
			t.Errorf("Generated ID is not increasing: %d <= %d", id, lastID)
			return
		}
		lastID = id
		time.Sleep(1 * time.Millisecond)
	}

}

func TestSnowflakeConcurrency(t *testing.T) {
	snowflake, err := NewSmallSnowflake(1)
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	// 启动多个 goroutine 并发生成ID
	numWorkers := 100
	done := make(chan struct{})

	ids := make(chan int64)
	defer close(ids)

	// 并发生成ID
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case ids <- snowflake.GenerateID():

				}
			}
		}()
	}

	// 检查生成的ID是否有重复
	seen := make(map[int64]bool)
	for i := 0; i < numWorkers*1000; i++ {
		id := <-ids
		if seen[id] {
			t.Errorf("Duplicate ID generated: %d", id)
		} else {
			t.Logf("gen id %d", id)
		}

		seen[id] = true
	}
	close(done)
	wg.Wait()
}

func TestSnowflakeInvalidMachineID(t *testing.T) {
	_, err := NewSmallSnowflake(-1)
	if err == nil {
		t.Errorf("Expected error for invalid machine ID")
	}
}

func TestSnowflakeClockBackwards(t *testing.T) {
	snowflake, err := NewSmallSnowflake(1)
	if err != nil {
		t.Errorf("Error creating Snowflake: %s", err)
	}

	// 模拟时钟回拨
	_ = snowflake.GenerateID()
	snowflake.lastTimestamp = time.Now().UnixNano()/1000000 + 1000

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for clock backwards")
		} else {
			t.Log(r)
		}

	}()

	_ = snowflake.GenerateID()
}
