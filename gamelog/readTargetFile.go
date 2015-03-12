package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

func splitInit(s rune) bool {
	if s == ' ' || s == '"' {
		return true
	}
	return false
}

func splitUnderline(s rune) bool {
	if s == '_' {
		return true
	}
	return false
}

func _loadStorageFile(path string) (name string, comStorage []Item, vipStorage []Item) {
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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 是否为".storage.o"结尾
		ret := strings.FieldsFunc(line, splitInit)
		if len(ret) > 1 {
			switch ret[0] {
			case "DEPOSIT_ITEMS":
				comStorage = parseItem(ret)
			case "POS_ITEMS":
				vipStorage = parseItem(ret)
			case "ACCOUNT":
				name = ret[1]
			}
		}
	}
	fmt.Println(name, comStorage, vipStorage)
	return
}

func parseItem(str []string) (ret []Item) {
	for _, v := range str {
		if v == "0_0_0" {
			return
		}
		if strings.Count(v, "_") == 2 {
			under := strings.FieldsFunc(v, splitUnderline)
			if len(under) > 2 {
				id, err1 := strconv.Atoi(under[0])
				use, err2 := strconv.Atoi(under[1])
				fix, err3 := strconv.Atoi(under[2])
				if err1 != nil || err2 != nil || err3 != nil {
					return
				}
				conut := _countUsable(use, fix)
				item := Item{id, conut}
				ret = append(ret, item)
			}
		}
	}
	return
}
