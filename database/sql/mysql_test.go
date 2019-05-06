package sql

import (
	"testing"
	"time"
)

func TestNewMySQL(t *testing.T) {

	db, err := NewMySQL(&Config{
		DSN:          "root:111111@(127.0.0.1:3306)/wallet",
		Active:       10, // pool
		Idle:         5,  // pool
		IdleTimeout:  time.Minute,
		QueryTimeout: time.Minute,
		ExecTimeout:  time.Minute,
		TranTimeout:  time.Minute,
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}
