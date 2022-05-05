//Cache system like redis with go concurrent
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//Complex version in terms of computing
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	return Fibonacci(n-1) + Fibonacci(n-2)
}

//Structs will, store all the results and keys
type Memory struct {
	f Function
	//Create cache
	cache map[int]FunctionResult
	lock  sync.Mutex
}

//Map the keys with the mubers desired to calc Fib
type Function func(key int) (interface{}, error)

type FunctionResult struct {
	value interface{}
	err   error
}

func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FunctionResult),
	}
}

//Method to return content of the cache
func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.Lock()

	result, exists := m.cache[key]
	m.lock.Unlock()

	//If it does not exist calc the value
	if !exists {
		m.lock.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		m.lock.Unlock()

	}
	//Else it does exist and return the value
	return result.value, result.err
}

func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

func main() {
	cache := NewCache(GetFibonacci)
	//Create a new slice
	fibo := []int{42, 40, 41, 42, 38}
	var wg sync.WaitGroup
	maxGoroutines := 2
	channel := make(chan int, maxGoroutines)
	//Iterate the slice to call the function
	for _, n := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			channel <- 1
			start := time.Now()
			value, err := cache.Get(index)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%d, %s, %d\n", index, time.Since(start), value)
			<-channel
		}(n)

	}
	wg.Wait()
}
