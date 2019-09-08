package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"dbName"`
	User         string `json:"user"`
	Password     string `json:"password"`
	MaxConns     int    `json:"maxConns"`
	MaxIdles     int    `json:"maxIdles"`
}

type MysqlConn struct {
	*sql.DB
}

func NewConn(conn_str string, max_conn, max_idle int) (*MysqlConn, error) {
	db, err := sql.Open("mysql", conn_str)
	if err != nil {
		return nil, err
	} else {
		db.SetMaxOpenConns(max_conn)
		db.SetMaxIdleConns(max_idle)
		err = db.Ping()
		if err != nil {
			return nil, err
		} else {
			return &MysqlConn{db}, nil
		}
	}
}

func (conn *MysqlConn) Stop() error {
	if conn.DB != nil {
		err := conn.DB.Close()
		if err != nil {
			return err
		} else {
			conn.DB = nil
			return nil
		}
	} else {
		return nil
	}
}
