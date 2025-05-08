package main

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
)

func Foo(in string) string {
	return in
}

type A struct{}

func (a A) Foo(in string) string  { return in }
func (a *A) Bar(in string) string { return in }

var Bar = 0

func TestMockXXX(t *testing.T) {
	mockey.PatchConvey("TestMockXXX", t, func() {
		mockey.Mock(Foo).Return("c").Build()        // mock函数
		mockey.Mock(A.Foo).Return("c").Build()      // mock方法
		mockey.Mock((*A).Bar).Return("c").Build()   // mock方法
		mockey.MockValue(&Bar).To(1)                // mock变量
		mockey.Mock(math.Floor).Return(0.1).Build() // mock方法

		convey.So(Foo("a"), convey.ShouldEqual, "c")        // 断言`Foo`成功mock
		convey.So(new(A).Foo("b"), convey.ShouldEqual, "c") // 断言`A.Foo`成功mock
		convey.So(new(A).Bar("d"), convey.ShouldEqual, "c") // 断言`A.Bar`成功mock
		convey.So(Bar, convey.ShouldEqual, 1)               // 断言`Bar`成功mock
		convey.So(math.Floor(1.1), convey.ShouldEqual, 0.1) // 断言
	})
	// `PatchConvey`外自动释放mock
	fmt.Println(Foo("a"))        // a
	fmt.Println(new(A).Foo("b")) // b
	fmt.Println(new(A).Bar("d")) // d
	fmt.Println(Bar)             // 0
	fmt.Println(math.Floor(1.2))
}

// func TestMockXXXX(t *testing.T) {
// 	PatchConvey("TestMockXXX", t, func() {
// 		Mock(Foo).Return("c").Build()   // mock函数
// 		Mock(A.Foo).Return("c").Build() // mock方法
// 		MockValue(&Bar).To(1)           // mock变量

// 		convey.So(Foo("a"), convey.ShouldEqual, "c")        // 断言`Foo`成功mock
// 		convey.So(new(A).Foo("b"), convey.ShouldEqual, "c") // 断言`A.Foo`成功mock
// 		convey.So(Bar, convey.ShouldEqual, 1)               // 断言`Bar`成功mock
// 	})
// 	// `PatchConvey`外自动释放mock
// 	fmt.Println(Foo("a"))        // a
// 	fmt.Println(new(A).Foo("b")) // b
// 	fmt.Println(Bar)             // 0
// }·

func main() {
	m := map[string]string{
		"123": "1",
		"qwe": "2",
		"asd": "3",
		"zxc": "4",
		"vbn": "5",
		"fgh": "6",
		"rty": "7",
		"456": "8",
		"789": "9",
		"uio": "0",
	}
	for i := 0; i < 10; i++ {
		v, err := json.Marshal(m)
		if err != nil {
			fmt.Println(i, "err", err)
			continue
		}
		fmt.Printf("%s\n", v)
		err = json.Unmarshal(v, &m)
		if err != nil {
			fmt.Println(i, "err", err)
			continue
		}
		fmt.Println(m)

		copy := map[string]string{}
		for k, v := range m {
			copy[k] = v
		}
		m = copy
	}

	TestMockXXX(&testing.T{})
}
