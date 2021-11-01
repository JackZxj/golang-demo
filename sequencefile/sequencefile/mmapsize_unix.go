// +build !windows

package sequencefile

import (
	"fmt"
	"os"
)

func mmapSize(size int, floor bool) (int, error) {
	if size > MAX_FILE_SIZE {
		return 0, fmt.Errorf("mmap too large, the MAX size: %d, but got: %d", MAX_FILE_SIZE, size)
	}

	if size < 0 {
		return 0, fmt.Errorf("mmap requires non-negative size, but got: %d", size)
	}

	pageSize := int64(os.Getpagesize())
	if size == 0 {
		if floor {
			return 0, nil
		}
		return int(pageSize), nil
	}
	sz := int64(size)
	if (sz % pageSize) != 0 {
		if floor {
			sz = (sz / pageSize) * pageSize
		} else {
			sz = ((sz / pageSize) + 1) * pageSize
		}
	}
	return int(sz), nil
}
