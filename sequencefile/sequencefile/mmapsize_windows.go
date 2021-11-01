package sequencefile

import (
	"fmt"
)

// Dwallocationgranularity is usually 64KB, you can get more info from links below.
// 	https://devblogs.microsoft.com/oldnewthing/20031008-00/?p=42223
// 	https://stackoverflow.com/questions/8583449/memory-map-file-offset-low
// 	https://docs.microsoft.com/en-us/windows/win32/memory/creating-a-view-within-a-file
const dwAllocationGranularity = 1 << 16 //64kb

func mmapSize(size int, floor bool) (int, error) {
	if size > MAX_FILE_SIZE {
		return 0, fmt.Errorf("mmap too large, the MAX size: %d, but got: %d", MAX_FILE_SIZE, size)
	}

	if size < 0 {
		return 0, fmt.Errorf("mmap requires non-negative size, but got: %d", size)
	}

	pageSize := int64(dwAllocationGranularity)
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
