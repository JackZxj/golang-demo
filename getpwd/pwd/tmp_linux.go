package pwd

import (
	"os"
)

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	return os.TempDir()
}
