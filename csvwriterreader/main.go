package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	writeCSV()
	fmt.Println("----- spliter -----")
	readCSV()
}

func writeCSV() {
	f, err := os.Create("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") //写入UTF-8 BOM,此处如果不写入就可能导致写入的汉字乱码
	w := csv.NewWriter(f)
	w.Write([]string{"id", "姓名", "分数"})
	w.Write([]string{"1", "张三", "90"})
	w.Write([]string{"2", "李四", "100"})
	w.Write([]string{"3", "王五", "79"})
	w.Write([]string{"4", "赵六", "82"})
	w.Flush() // 将 writer 缓冲中的数据都推送到 csv 文件，至此就完成了数据写入到 csv 文件
}

func readCSV() {
	f, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// 循环迭代，遍历并打印每个字符串切片
	for _, item := range records {
		fmt.Println(item)
	}
}
