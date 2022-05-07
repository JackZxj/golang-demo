package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	f6()
	time.Sleep(2 * time.Second)
	f7()
}

func f1(ctx context.Context) {
	i := 1
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("hello +%ds\n", i)
			i++
		case <-ctx.Done():
			fmt.Println("f1 done!")
			return
		}
	}

}

// Background/TODO 是不可取消、无法截止、无传值的 context
// Background 应当用于 main
// TODO 可以用于中途未知的 context
func f2() {
	ctx := context.Background()
	go f1(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("f2 done")
}

// WithCancel 可取消的 context
// 父线程调用 cancel 时，子线程执行 done
func f3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go f1(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("f3 done")
}

// 截止同理
func f4() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Microsecond))
	defer cancel()
	go f1(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("f4 done")
}

// 超时同理，等价于 WithDeadline(parent, time.Now().Add(timeout))
func f5() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()
	go f1(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("f5 done")
}

// WithValue 可以用于传值
// 任何派生于此 context 的子线程都可以取得该值
// context 传值的 key 需要为可比较类型，且推荐不为内置的数据类型
func f6() {
	type mykey string
	ctx := context.WithValue(context.Background(), mykey("k"), "vvvv")
	var f func(ctx context.Context, k mykey)
	f = func(ctx context.Context, k mykey) {
		if v := ctx.Value(k); v != nil {
			fmt.Printf("got %q: %v\n", k, v)
			return
		}
		fmt.Printf("can not get %q\n", k)
		nctx := context.WithValue(ctx, k, "nctx")
		go f(nctx, "k")
	}
	go f(ctx, mykey("k"))
	go f(ctx, mykey("kk"))
	time.Sleep(2 * time.Second)
	fmt.Println("f6 done")
}

func f7() {
	ch := make(chan bool, 1)
	fn := func(ctx context.Context, num int) {
		// for {
		select {
		case <-ch:
			fmt.Println(num, "got chan")
		case <-ctx.Done():
			fmt.Println(num, "done")
			// return
		}
		// }
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go fn(ctx, 1)
	go fn(ctx, 2)
	go fn(ctx, 3)

	time.Sleep(200 * time.Millisecond)
	ch <- false

	time.Sleep(2 * time.Second)
}
