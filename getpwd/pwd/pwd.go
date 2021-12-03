package pwd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func Run() string {
	// 获取二进制所在的目录
	pwd := getCurrentAbPathByExecutable()
	path := filepath.Join(pwd, "go.mod")
	// 判断二进制目录是否是临时目录
	if strings.Contains(pwd, getTmpDir()) {
		// 是的话获取代码所在目录
		pwd = getCurrentAbPathByCaller()
		path = filepath.Join(pwd, "..", "go.mod")
	}
	str, _ := ioutil.ReadFile(path)
	return string(str)
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
