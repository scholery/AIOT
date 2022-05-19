package db

import (
	"os"

	"koudai-box/conf"

	"koudai-box/iot/db/common"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" // driver
	// "aiot/db/common"
	// "aiot/db/dbweb"
)

func DBInit() {
	orm.DefaultRowsLimit = common.MAX_LIMIT
	orm.Debug = false
	//os.MkdirAll("db", 0660)
	os.MkdirAll(conf.GetConf().DbPath, os.ModePerm)
	os.MkdirAll(conf.GetConf().TempPath, os.ModePerm)
	os.MkdirAll(conf.GetConf().IotImagePath, os.ModePerm)
	CreateTables()
	RegisterModels()
	Init()
	InitDict()
}
