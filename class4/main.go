package main

import (
	"fmt"
	"sync"
)

var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	b := balance
	balance = b + amount
	lock.Unlock()
}

//Will use RLock and RUnlock to avoid getting blocked
func Balance(lock *sync.RWMutex) int {
	lock.RLock()
	//Reading balance
	b := balance
	lock.RUnlock()
	return b
}

// It could be only 1 doing deposits but it could be N reading

func main() {
	var wg sync.WaitGroup
	var lock sync.RWMutex

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg, &lock)
	}
	wg.Wait()
	fmt.Println(Balance(&lock))
}
