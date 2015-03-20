package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NicholeGit/work/gamelog/helper"
)

// 读取upinfo.tmp，得到这次同步需要有那些文件需要处理
func _loadTargetFile(path string) (ret *UserFileSet, err error) {
	f, err := os.Open(path)
	if err != nil {
		helper.ERR("path cannot find")
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	ret = NewUserFileSet()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 是否为".storage.o"结尾
		if strings.HasSuffix(line, ".storage.o") {
			ret.Add(line)
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
func _loadStorageFile(path string) (user *User, err error) {
	u := new(User)
	f, err := os.Open(path)
	if err != nil {
		helper.WARN(fmt.Sprintf("%s can't open err(%v)", path, err))
		return nil, err
	}
	defer f.Close()
	u.ComStorage = make([]Item, 0, 120)
	u.VipStorage = make([]Item, 0, 240)

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
