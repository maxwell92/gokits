package mysql

import (
	mylog "gokits/log"
	config "app/backend/common/yce/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

var log = mylog.Log

type MysqlClient struct {
	DB *sql.DB
	/*
		host     string
		user     string
		password string
		database string
		pool     int
	*/
}

var instance *MysqlClient

var once sync.Once

// func NewMysqlClient(host, user, password, database string, pool int) *MysqlClient {
func MysqlInstance() *MysqlClient {
	once.Do(func() {
		instance = new(MysqlClient)
	})
	return instance
}

func (c *MysqlClient) Open() {
	// endpoint := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME + DB_CONNECTION_SUFFIX

	endpoint := config.Instance().GetDbEndpoint()

	db, err := sql.Open(config.DATABASE_DRIVER, endpoint)

	if err != nil {
		log.Fatalf("MysqlClient Open Error: err=%s", err)
		return
	}

	// Set Connection Pool
	db.SetMaxOpenConns(config.Instance().RedisMaxActiveConn)
	db.SetMaxIdleConns(config.Instance().RedisMaxIdleConn)

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
	case <-time.After(time.Millisecond * time.Duration(config.DELAY_MILLISECONDS)):
		err := c.DB.Ping()
		if err != nil {
			log.Fatalf("MysqlClient Ping Error: err=%s", err)
			c.Open()
		}
	}
}
