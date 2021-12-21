package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func main() {
	// TestStruct()
	// TestStu()
	// TestStructMerge()
	// TestNoJsonInline()
	TestEncodeBytes()
}

type Project struct {
	Key1 int  `json:"key1,omitempty"`
	Key2 *int `json:"key2,omitempty"`

	Key3 string  `json:"key3,omitempty"`
	Key4 *string `json:"key4,omitempty"`
}

type Data struct {
	Project `json:",inline"`

	Summary     string `json:"summary"`
	Description string `json:"description"`

	Opt     string  `json:"opt,omitempty"`
	OptP    *string `json:"optP,omitempty"`
	OptInt  int     `json:"optInt,omitempty"`
	OptIntP *int    `json:"optIntP,omitempty"`

	CP0 Project `json:",omitempty"`
	CP1 Project
	CP2 *Project
	CP3 Project  `json:"cp3,omitempty"`
	CP4 *Project `json:"cp4,omitempty"`
}

func TestStruct() {
	num := 2
	str := "hello"

	dataProject := Project{
		Key1: 1,
		Key2: &num,

		Key3: "value",
		Key4: &str,
	}

	data := &Data{
		Project: dataProject,

		Summary:     "Summary",
		Description: "Description",

		Opt:     "opts",
		OptP:    &str,
		OptInt:  3,
		OptIntP: &num,

		CP0: dataProject,
		CP1: dataProject,
		CP2: &dataProject,
		CP3: dataProject,
		CP4: &dataProject,
	}
	d, _ := json.Marshal(data)
	fmt.Println(string(d)) // {"key1":1,"key2":2,"key3":"value","key4":"hello","summary":"Summary","description":"Description","opt":"opts","optP":"hello","optInt":3,"optIntP":2,"CP0":{"key1":1,"key2":2,"key3":"value","key4":"hello"},"CP1":{"key1":1,"key2":2,"key3":"value","key4":"hello"},"CP2":{"key1":1,"key2":2,"key3":"value","key4":"hello"},"cp3":{"key1":1,"key2":2,"key3":"value","key4":"hello"},"cp4":{"key1":1,"key2":2,"key3":"value","key4":"hello"}}
	fmt.Println("--------------")

	data = &Data{
		Summary:     "Summary",
		Description: "Description",
	}
	d, _ = json.Marshal(data)
	fmt.Println(string(d)) // {"summary":"Summary","description":"Description","CP0":{},"CP1":{},"CP2":null,"cp3":{}}
	fmt.Println("--------------")

	str0 := `{
		"key1":1,
		"key2":2,
		"key3":"value",
		"key4":"hello",
		"summary":"Summary",
		"description":"Description",
		"opt":"opts",
		"optP":"hello",
		"optInt":3,
		"optIntP":2,
		"CP0":{"key1":1,"key2":2,"key3":"value","key4":"hello"},
		"CP1":{"key1":1,"key2":2,"key3":"value","key4":"hello"},
		"CP2":{"key1":1,"key2":2,"key3":"value","key4":"hello"},
		"cp3":{"key1":1,"key2":2,"key3":"value","key4":"hello"},
		"cp4":{"key1":1,"key2":2,"key3":"value","key4":"hello"}
		}`
	str1 := `{"summary":"Summary","description":"Description"}`

	var res0 Data
	var res1 Data
	json.Unmarshal([]byte(str0), &res0)
	json.Unmarshal([]byte(str1), &res1)
	fmt.Println(res0) // {{1 0xc0000ba510 value 0xc000096450} Summary Description opts 0xc000096460 3 0xc0000ba548 {0 <nil>  <nil>} {0 <nil>  <nil>} <nil> {1 0xc0000ba550 value 0xc000096480} 0xc000098450}
	fmt.Println("--------------")
	fmt.Println(res1) // {{0 <nil>  <nil>} Summary Description  <nil> 0 <nil> {0 <nil>  <nil>} {0 <nil>  <nil>} <nil> {0 <nil>  <nil>} <nil>}
	fmt.Println("--------------")
}

type Stu struct {
	Name string `json:"name"`
	Age  string
}

func TestStu() {
	s := Stu{
		Name: "sss",
	}
	ss, e := json.Marshal(s)
	fmt.Println(string(ss), e)
	var sss Stu
	e = json.Unmarshal(ss, &sss)
	fmt.Println(sss, e)
}

type A struct {
	ID     int
	Name   string
	Gender string
	Age    int
	Date   time.Time
}

//binding type interface 要修改的结构体
//value type interace 有数据的结构体
func structAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		// fmt.Println(name,
		// 	vVal.Field(i).IsValid(),
		// 	reflect.ValueOf(vVal.Field(i).Interface()),
		// 	reflect.Zero(vVal.Field(i).Type()),
		// 	vVal.Field(i).Type())

		//field.Interface() 当前持有的值
		//reflect.Zero 根据类型获取对应的 零值
		//这个必须调用 Interface 方法 否则为 reflect.Value 构造体的对比 而不是两个值的对比
		//这个地方不要用等号去对比 因为golang 切片类型是不支持 对比的
		if reflect.DeepEqual(vVal.Field(i).Interface(), reflect.Zero(vVal.Field(i).Type()).Interface()) {
			fmt.Println(name, "000000000000000000000000000zero")
			continue
		}
		if bVal.FieldByName(name).IsValid() {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}

func TestStructMerge() {
	as := A{ID: 0, Name: "sss", Age: 12, Date: time.Now()}
	bs := A{Name: "wfy", Gender: "man"}
	fmt.Println("### before:", as)
	structAssign(&as, &bs)
	fmt.Println("### after:", as)
}

type HH struct {
	H1 string `json:"h1"`
	H2 string `json:"h2"`
}

type HHH struct {
	HH
	H3 string `json:"h3"`
}

func TestNoJsonInline() {
	h := HHH{HH{"1", "2"}, "3"}
	v, _ := json.Marshal(h)
	fmt.Println(string(v)) // {"h1":"1","h2":"2","h3":"3"}
}

func TestEncodeBytes() {
	type Bytes struct {
		Value []byte `json:"value"` // []byte默认会被base64编码
	}
	b := Bytes{Value: []byte{'a', 'b', 'c'}}
	v, _ := json.Marshal(b)
	fmt.Println(string(v)) // {"value":"YWJj"}
	json.Unmarshal(v, &b)
	fmt.Println(b) // {[97 98 99]} // ascii a,b,c
}
