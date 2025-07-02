package internal

import (
	"URLShortener/utils"
	"testing"
)

func TestMachineIdValidation(t *testing.T) {
	validMachineIds := []int{0, 3, 4095}
	invalidMachineIds := []int{4096, 10000, -250}

	// Test valid machine IDs
	for _, machineId := range validMachineIds {
		if !utils.ValidateMachineId(machineId) {
			t.Errorf("Machine ID %d should be valid but was not", machineId)
		}
	}

	// Test invalid machine IDs
	for _, machineId := range invalidMachineIds {
		if utils.ValidateMachineId(machineId) {
			t.Errorf("Machine ID %d should be invalid but was marked valid", machineId)
		}
	}
}
