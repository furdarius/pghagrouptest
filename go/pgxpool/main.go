package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx"
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
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: config})
	if err != nil {
		fmt.Printf("failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	for {
		fmt.Println("try to Acquire")
		conn, err := pool.Acquire()
		if err != nil {
			fmt.Printf("failed to acquire connection: %v\n", err)

			time.Sleep(1 * time.Second)

			continue
		}

		err = tryToSelect(conn)
		if err != nil {
			fmt.Println("close conn")
			conn.Close()
			pool.Release(conn)

			time.Sleep(1 * time.Second)

			continue
		}

		fmt.Println("successfully selected!")

		err = tryToInsert(conn)
		if err != nil {
			fmt.Println("close conn")
			conn.Close()
			pool.Release(conn)

			time.Sleep(1 * time.Second)

			continue
		}

		fmt.Println("successfully inserted!")

		pool.Release(conn)

		time.Sleep(5 * time.Second)
	}
}

func tryToSelect(conn *pgx.Conn) error {
	rows, err := conn.Query("select name, value from t;")
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
		return err
	}

	return nil
}

func tryToInsert(conn *pgx.Conn) error {
	const q = `insert into t (name, value) values ('Anton', '{"test2": 2}'::jsonb)`

	_, err := conn.Exec(q)
	if err != nil {
		fmt.Printf("failed to exec: %v\n", err)

		return err
	}

	return nil
}