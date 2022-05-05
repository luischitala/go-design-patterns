//Cache system like redis with go non concurrent
package main

import (
	"fmt"
	"log"
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
	result, exists := m.cache[key]
	//If it does not exist calc the value
	if !exists {
		result.value, result.err = m.f(key)
		m.cache[key] = result
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
	//Iterate the slice to call the function
	for _, n := range fibo {
		start := time.Now()
		value, err := cache.Get(n)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%d, %s, %d\n", n, time.Since(start), value)
	}
}
