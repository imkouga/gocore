package mysqlx

import "testing"

func TestDefaultDsnCase1(t *testing.T) {
	d := DefaultDsn()
	t.Log(d.String())
}

func TestDsnCase1(t *testing.T) {
	d := Dsn(WithUsername("admin"))
	t.Log(d.String())
}

func TestDsnCase2(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"))
	t.Log(d.String())
}

func TestDsnCase3(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"))
	t.Log(d.String())
}

func TestDsnCase4(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"), WithPort("3307"))
	t.Log(d.String())
}

func TestDsnCase5(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"), WithPort("3307"), WithDatabase("database"))
	t.Log(d.String())
}

func TestDsnCase6(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"), WithPort("3307"), WithDatabase("database"), WithCharset("utf8"))
	t.Log(d.String())
}

func TestDsnCase7(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"), WithPort("3307"), WithDatabase("database"), WithCharset("utf8"), WithParseTime("False"))
	t.Log(d.String())
}

func TestDsnCase8(t *testing.T) {
	d := Dsn(WithUsername("admin"), WithPassword("123456"), WithHost("localhost"), WithPort("3307"), WithDatabase("database"), WithCharset("utf8"), WithParseTime("False"), WithLocal("location"))
	t.Log(d.String())
}
