package db

import (
	"fmt"

	"github.com/antsgo/antsgo/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(c conf.ConfigDB) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		c.Username,
		c.Password,
		c.Addr,
		c.DbName,
		c.Config)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("Mysql数据库链接异常:%+v", err)
		return
	}
	conn, err := db.DB()
	if err != nil {
		err = fmt.Errorf("Mysql数据库ping异常:%+v", err)
		return
	}
	err = conn.Ping()
	if err != nil {
		err = fmt.Errorf("Mysql数据库ping异常:%+v", err)
		return
	}
	conn.SetMaxIdleConns(c.MaxIdleConn)
	conn.SetMaxOpenConns(c.MaxOpenConn)
	fmt.Println("Mysql数据库已链接...")
	return
}
