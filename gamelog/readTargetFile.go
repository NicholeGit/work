package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func _loadTargetFile(path string) (ret []string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}
	ret = make([]string, 0, 1024)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 是否为".storage.o"结尾
		if strings.HasSuffix(line, ".storage.o") {
			ret = append(ret, line)
		}
	}

	return
}

// use 使用次数
// fix 耐久
// 返回总是用次数
func _countUsable(use int, fix int) int {
	// use = 5   fix = 8
	// (8-1 + 1) * (8-1) / 2 + use
	return fix*(fix-1)/2 + use
}

func _loadStorageFile(path string) (comStorage []Item, vipStorage []Item) {
	//fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}
	comStorage = make([]Item, 0, 120)
	vipStorage = make([]Item, 0, 240)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	//	re := regexp.MustCompile(`[\t ]*DEPOSIT_ITEMS [({"]*`)
	re := regexp.MustCompile(`\"(.*?)\"`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 是否为".storage.o"结尾
		if strings.HasPrefix(line, "DEPOSIT_ITEMS") {
			slice := re.FindStringSubmatch(line)
			fmt.Println(slice)

			//			if slice != nil {
			//				ret[slice[1]] = slice[2]
			//				log.Println(slice[1], "=", slice[2])
			//			}
		}
	}

	return
}
