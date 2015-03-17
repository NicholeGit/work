// main
package main

import (
	"flag"
	"log"
	"os"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/core"
)

func main() {
	upinfo := flag.String("path", "data/upinfo.tmp", "需要同步的记录文件的路径")
	flag.Parse()
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

	cfg.NOTICE("upinfo name is ", *upinfo)
	inco, err := core.NewInsertCore(*upinfo)
	if err != nil {
		cfg.ERR("insertCore init is error")
		os.Exit(1)
	}
	inco.Run()
	inco.Print()
	log.Println("run end")
}
