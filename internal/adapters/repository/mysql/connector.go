package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/rhuandantas/xm-challenge/config"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type DBConnector interface {
	GetORM() *xorm.Engine
	Close()
	SyncTables(beans ...interface{}) error
}

type MySQLConnector struct {
	engine  *xorm.Engine
	configs *config.Config
}

func (m *MySQLConnector) GetORM() *xorm.Engine {
	return m.engine
}

func (m *MySQLConnector) Close() {
	err := m.engine.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewMySQLConnector(configs *config.Config) DBConnector {
	// TODO put in env vars
	var (
		err error
	)

	dbUrl := configs.Database.Url
	if dbUrl == "" {
		log.Fatal("make sure your db variable are configured properly")
	}

	engine, err := xorm.NewEngine("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(true) // TODO it should come from env
	//engine.Logger().SetLevel(log.DEBUG)
	engine.SetMapper(names.SnakeMapper{})

	return &MySQLConnector{
		engine: engine,
	}
}

func (m *MySQLConnector) SyncTables(beans ...interface{}) error {
	return m.engine.Sync(beans...)
}
