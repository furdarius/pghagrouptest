package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx"
	stdpgx "github.com/jackc/pgx/stdlib"
)

func main() {
	config := pgx.ConnConfig{
		Host:     "127.0.0.1:5433,127.0.0.1:5432",
		User:     "haha_user",
		Password: "secret",
		Database: "haha",
		Dial: (&net.Dialer{
			KeepAlive: 5 * time.Minute,
			Timeout:   100 * time.Millisecond,
		}).Dial,
		TargetSessionAttrs: pgx.ReadWriteTargetSession, // read-write used to select writable (master) db.
		RuntimeParams: map[string]string{
			"standard_conforming_strings": "on",
		},
		PreferSimpleProtocol: true,
	}

	db := stdpgx.OpenDB(config)

	err := db.Ping()
	if err != nil {
		fmt.Printf("failed to ping database: %v\n", err)
		os.Exit(1)
	}

	db.SetMaxIdleConns(32)
	db.SetMaxOpenConns(32)

	for {
		fmt.Println("try to select")

		err = tryToSelect(db)
		if err != nil {
			time.Sleep(1 * time.Second)

			continue
		}

		fmt.Println("successfully selected!")

		fmt.Println("try to insert")

		err = tryToInsert(db)
		if err != nil {
			time.Sleep(1 * time.Second)

			continue
		}

		fmt.Println("successfully inserted!")

		time.Sleep(5 * time.Second)
	}
}

func tryToSelect(db *sql.DB) error {
	rows, err := db.Query("select name, value from t;")
	if err != nil {
		fmt.Printf("failed to query data: %v\n", err)

		return err
	}
	defer rows.Close() // nolint: errcheck

	for rows.Next() {
		var name, value string

		err = rows.Scan(&name, &value)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)

			continue
		}

		fmt.Printf("selected name: %s, value %s\n", name, value)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("rows error: %v\n", err)

		return err
	}

	return nil
}

func tryToInsert(db *sql.DB) error {
	const q = `insert into t (name, value) values ('Anton', '{"test2": 2}'::jsonb)`

	_, err := db.Exec(q)
	if err != nil {
		fmt.Printf("failed to exec: %v\n", err)

		return err
	}

	return nil
}