package core

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/util"
)

// state.
type RunState struct {

	// 总共需要处理的数据
	Total int

	//读取文件统计
	LoadSucceed uint64 //完成读取数量
	LoadFailure uint64 //读取失败数量

	//数据库统计
	DeleteSucceed uint64 //删除成功数量
	DeleteFailure uint64 //删除失败数量

	InsertSucceed uint64 //插入成功数量
	InsertFailure uint64 //插入失败数量

	AllSucceed uint64 //最后成功完成的账号数
	AllFailure uint64 //最后失败完成的账号数

}

type InsertCore struct {
	mutexUser sync.Mutex
	userList  []util.User
	filePath  []string
	wg        sync.WaitGroup

	DB *DataBase

	//读取文件统计
	loadSucceed uint64 //完成读取数量
	loadFailure uint64 //读取失败数量

	//数据库统计
	deleteSucceed uint64 //删除成功数量
	deleteFailure uint64 //删除失败数量

	insertSucceed uint64 //插入成功数量
	insertFailure uint64 //插入失败数量

	allSucceed uint64 //最后成功完成的账号数
	allFailure uint64 //最后失败完成的账号数

	//	readChan chan *CallInfo
}

func NewInsertCore() (*InsertCore, error) {
	s := new(InsertCore)
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (this *InsertCore) init() error {
	this.userList = make([]util.User, 0, 1024)
	this.filePath = _loadTargetFile("data/upinfo.tmp")
	//log.Println(filePath)
	//test:gamlaxy@tcp(10.100.12.95:3306)/golang?charset=utf8
	config := cfg.Get()
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config["account"], config["password"], config["ip"], config["db"])
	cfg.DEBUG("dataSourceName:\t", str)
	var err error
	this.DB, err = Open("mysql", str)
	return err
}

func (this *InsertCore) addUser(user *util.User) error {
	//fmt.Println(user)
	defer this.mutexUser.Unlock()
	if len(user.Account) == 0 {
		return errors.New("Account is empty")
	}
	this.mutexUser.Lock()
	this.userList = append(this.userList, *user)
	return nil
}

// 准备userList
func (this *InsertCore) readfile() {
	config := cfg.Get()
	usedataPath := "./"
	if config["usedataPath"] != "" {
		usedataPath += config["usedataPath"]
	}
	//	fmt.Println(usedataPath)
	for _, value := range this.filePath {
		//		fmt.Println(value)
		this.wg.Add(1)
		go func(path string) {
			defer this.wg.Done()
			if err := this.addUser(_loadStorageFile(usedataPath + path)); err != nil {
				atomic.AddUint64(&this.deleteFailure, 1)
			} else {
				atomic.AddUint64(&this.loadSucceed, 1)
			}
			return
		}(value)
	}
}

func (this *InsertCore) install() {
	for _, value := range this.userList {
		//		fmt.Println(value)
		this.wg.Add(1)
		go func(user util.User) {
			defer this.wg.Done()
			acc := user.Account
			if _, err := this.DB.DeleteUser(acc); err != nil {
				atomic.AddUint64(&this.deleteFailure, 1)
				cfg.ERR("delete %s is err(%v) ", acc, err)
				return
			} else {
				atomic.AddUint64(&this.deleteSucceed, 1)
			}
			var isSucceed = true
			for _, comObject := range user.ComStorage {
				if _, err := this.DB.InsertUser(acc, comObject.ID, comObject.Count, util.COMMON); err != nil {
					atomic.AddUint64(&this.insertFailure, 1)
					isSucceed = false
					cfg.WARN("insert %s %d(%d)is by com err(%v) ", acc, comObject.ID, comObject.Count, err)
				} else {
					atomic.AddUint64(&this.insertSucceed, 1)
				}
			}
			for _, vipObject := range user.VipStorage {
				if _, err := this.DB.InsertUser(acc, vipObject.ID, vipObject.Count, util.VIP); err != nil {
					atomic.AddUint64(&this.insertFailure, 1)
					isSucceed = false
					cfg.WARN("insert %s %d(%d) by vip is err(%v) ", acc, vipObject.ID, vipObject.Count, err)
				} else {
					atomic.AddUint64(&this.insertSucceed, 1)
				}
			}
			if isSucceed {
				atomic.AddUint64(&this.allSucceed, 1)
			} else {
				atomic.AddUint64(&this.allFailure, 1)
				cfg.ERR("insert %s is err", acc)
			}
			return
		}(value)
	}
}

func (this *InsertCore) Run() {
	cfg.NOTICE(this.GetRunState())
	this.readfile()
	this.wg.Wait()
	cfg.NOTICE(this.GetRunState())
	cfg.NOTICE("读取文件完成")
	this.install()
	this.wg.Wait()
	cfg.NOTICE("插入完成")
	cfg.NOTICE(this.GetRunState())
}

// Get buffer pool state.
func (this *InsertCore) GetRunState() RunState {
	return RunState{
		Total:         len(this.filePath),
		LoadSucceed:   atomic.LoadUint64(&this.loadSucceed),
		LoadFailure:   atomic.LoadUint64(&this.loadFailure),
		DeleteSucceed: atomic.LoadUint64(&this.deleteSucceed),
		DeleteFailure: atomic.LoadUint64(&this.deleteFailure),
		InsertSucceed: atomic.LoadUint64(&this.insertSucceed),
		InsertFailure: atomic.LoadUint64(&this.insertFailure),
		AllSucceed:    atomic.LoadUint64(&this.allSucceed),
		AllFailure:    atomic.LoadUint64(&this.allFailure),
	}
}

func (this *InsertCore) Print() {
	for _, value := range this.userList {
		fmt.Println(value)
	}
}
