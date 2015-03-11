// main
package main

import (
	"fmt"
	"log"

	"github.com/NicholeGit/work/gamelog/cfg"
)

func main() {
	fmt.Println("Hello World!")
	//收集错误信息
	defer func() {
		if x := recover(); x != nil {
			log.Println("caught panic in main()", x)
		}
	}()

	// start logger
	config := cfg.Get()
	if config["log"] != "" {
		cfg.StartLogger(config["log"])
	}
	filePath := _loadTargetFile("data/upinfo.tmp")
	log.Println(len(filePath), cap(filePath))
	log.Println(filePath)

	usedataPath := "./"
	if config["usedataPath"] != "" {
		usedataPath += config["usedataPath"]
	}
	fmt.Println(usedataPath)

	for _, value := range filePath {
		_loadStorageFile(usedataPath + value)
	}

}
