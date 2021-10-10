package group_task

import (
	"fmt"
	"testing"
	"time"
)

func TestGroupTask(t *testing.T) {
	gt := NewGroupTask()

	go func() {
		gt.Do("1", taskA)
		gt.Do("1", taskB)
	}()

	go func() {
		gt.Do("2", taskC)
		gt.Do("2", taskD)
	}()

	time.Sleep(time.Hour)
}

func taskA() {
	time.Sleep(time.Second)
	fmt.Println("A")
}
func taskB() {
	time.Sleep(time.Second)
	fmt.Println("B")
}
func taskC() {
	time.Sleep(time.Second)
	fmt.Println("C")
}
func taskD() {
	time.Sleep(time.Second)
	fmt.Println("D")
}
