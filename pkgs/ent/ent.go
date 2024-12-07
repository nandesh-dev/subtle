package ent

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/nandesh-dev/subtle/generated/ent"
	"modernc.org/sqlite"
)

type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys %m", err)
	}
	return conn, nil
}

func Open(filepath string) (*ent.Client, error) {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})

	client, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%s?cache=shared", filepath))
	if err != nil {
		return nil, err
	}

	return client, nil
}
