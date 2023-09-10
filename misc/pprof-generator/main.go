package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func run() error {
	f, err := os.Create("tmp/cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
	t := 1
	for i := 0; i < 3000; i++ {
		t += rand.Intn(100)
		time.Sleep(time.Millisecond)
	}
	fmt.Println(t)
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(fmt.Sprintf("Error:%+v", err))
	}
}
