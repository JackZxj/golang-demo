package pwd

import (
	"os"
	"path/filepath"
)

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	res, _ := filepath.EvalSymlinks(os.TempDir())
	return res
}
