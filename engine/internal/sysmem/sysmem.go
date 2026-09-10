package sysmem

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	memoryFactor = 0.6
)

func TotalMemory() (int64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}

		var memKB int64
		_, err = fmt.Sscanf(line, "MemTotal: %d kB", &memKB)
		if err == nil {
			totalBytes := memKB * 1024
			safeLimit := float64(totalBytes) * memoryFactor
			return int64(safeLimit), nil
		}
	}

	return 0, errors.New("MemTotal not found in /proc/meminfo")
}
