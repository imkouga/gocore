package mysqlx

import (
	"fmt"

	"github.com/imkouga/gocore/loger"
)

type dsn struct {
	username  string
	password  string
	host      string
	port      string
	database  string
	charset   string
	parseTime string
	loc       string
}

type dsnOption interface {
	apply(*dsn)
}

type funcDsnOption struct {
	f func(*dsn)
}

func (f *funcDsnOption) apply(d *dsn) {
	f.f(d)
}

func (d *dsn) String() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", d.username, d.password, d.host, d.port, d.database, d.charset, d.parseTime, d.loc)
	loger.Trace("mysql dsn: ", s)
	return s
}

func DefaultDsn() *dsn {
	return &dsn{
		username:  "root",
		password:  "",
		host:      "127.0.0.1",
		port:      "3306",
		database:  "default",
		charset:   "utf8mb4",
		parseTime: "True",
		loc:       "Local",
	}
}

func Dsn(opts ...dsnOption) *dsn {

	d := DefaultDsn()

	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

func newFuncDsnOption(f func(*dsn)) *funcDsnOption {
	return &funcDsnOption{
		f: f,
	}
}

func WithUsername(usename string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.username = usename
	})
}

func WithPassword(password string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.password = password
	})
}

func WithHost(host string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.host = host
	})
}

func WithPort(port string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.port = port
	})
}

func WithDatabase(database string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.database = database
	})
}

func WithCharset(charset string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.charset = charset
	})
}

func WithParseTime(parseTime string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.parseTime = parseTime
	})
}

func WithLocal(loc string) dsnOption {
	return newFuncDsnOption(func(d *dsn) {
		d.loc = loc
	})
}
