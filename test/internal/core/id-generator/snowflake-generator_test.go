package id_generator

import (
	idgenerator "URLShortener/internal/core/id-generator"
	"sync"
	"testing"
	"time"
)

// TestSnowflakeGeneratorBasic tracks basic id generation characteristics
func TestSnowflakeGeneratorBasic(t *testing.T) {
	generator := idgenerator.SnowflakeGenerator{
		MachineID: 14,
	}

	now := time.Now().UnixMilli()
	id := generator.NextId()

	// Assert the encoded timestamp is within 1 millisecond of its request time
	timestamp := int64((id >> 24) & ((1 << 39) - 1))
	if timestamp-now > 1 {
		t.Errorf("Expected timestamp closer to %d, got delta of %d for encoded value %d", now, timestamp-now, timestamp)
	}

	// Assert the machineId is encoded properly
	machineId := int((id >> 12) & ((1 << 12) - 1))
	if machineId != 14 {
		t.Errorf("Expected machineId to be 14, got %d", machineId)
	}

	// Assert the sequencenum is encoded correctly
	sequencenum := int64(id & (1 << 12))
	if sequencenum != 0 {
		t.Errorf("Expected sequence number to be 0, got %d", sequencenum)
	}
}

// TestSnowflakeGeneratorSequenceNumberIncrements tests whether multiple generations close together
//
//	generate ids with the same timestamp, but different sequence numbers, to maintain uniqueness
//	under high load. This does not test rollover (at >4096 requests in one ms), just that unique
//	sequence numbers are produced.
func TestSnowflakeGeneratorSequenceNumberIncrements(t *testing.T) {
	// Initialize the SnowflakeGenerator
	generator := idgenerator.SnowflakeGenerator{
		MachineID:       1,
		SequenceNum:     0,
		LastRequestTime: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	var wg sync.WaitGroup
	numTests := 1000
	uniqueSequenceNumbers := make(map[int]bool) // To store unique sequence numbers
	var mu sync.Mutex                           // Mutex to synchronize access to the map
	timestamp := uint64(0)

	// Wait for the next millisecond boundary to align the test execution
	waitUntilNextMillisecond()

	// Run the nextId() method 1000 times concurrently
	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			id := generator.NextId()

			// Extract the sequence number
			extractedSequence := int(id & 0xFFF)

			// Lock the map to ensure thread-safe access
			mu.Lock()
			// Check if this sequence number already exists
			if uniqueSequenceNumbers[extractedSequence] {
				t.Errorf("Test failed: Duplicate sequence number %d found at iteration %d", extractedSequence, i)
				mu.Unlock()
				return
			}

			// Add the sequence number to the map
			uniqueSequenceNumbers[extractedSequence] = true
			mu.Unlock()

			extractedTimestamp := id >> 24
			if timestamp == 0 {
				timestamp = extractedTimestamp
			} else if extractedTimestamp != timestamp {
				t.Errorf("Test failed: Timestamp does not match %d for timestamps %d", extractedTimestamp, timestamp)

			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Ensure there are exactly 4096 unique sequence numbers
	if len(uniqueSequenceNumbers) != numTests {
		t.Errorf("Test failed: Expected %d unique sequence numbers, but got %d", numTests, len(uniqueSequenceNumbers))
	} else {
		t.Logf("Test passed: All %d sequence numbers are unique, and all timestamps are the same.", numTests)
	}
}

func waitUntilNextMillisecond() {
	// Get the current time
	now := time.Now()

	// Calculate how much time is left until the next millisecond
	timeToWait := time.Millisecond - time.Duration(now.Nanosecond())%time.Millisecond

	// Wait until the next millisecond boundary
	time.Sleep(timeToWait)
}
