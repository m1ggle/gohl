package test

import (
	"fmt"
	"testing"

	"github.com/yahve/gohl/conf"
)

func TestConf(t *testing.T) {
	c := conf.LoadConf()
	fmt.Println(c)
	if c.Database.Host != "localhost" {
		t.Error("Host is not localhost")
	}
	if c.Database.Dbname != "mysql" {
		t.Error("Dbname is not mysql")
	}
	if c.Database.User != "root" {
		t.Error("User is not root")
	}
	if c.Database.Password != "a12bCd3_W45pUq6" {
		t.Error("Password is not a12bCd3_W45pUq6")
	}
	if c.Database.Port != 3306 {
		t.Error("Port is not 3306")
	}

}
