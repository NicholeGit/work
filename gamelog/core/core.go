package core

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/NicholeGit/work/gamelog/cfg"
	"github.com/NicholeGit/work/gamelog/util"
)

type InsertCore struct {
	mutexUser sync.Mutex
	userList  []util.User
	filePath  []string
	wg        sync.WaitGroup

	DB *DataBase

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
	var err error
	this.DB, err = Open("mysql", "sgcuser:gamlaxy@tcp(10.100.0.181:3306)/testStorage?charset=utf8")
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
			this.addUser(_loadStorageFile(usedataPath + path))
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

			return
		}(value)
	}
}

func (this *InsertCore) Run() {
	this.readfile()
	this.wg.Wait()
	log.Println("插入userlist完成")
	this.install()
	this.wg.Wait()
	log.Println("插入完成")
}

func (this *InsertCore) Print() {
	for _, value := range this.userList {
		fmt.Println(value)
	}
}
