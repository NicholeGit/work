// main
package main

import (
	"log"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/core"
)

func main() {
	// start logger
	config := cfg.Get()
	if config["log"] != "" {
		cfg.StartLogger(config["log"])
	}
	log.Println("run start")
	//收集错误信息
	defer func() {
		if x := recover(); x != nil {
			log.Println("caught panic in main()", x)
		}
	}()

	inco, err := core.NewInsertCore()
	_ = err
	inco.Run()
	inco.Print()
	log.Println("run end")
}
