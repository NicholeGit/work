package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NicholeGit/work/gamelog/helper"
	"github.com/NicholeGit/work/gamelog/util"
)

func _loadTargetFile(path string) (ret []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		helper.ERR("path cannot find")
		return nil, err
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

	return ret, nil
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

//读取具体文件
func _loadStorageFile(path string) (user *util.User, err error) {
	u := new(util.User)
	f, err := os.Open(path)
	if err != nil {
		helper.WARN(fmt.Sprintf("%s can't open err(%v)", path, err))
		return nil, err
	}
	defer f.Close()
	u.ComStorage = make([]util.Item, 0, 120)
	u.VipStorage = make([]util.Item, 0, 240)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 是否为".storage.o"结尾
		ret := strings.FieldsFunc(line, splitInit)
		if len(ret) > 1 {
			switch ret[0] {
			case "DEPOSIT_ITEMS":
				u.ComStorage = parseItem(ret)
			case "POS_ITEMS":
				u.VipStorage = parseItem(ret)
			case "ACCOUNT":
				u.Account = ret[1]

			}
		}
	}
	return u, nil
}

func parseItem(str []string) (ret []util.Item) {
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
				item := util.Item{id, conut}
				ret = append(ret, item)
			}
		}
	}
	return
}
