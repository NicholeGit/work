package core

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/helper"
	"github.com/NicholeGit/work/gamelog/misc/timer"
)

const (
	PRINT_INTERVAL  = 60 //1分钟，输入服务器状态的频率
	GOROUTINE_COUNT = 64 //同时跑的goroutine数量
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
	DeleteRows    uint64

	InsertSucceed uint64 //插入成功数量
	InsertFailure uint64 //插入失败数量

	AllSucceed uint64 //最后成功完成的账号数
	AllFailure uint64 //最后失败完成的账号数

}

func (this RunState) Print() string {
	return fmt.Sprintf("总共用户:%v\t完成:%v(%v/%v)\t完成度:%v%%\t总共读取用户:%v(%v/%v)\t总共删除用户:%v(%v/%v)\t总共删除宝物:%v\t总共插入宝物:%v(%v/%v)",
		this.Total,
		this.AllSucceed+this.AllFailure, this.AllSucceed, this.AllFailure,
		(float64(this.AllSucceed)+float64(this.AllFailure))/float64(this.Total)*100,
		this.LoadSucceed+this.LoadFailure, this.LoadSucceed, this.LoadFailure,
		this.DeleteSucceed+this.DeleteFailure, this.DeleteSucceed, this.DeleteFailure,
		this.DeleteRows,
		this.InsertSucceed+this.InsertFailure, this.InsertSucceed, this.InsertFailure)
}

// core
type InsertCore struct {
	targetFile  string
	usedataPath string
	//	mutexUser   sync.Mutex
	//	userList    []util.User
	filePath *UserFileSet

	DB *DataBase

	//channel
	gchan chan bool //控制goroutine数量

	//	gReadChan chan bool //控制goroutine数量

	//读取文件统计
	loadSucceed uint64 //完成读取数量
	loadFailure uint64 //读取失败数量

	//数据库统计
	deleteSucceed uint64 //删除成功数量
	deleteFailure uint64 //删除失败数量
	deleteRows    uint64 //删除的行数

	insertSucceed uint64 //插入成功数量
	insertFailure uint64 //插入失败数量

	allSucceed uint64 //最后成功完成的账号数
	allFailure uint64 //最后失败完成的账号数

}

func NewInsertCore(fileName string, path string) (*InsertCore, error) {
	s := new(InsertCore)
	s.targetFile = fileName
	s.usedataPath = path
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (this *InsertCore) init() error {
	//	this.userList = make([]util.User, 0, 1024)
	var err error
	if this.filePath, err = _loadTargetFile(this.targetFile); err != nil {
		return err
	}
	config := cfg.Get()
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Asia%%2FShanghai",
		config["username"], config["password"], config["ip"], config["db"])
	helper.NOTICE("dataSourceName:\t", str)
	this.DB, err = Open("mysql", str)

	var count = GOROUTINE_COUNT
	if config["goroutine_count"] != "" {
		if id, err := strconv.Atoi(config["goroutine_count"]); err == nil {
			count = id
		}
	}
	this.gchan = make(chan bool, count)
	//	this.gReadChan = make(chan bool, 1024)
	return err
}

func (this *InsertCore) addUser(user *User) error {
	//fmt.Println(user)
	if len(user.Account) == 0 {
		return errors.New("Account is empty")
	}
	//	this.mutexUser.Lock()
	//	defer this.mutexUser.Unlock()
	//	this.userList = append(this.userList, *user)
	return nil
}

func (this *InsertCore) readfile(fileName string) (user *User, err error) {
	path := this.usedataPath + fileName
	user, err = _loadStorageFile(path)
	if err != nil {
		atomic.AddUint64(&this.loadFailure, 1)
		return nil, errors.New(fmt.Sprintf("%s is loadFailure", path))
	} else {
		atomic.AddUint64(&this.loadSucceed, 1)
		if err := this.addUser(user); err != nil {
			return nil, err
		} else {
			return user, nil
		}
	}
}

func (this *InsertCore) install(user *User) {
	acc := user.Account
	if res, err := this.DB.DeleteUser(acc); err != nil {
		atomic.AddUint64(&this.deleteFailure, 1)
		helper.ERR(fmt.Sprintf("delete %s is err(%v) ", acc, err))
		return
	} else {
		atomic.AddUint64(&this.deleteSucceed, 1)
		if affect, err := res.RowsAffected(); err != nil {
			helper.WARN(fmt.Sprintf("delete %s RowsAffected is err(%v)", acc, err))
		} else {
			atomic.AddUint64(&this.deleteRows, uint64(affect))
		}
	}
	var isSucceed = true
	for _, comObject := range user.ComStorage {
		if _, err := this.DB.InsertUser(acc, comObject.ID, comObject.Count, COMMON); err != nil {
			atomic.AddUint64(&this.insertFailure, 1)
			isSucceed = false
			helper.WARN(fmt.Sprintf("insert %s %d(%d)is by com err(%v) ", acc, comObject.ID, comObject.Count, err))
		} else {
			atomic.AddUint64(&this.insertSucceed, 1)

		}
	}
	for _, vipObject := range user.VipStorage {
		if _, err := this.DB.InsertUser(acc, vipObject.ID, vipObject.Count, VIP); err != nil {
			atomic.AddUint64(&this.insertFailure, 1)
			isSucceed = false
			helper.WARN(fmt.Sprintf("insert %s %d(%d) by vip is err(%v) ", acc, vipObject.ID, vipObject.Count, err))
		} else {
			atomic.AddUint64(&this.insertSucceed, 1)
		}
	}
	if isSucceed {
		atomic.AddUint64(&this.allSucceed, 1)
	} else {
		atomic.AddUint64(&this.allFailure, 1)
		helper.ERR(fmt.Sprintf("insert %s is err", acc))
	}
	return

}

func (this *InsertCore) updateDB(file string) {
	if user, err := this.readfile(file); err == nil {
		this.install(user)
	} else {
		helper.WARN(fmt.Sprintf("updateDB %v is err(%v)", user, err))
	}
}

func (this *InsertCore) StatsAgent() {
	queue_timer := make(chan int32, 1)
	queue_timer <- 1
	for {
		select {
		case <-queue_timer:
			helper.INFO(this.GetRunState().Print())
			timer.Add(TinsertCoreStat, time.Now().Unix()+PRINT_INTERVAL, queue_timer)
		}
	}
}

func (this *InsertCore) Run() {
	helper.NOTICE("run start")
	go this.StatsAgent()
	var wg sync.WaitGroup
	for _, value := range this.filePath.List() {
		this.gchan <- true
		wg.Add(1)
		go func(file string) {
			defer func() {
				<-this.gchan
				wg.Done()
			}()
			this.updateDB(file)
		}(value)
	}
	wg.Wait()
	helper.NOTICE("InsertCore end\t", this.GetRunState().Print())
}

// 得到运行状态
func (this *InsertCore) GetRunState() RunState {
	return RunState{
		Total:         this.filePath.Len(),
		LoadSucceed:   atomic.LoadUint64(&this.loadSucceed),
		LoadFailure:   atomic.LoadUint64(&this.loadFailure),
		DeleteSucceed: atomic.LoadUint64(&this.deleteSucceed),
		DeleteFailure: atomic.LoadUint64(&this.deleteFailure),
		DeleteRows:    atomic.LoadUint64(&this.deleteRows),
		InsertSucceed: atomic.LoadUint64(&this.insertSucceed),
		InsertFailure: atomic.LoadUint64(&this.insertFailure),
		AllSucceed:    atomic.LoadUint64(&this.allSucceed),
		AllFailure:    atomic.LoadUint64(&this.allFailure),
	}
}
