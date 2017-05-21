package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	mylog "github.com/maxwell92/gokits/log"
	"sync"
	"time"
)

var log = mylog.Log

const (
	DELAY_MILLISECONDS = 5000
	MAX_ACTIVECONN     = 10
	MAX_IDLECONN       = 1
	CONNECTION_SUFFIX =  "?parseTime=true"
	DRIVER             = "mysql"
)

type MysqlClient struct {
	DB       *sql.DB
	host     string
	user     string
	password string
	database string
	pool     int
}

var instance *MysqlClient
var once sync.Once

func NewMysqlClient(host, user, password, database string, pool int) *MysqlClient {
	once.Do(func() {
		instance = &MysqlClient{
			host:     host,
			user:     user,
			password: password,
			database: database,
			pool:     pool,
		}
	)}
	return instance
}

func (c *MysqlClient) Open() {
	endpoint := c.user+ ":" + c.password + "@tcp(" + c.host + ")/" + c.database + CONNECTION_SUFFIX
	db, err := sql.Open(DRIVER, endpoint)
	if err != nil {
		log.Fatalf("MysqlClient Open Error: err=%s", err)
		return
	}

	// Set Connection Pool
	db.SetMaxOpenConns(MAX_ACTIVECONN)
	db.SetMaxIdleConns(MAX_IDLECONN)

	c.DB = db

}

func (c *MysqlClient) Close() {
	c.DB.Close()
}

func (c *MysqlClient) Conn() *sql.DB {
	return c.DB
}

// Ping the connection, keep connection alive
func (c *MysqlClient) Ping() {
	select {
	case <-time.After(time.Millisecond * time.Duration(DELAY_MILLISECONDS)):
		err := c.DB.Ping()
		if err != nil {
			log.Fatalf("MysqlClient Ping Error: err=%s", err)
			c.Open()
		}
	}
}
