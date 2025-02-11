package main

import (
	"fmt"
	"runtime"
	"strings"
)

func IsLockedToThread() bool {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	fmt.Println(stack)

	return strings.Contains(stack, "locked to thread")
}

func main() {
	runtime.LockOSThread()

	if IsLockedToThread() {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
