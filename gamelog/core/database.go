package core

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/NicholeGit/work/gamelog/cfg"
	_ "github.com/go-sql-driver/mysql"
)

var (
	deleteUser = "delete from %s where account=?"
	insertUser = "insert into %s(account, objectId,count, type, created) values(?,?,?,?,?)"
)

type DataBase struct {
	mysql *sql.DB

	deleteUser *sql.Stmt //删除user
	insertUser *sql.Stmt //插入物品
}

func addTableName(str string) string {
	config := cfg.Get()
	table := config["table"]
	return fmt.Sprintf(str, table)
}

func Open(driverName, dataSourceName string) (*DataBase, error) {
	base := new(DataBase)
	var err error
	base.mysql, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	base.deleteUser, err = base.mysql.Prepare(addTableName(deleteUser))
	if err != nil {
		return nil, err
	}

	base.insertUser, err = base.mysql.Prepare(addTableName(insertUser))
	if err != nil {
		return nil, err
	}

	return base, nil
}

func (this *DataBase) InsertUser(account string, objectId int, count int, ty int) (sql.Result, error) {
	return this.insertUser.Exec(account, objectId, count, ty, time.Now())

}

func (this *DataBase) DeleteUser(account string) (sql.Result, error) {
	return this.deleteUser.Exec(account)
}
