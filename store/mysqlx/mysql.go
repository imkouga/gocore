package mysqlx

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySQL struct {
	db *gorm.DB
}

var defaultMySQL *MySQL

func (ms *MySQL) OK() error {
	return ms.db.Error
}

func (ms *MySQL) Error() string {
	if err := ms.OK(); nil == err {
		return ""
	}
	return ms.db.Error.Error()
}

func (ms *MySQL) Found() bool {
	if ms.db.RowsAffected > 0 {
		return true
	}
	return false
}

func (ms *MySQL) Debug() *MySQL {
	db := ms.db.Debug()
	return NewMySQL(db)
}

func (ms *MySQL) Model(value interface{}) *MySQL {
	db := ms.db.Model(value)
	return NewMySQL(db)
}

func (ms *MySQL) Table(name string) *MySQL {
	db := ms.db.Table(name)
	return NewMySQL(db)
}

func (ms *MySQL) Create(value interface{}) *MySQL {
	db := ms.db.Create(value)
	return NewMySQL(db)
}

func (ms *MySQL) Delete(value interface{}) *MySQL {
	db := ms.db.Delete(value)
	return NewMySQL(db)
}

func (ms *MySQL) FindOne(dest interface{}) *MySQL {
	db := ms.db.First(dest)
	return NewMySQL(db)
}

func (ms *MySQL) FindAll(dest interface{}) *MySQL {
	db := ms.db.Find(dest)
	return NewMySQL(db)
}

func (ms *MySQL) Save(value interface{}) *MySQL {
	db := ms.db.Save(value)
	return NewMySQL(db)
}

func (ms *MySQL) Update(column string, value interface{}) *MySQL {
	db := ms.db.Update(column, value)
	return NewMySQL(db)
}

func (ms *MySQL) Updates(values interface{}) *MySQL {
	db := ms.db.Updates(values)
	return NewMySQL(db)
}

func (ms *MySQL) Where(query interface{}, args ...interface{}) *MySQL {
	db := ms.db.Where(query, args...)
	return NewMySQL(db)
}

func (ms *MySQL) Or(query interface{}, args ...interface{}) *MySQL {
	db := ms.db.Or(query, args)
	return NewMySQL(db)
}

func (ms *MySQL) Order(value interface{}) *MySQL {
	db := ms.db.Order(value)
	return NewMySQL(db)
}

func (ms *MySQL) Group(name string) *MySQL {
	db := ms.db.Group(name)
	return NewMySQL(db)
}

func (ms *MySQL) Begin() *MySQL {
	db := ms.db.Begin()
	return NewMySQL(db)
}

func (ms *MySQL) Commit() *MySQL {
	db := ms.db.Commit()
	return NewMySQL(db)
}

func (ms *MySQL) Rollback() *MySQL {
	db := ms.db.Rollback()
	return NewMySQL(db)
}

func (ms *MySQL) Ping() error {

	db, err := ms.db.DB()
	if nil != err {
		return err
	}
	return db.Ping()
}

func (ms *MySQL) setMaxIdleConns(n int) error {

	db, err := ms.db.DB()
	if nil != err {
		return err
	}

	db.SetMaxIdleConns(n)
	return nil
}

func (ms *MySQL) setMaxOpenConns(n int) error {

	db, err := ms.db.DB()
	if nil != err {
		return err
	}

	db.SetMaxOpenConns(n)
	return nil
}

func (ms *MySQL) setConnMaxLifetime(d time.Duration) error {

	db, err := ms.db.DB()
	if nil != err {
		return err
	}

	db.SetConnMaxLifetime(d * time.Second)
	return nil
}

func (ms *MySQL) setConnMaxIdleTime(d time.Duration) error {

	db, err := ms.db.DB()
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
	return &MySQL{db: db}
}

func CloneMySQL(ms *MySQL) *MySQL {

	tx := &gorm.DB{Config: ms.db.Config}
	tx.Statement = &gorm.Statement{
		DB:       tx,
		ConnPool: ms.db.Statement.ConnPool,
		Context:  ms.db.Statement.Context,
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
