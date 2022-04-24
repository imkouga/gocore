package mysqlx

import (
	"testing"
)

type DefaultTable struct {
	ID   int64  `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
}

func (t *DefaultTable) TableName() string {
	return "default"
}

func loadDefault(t *testing.T) *MySQL {
	if err := LoadMySQL(DefaultDsn()); nil != err {
		t.Fatal(err)
	}
	db, err := GetMySQL()
	if nil != err {
		t.Fatal(err)
	}
	return db
}

func TestLoadMySQLCase1(t *testing.T) {
	defaultDsn := DefaultDsn()
	if err := LoadMySQL(defaultDsn); nil != err {
		t.Fatal(err)
	}
}
