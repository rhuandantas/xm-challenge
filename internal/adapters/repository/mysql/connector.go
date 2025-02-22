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

type Connector struct {
	engine *xorm.Engine
}

func (m *Connector) GetORM() *xorm.Engine {
	return m.engine
}

func (m *Connector) Close() {
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

	return &Connector{
		engine: engine,
	}
}

func (m *Connector) SyncTables(beans ...interface{}) error {
	return m.engine.Sync(beans...)
}
