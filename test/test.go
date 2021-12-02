package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func RunTask(d int) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("%d:%d\n", d, 111)
		time.Sleep(time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("%d:%d\n", d, 222)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("%d:%d\n", d, 333)
	}()
}

func main() {

	// 检测队列是否没有未消耗

	// 监控协程使用情况，过高发预警
	maxG := make(chan int)
	go func() {
		for {
			num := runtime.NumGoroutine()
			if num > 10 {
				maxG <- num
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	// 监听协程数量
	go func() {
		for {
			select {
			case <-maxG:
				panic("协程数量超过限制")
			}
		}
	}()

	i := 1
	for {
		ch := make(chan bool)
		go func(d int) {

			RunTask(d)
			wg.Wait()
			fmt.Println("=====")
			ch <- true

		}(i)

		<-ch

		i++
	}

}

// func main() {

// 	// 检测队列是否没有未消耗
// 	task := make(chan string)
// 	stop := make(chan bool)
// 	go func() {
// 		for {
// 			select {
// 			case task := <-task:
// 				var wg sync.WaitGroup

// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					fmt.Printf("%s:%d\n", task, 111)
// 					time.Sleep(time.Second)
// 				}()

// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					fmt.Printf("%s:%d\n", task, 222)
// 				}()

// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					fmt.Printf("%s:%d\n", task, 333)
// 				}()

// 				wg.Wait()
// 			case <-stop:
// 				fmt.Println("stop")
// 				return
// 			}

// 		}
// 	}()

// 	i := 1
// 	for {
// 		ch := make(chan bool)
// 		go func(d int) {

// 			if i == 10 {
// 				stop <- true
// 			} else {
// 				task <- fmt.Sprintf("tasking %d", i)
// 			}
// 			time.Sleep(time.Second)
// 			ch <- true
// 		}(i)

// 		<-ch
// 		i++
// 	}

// }
