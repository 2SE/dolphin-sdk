package sql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Config struct {
	Addr         string // for trace
	DSN          string // write data source name.
	Active       int    // pool
	Idle         int    // pool
	ConnTimeout  time.Duration
	IdleTimeout  time.Duration // connect max life time.
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
	TranTimeout  time.Duration // transaction sql timeout

}

//只支持单个db,读写分离连接在业务层配置处理，集群库配置数据库只读
//如需功能型拓展，可以做
type DB struct {
	*sqlx.DB
	Cnf *Config
}

func NewMySQL(c *Config) (*DB, error) {
	db := new(DB)
	dbi, err := connect(c, c.DSN)
	if err != nil {
		return nil, err
	}
	db.DB = dbi
	db.Cnf = c
	return db, nil
}

func connect(c *Config, dataSourceName string) (*sqlx.DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), c.ConnTimeout)
	db, err := sqlx.ConnectContext(ctx, "mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(c.IdleTimeout)
	return db, nil
}
