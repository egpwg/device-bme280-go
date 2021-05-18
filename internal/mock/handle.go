package mock

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// SensorHandle ip&mac handle
type SensorHandle struct {
	Locker        sync.Mutex
	SensorInfoMap []float32
}

// Scout scout ip and mac
// input command string
func Scout(command string) ([]float32, error) {

	return scout(strings.Split(command, " "))
}

func scout(command []string) ([]float32, error) {
	if len(command) < 2 {
		panic("comman len ERROR")
	}
	result := []float32{}

	cmd := exec.Command(command[0], command[1:]...)
	consoleOut, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	outStr := string(consoleOut)
	kvArr := strings.Split(outStr, "\n")
	for i, kv := range kvArr {
		if i == 0 {
			continue
		}
		arr := strings.Split(kv, ":")
		if len(arr) != 2 {
			continue
		}
		vTemp := strings.Trim(arr[1], " ")
		vTrue, err := strconv.ParseFloat(vTemp, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, float32(vTrue))
	}

	return result, nil
}
