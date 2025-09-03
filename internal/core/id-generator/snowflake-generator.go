package id_generator

import (
	"sync"
	"time"
)

type SnowflakeGenerator struct {
	MachineID       int        // a value unique to the server processing the request
	SequenceNum     int        // a value that increments when multiple requests are processed quickly
	LastRequestTime uint64     // the last millisecond that a request occurred on, for SN incrementing
	mutex           sync.Mutex // thread lock for concurrent operations
}

func NewSnowflakeGenerator(machineId int) *SnowflakeGenerator {
	return &SnowflakeGenerator{MachineID: machineId}
}

// NextId produces a Snowflake algorithm id
// an ID has the following form:
// 64 bits: &++++++++++++++++++++++++++++++++++++++$$$$$$$$$$$$############ where:
//
//	& - sign bit, always 0
//	+ - 39 bits to represent timestamp
//	$ - 12 bits to represent machine code
//	# - 12 bits to represent sequence numbers
func (gen *SnowflakeGenerator) NextId() uint64 {
	var resId uint64 = 0
	gen.mutex.Lock()
	defer gen.mutex.Unlock()

	// verify request timing and adjust sequence number
	requestTime := uint64(time.Now().UnixMilli())
	if requestTime > gen.LastRequestTime {
		gen.SequenceNum = 0
	} else {
		gen.SequenceNum = gen.SequenceNum + 1
	}
	gen.LastRequestTime = requestTime

	// If there has been more than 4096 requests in the same millisecond, wait a millisecond
	if gen.SequenceNum > 4096 {
		gen.SequenceNum = 0
		requestTime++
		time.Sleep(1 * time.Millisecond)
	}

	// Construct the ID
	resId |= requestTime << 24
	resId |= uint64(gen.MachineID) << 12
	resId |= uint64(gen.SequenceNum)
	resId &= 0x7fffffffffffffff // force MSB to zero

	return resId
}
