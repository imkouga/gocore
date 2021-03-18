package mysqlx

import (
	"testing"
	"time"
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

func TestPingCase1(t *testing.T) {

	if err := LoadMySQL(DefaultDsn()); nil != err {
		t.Fatal(err)
	}
	db, err := GetMySQL()
	if nil != err {
		t.Fatal(err)
	}
	if err := db.Ping(); nil != err {
		t.Fatal(err)
	}
}

func TestCreateCase1(t *testing.T) {

	db := loadDefault(t)
	m := &DefaultTable{
		ID:   1,
		Name: "string",
	}
	db.Create(m)
}

// trans commit
func TestCreateCase2(t *testing.T) {

	db := loadDefault(t)
	db = db.Begin()
	m := &DefaultTable{
		Name: "string",
	}
	db = db.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	db = db.Create(m)
	db.Commit()
}

// trans not commit
func TestCreateCase3(t *testing.T) {

	db := loadDefault(t)
	db = db.Begin()
	m := &DefaultTable{
		Name: "string",
	}
	db = db.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	db = db.Create(m)
}

// trans rollback
func TestCreateCase4(t *testing.T) {

	db := loadDefault(t)
	db = db.Begin()
	m := &DefaultTable{
		Name: "string",
	}
	db = db.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	db = db.Create(m)
	db.Rollback()
}

// 2个trans并发  一个 commit 一个不 commit
func TestCreateCase5(t *testing.T) {
	db := loadDefault(t)
	notCommit(db, t)
	commit(db, t)
}

// 2个trans并发  一个 commit 一个延迟commit
func TestCreateCase6(t *testing.T) {
	db := loadDefault(t)
	go commitDelay(db, t)
	commit(db, t)
	time.Sleep(time.Second * 35)
}

func commit(db *MySQL, t *testing.T) {
	b := db.Begin().Debug()
	m := &DefaultTable{
		Name: "string",
	}
	b = b.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	b = b.Create(m)
	b.Commit()
}

func commitDelay(db *MySQL, t *testing.T) {
	b := db.Begin()
	m := &DefaultTable{
		Name: "string",
	}
	b = b.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	b = b.Create(m)
	time.Sleep(time.Second * 30)
	b.Commit()
}

func notCommit(db *MySQL, t *testing.T) {
	b := db.Begin()
	m := &DefaultTable{
		Name: "string",
	}
	b = b.Create(m)
	m = &DefaultTable{
		Name: "34",
	}
	b = b.Create(m)
}

func TestFindCase1(t *testing.T) {

	db := loadDefault(t)
	m := &DefaultTable{}
	result := db.FindOne(m)
	if err := result.OK(); nil != err {
		t.Fatal(err)
	}
	t.Log(m)
}

func TestFindAndSaveCase2(t *testing.T) {

	db := loadDefault(t).Debug()
	m := &DefaultTable{}
	result := db.FindOne(m)
	if err := result.OK(); nil != err {
		t.Fatal(err)
	}
	t.Log(m)
	m.Name = "789456"
	result = db.Save(m)
	if err := result.OK(); nil != err {
		t.Fatal(err)
	}
}

func TestFindAllCase1(t *testing.T) {

	db := loadDefault(t).Debug()
	var m []*DefaultTable
	result := db.FindAll(&m)
	if err := result.OK(); nil != err {
		t.Fatal(err)
	}
	t.Log(m)
}

func TestFindAllWhereCase1(t *testing.T) {
	db := loadDefault(t).Debug()
	var m []*DefaultTable
	if result := db.Where("name = ?", "789456").FindAll(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
	t.Log(len(m))
}

func TestFindAllWhereOrCase1(t *testing.T) {
	db := loadDefault(t).Debug()
	var m []*DefaultTable
	if result := db.Where("name = ?", "789456").Or("id = ?", 2).FindAll(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
	t.Log(len(m))
}

func TestFindAllWhereAndCase1(t *testing.T) {
	db := loadDefault(t).Debug()
	var m []*DefaultTable
	if result := db.Where("name = ?", "789456").Where("id = ?", 2).FindAll(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
	t.Log(len(m))
}

func TestFindAllWhereOrOrCase1(t *testing.T) {
	db := loadDefault(t).Debug()
	var m []*DefaultTable
	if result := db.Where("name = ?", "789456").Where("id = ? OR id = ?", 2, 4).FindAll(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
	t.Log(len(m))
}

func TestDeleteCase1(t *testing.T) {
	db := loadDefault(t).Debug().Begin()
	var m DefaultTable
	if result := db.Delete(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
}

func TestDeleteWhereCase1(t *testing.T) {
	db := loadDefault(t).Debug()
	var m DefaultTable
	if result := db.Where("id = ?", 100).Delete(&m); result.OK() != nil {
		t.Fatal(result.Error())
	}
}
