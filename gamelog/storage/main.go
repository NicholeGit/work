// main
package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/helper"
	"github.com/NicholeGit/work/gamelog/storage/core"
)

//storage.exe -path E:/tmp/new/userdata/ -file E:\tmp\new\upinfo.tmp

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//	file := flag.String("file", "../data/upinfo.tmp", "需要同步的记录文件的路径")
	//	path := flag.String("path", "../data/userdata/", "需要同步的记录文件的路径")
	file := flag.String("file", "", "需要同步的记录文件的路径")
	path := flag.String("path", "", "需要同步的记录文件的路径")
	flag.Parse()
	// start logger
	config := cfg.Get()
	if config["log"] != "" {
		cfg.StartLogger(config["log"])
	}

	//收集错误信息
	defer func() {
		if x := recover(); x != nil {
			helper.ERR("caught panic in main()", x)
		}
	}()
	if len(*file) == 0 || len(*path) == 0 {
		helper.ERR("需要提供运行目录和upinfo.tmp")
		os.Exit(1)
	}
	inco, err := core.NewInsertCore(*file, *path)
	if err != nil {
		helper.ERR("insertCore init is error:", err)
		os.Exit(1)
	}
	inco.Run()
}
