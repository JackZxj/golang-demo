package main

import (
	"fmt"
	"time"

	"github.com/JackZxj/golang-demo/cache/base"
)

func main() {
	type hh struct {
		base.Base
	}

	u1 := hh{
		Base: base.Base{Name: "hh"},
	}
	u1.BasePrepare()
	fmt.Printf("init\t\t%+v\n", u1.CreatorInfo)

	time.Sleep(4 * time.Second)
	u2 := hh{
		Base: base.Base{Name: "hh"},
	}
	u2.BasePrepare()
	fmt.Printf("cache\t\t%+v\n", u2.CreatorInfo)

	time.Sleep(4 * time.Second)
	u3 := hh{
		Base: base.Base{Name: "hhh"},
	}
	u3.BasePrepare()
	fmt.Printf("init\t\t%+v\n", u3.CreatorInfo)

	time.Sleep(4 * time.Second)
	u4 := hh{
		Base: base.Base{Name: "hh"},
	}
	u4.BasePrepare()
	fmt.Printf("timeout\t\t%+v\n", u4.CreatorInfo)

	time.Sleep(4 * time.Second)
	u5 := hh{
		Base: base.Base{Name: "hh"},
	}
	u5.BasePrepare()
	fmt.Printf("cache\t\t%+v\n", u5.CreatorInfo)

	time.Sleep(4 * time.Second)
	base.ReadCache()
}
