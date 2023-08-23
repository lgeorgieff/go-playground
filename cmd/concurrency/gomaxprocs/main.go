package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("runtime.NumCPU():", runtime.NumCPU())
	// Controlls the number of OS threads that will process work-queues.
	// This is automatcally set to the number of CPU cores.
	// A value of 0 does only return back the current setting.
	fmt.Println("runtime.GOMAXPROCS(runtime.NumCPU()):", runtime.GOMAXPROCS(runtime.NumCPU()))
}
