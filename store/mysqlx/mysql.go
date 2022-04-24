package mysqlx

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySQL struct {
	*gorm.DB
}

var defaultMySQL *MySQL

func (ms *MySQL) OK() error {
	return ms.DB.Error
}

func (ms *MySQL) Error() string {
	if err := ms.OK(); nil == err {
		return ""
	}
	return ms.DB.Error.Error()
}

func (ms *MySQL) Found() bool {
	if ms.DB.RowsAffected > 0 {
		return true
	}
	return false
}

func (ms *MySQL) setMaxIdleConns(n int) error {

	db, err := ms.DB.DB()
	if nil != err {
		return err
	}

	db.SetMaxIdleConns(n)
	return nil
}

func (ms *MySQL) setMaxOpenConns(n int) error {

	db, err := ms.DB.DB()
	if nil != err {
		return err
	}

	db.SetMaxOpenConns(n)
	return nil
}

func (ms *MySQL) setConnMaxLifetime(d time.Duration) error {

	db, err := ms.DB.DB()
	if nil != err {
		return err
	}

	db.SetConnMaxLifetime(d * time.Second)
	return nil
}

func (ms *MySQL) setConnMaxIdleTime(d time.Duration) error {

	db, err := ms.DB.DB()
	if nil != err {
		return err
	}

	db.SetConnMaxIdleTime(d * time.Second)
	return nil
}

func checkMySQL(ms *MySQL) error {
	if nil == ms {
		return errors.New("mysql driver not init!!")
	}
	return nil
}

func NewMySQL(db *gorm.DB) *MySQL {
	return &MySQL{db}
}

func CloneMySQL(ms *MySQL) *MySQL {

	tx := &gorm.DB{Config: ms.DB.Config}
	tx.Statement = &gorm.Statement{
		DB:       tx,
		ConnPool: ms.DB.Statement.ConnPool,
		Context:  ms.DB.Statement.Context,
		Clauses:  map[string]clause.Clause{},
		Vars:     make([]interface{}, 0, 8),
	}
	return NewMySQL(tx)
}

func LoadMySQL(dsn *dsn) error {

	db, err := gorm.Open(mysql.Open(dsn.String()), &gorm.Config{})
	if nil != err {
		return err
	}

	defaultMySQL = &MySQL{db}

	return nil
}

func GetMySQL() (*MySQL, error) {
	if err := checkMySQL(defaultMySQL); nil != err {
		return nil, err
	}
	return CloneMySQL(defaultMySQL), nil
}

func SetMaxIdleConns(ms *MySQL, n int) error {
	if err := checkMySQL(ms); nil != err {
		return err
	}
	return ms.setMaxIdleConns(n)
}

func SetMaxOpenConns(ms *MySQL, n int) error {
	if err := checkMySQL(ms); nil != err {
		return err
	}
	return ms.setMaxOpenConns(n)
}

func SetConnMaxLifetime(ms *MySQL, d time.Duration) error {
	if err := checkMySQL(ms); nil != err {
		return err
	}
	return ms.setConnMaxLifetime(d)
}

func SetConnMaxIdleTime(ms *MySQL, d time.Duration) error {
	if err := checkMySQL(ms); nil != err {
		return err
	}
	return ms.setConnMaxIdleTime(d)
}
