package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/yahve/gohl/conf"
	"golang.org/x/term"
)

func promptPassword() string {
	fmt.Print("Enter password: ")
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return ""
	}
	return string(pw)
}

func GetDBConnection(cmd *cobra.Command) (*sql.DB, error) {
	// Viper already has the values, prioritized correctly (Flag > Config > Default)
	cfg := conf.LoadConf()
	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.User
	password := cfg.Database.Password
	dbname := cfg.Database.Dbname
	log.Printf("Connecting to %s:%d as %s", host, port, user)

	if password == "" {
		password = promptPassword()
	}

	if dbname == "" {
		dbname = "mysql"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return db, nil
}
