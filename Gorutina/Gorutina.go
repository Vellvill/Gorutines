package main

import (
	"fmt"
	"sync"
	"time"
)

var start time.Time

func init() {
	start = time.Now()
}

func main() {
	in1 := make(chan int)
	in2 := make(chan int)
	out := make(chan int)
	n := 3
	merge2Channels(fn, in1, in2, out, n)
	for i := 0; i < 3; i++ {
		in1 <- i
		in2 <- 10 + i
	}
	fmt.Println("not blocked if printed first")
	time.Sleep(time.Second * 3)

	for i := 0; i < n; i++ {
		fmt.Println(<-out)
	}
	fmt.Println(time.Since(start))
}
func fn(x int) int {
	time.Sleep(time.Millisecond * 100)
	return x + 1
}
func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	array1 := []int{}
	array2 := []int{}
	mu := new(sync.Mutex)
	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(n * 2)
		go func(mu *sync.Mutex, wg *sync.WaitGroup) {
			for {
				select {
				case x := <-in1:
					mu.Lock()
					array1 = append(array1, x)
					mu.Unlock()
					wg.Done()
					if len(array2) == n {
						return
					}
				}
			}
		}(mu, wg)
		go func(mu *sync.Mutex, wg *sync.WaitGroup) {
			for {
				select {
				case x := <-in2:
					mu.Lock()
					array2 = append(array2, x)
					mu.Unlock()
					wg.Done()
					if len(array2) == n {
						return
					}
				}
			}
		}(mu, wg)
		wg.Wait()
		go func() {
			wg2 := new(sync.WaitGroup)
			mu2 := new(sync.Mutex)
			wg2.Add(n * 2)
			go func(wg2 *sync.WaitGroup, mu2 *sync.Mutex) {
				go func() {
					for i, v := range array1 {
						array1[i] = fn(v)
						wg2.Done()
					}
				}()
				go func() {
					for i, v := range array2 {
						array2[i] = fn(v)
						wg2.Done()
					}
				}()
			}(wg2, mu2)
			wg2.Wait()
			for i, v := range array1 {
				out <- v + array2[i]
			}
		}()
	}()
}
