package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/NicholeGit/work/gamelog/cfg"
)

func startPProf() {
	config := cfg.Get()
	cpuprofile := config["cpuprofile"]
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func stopPProf() {
	config := cfg.Get()
	cpuprofile := config["cpuprofile"]
	if cpuprofile != "" {
		pprof.StopCPUProfile()
	}
}

func SignalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	for {
		_ = <-ch
		stopPProf()
		os.Exit(0)
	}
}
