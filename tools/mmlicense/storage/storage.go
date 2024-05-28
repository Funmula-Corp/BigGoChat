package storage

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var (
	pgUsername string = "mmuser"
	pgPassword string = "mostest"
	pgHost     string = "localhost"
	pgPort     int    = 5432
	pgDB       string = "mattermost_test"
)

func init() {
	flag.StringVar(&pgUsername, "pg_username", pgUsername, "set the postgres username")
	flag.StringVar(&pgPassword, "pg_password", "mostest", "set the postgres password")
	flag.StringVar(&pgHost, "pg_host", pgHost, "set the postgres host")
	flag.StringVar(&pgDB, "pg_db", pgDB, "set the postgres database")
	flag.IntVar(&pgPort, "pg_port", pgPort, "set the postgres port")
}

func initDBConn() (pgConn *pgx.Conn) {
	var err error
	if pgConn, err = pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", pgUsername, pgPassword, pgHost, pgPort, pgDB)); err != nil {
		log.Fatalln("database connection error:", err)
	}
	return
}

func ActivateLicense(id string) (err error) {
	db := initDBConn()
	defer db.Close(context.Background())

	var tx pgx.Tx
	if tx, err = db.Begin(context.Background()); err != nil {
		db.Close(context.Background())
		log.Println(err)
	} else {
		defer tx.Rollback(context.Background())
		if _, err = tx.Exec(context.Background(), "INSERT INTO systems (name, value) VALUES ('ActiveLicenseId', $1) ON CONFLICT (name) DO UPDATE SET value = $1;", id); err == nil {
			err = tx.Commit(context.Background())
		}
	}
	return
}

func InsertLicense(id string, createdAt int64, buffer []byte) (err error) {
	db := initDBConn()
	defer db.Close(context.Background())

	var tx pgx.Tx
	if tx, err = db.Begin(context.Background()); err != nil {
		db.Close(context.Background())
		log.Fatalln(err)
	} else {
		defer tx.Rollback(context.Background())
		if _, err = tx.Exec(context.Background(), "INSERT INTO licenses (id, createat, bytes) VALUES ($1, $2, $3);", id, createdAt, buffer); err == nil {
			err = tx.Commit(context.Background())
		}
	}
	return
}

func GetActiveLicense() (licenseId string, err error) {
	db := initDBConn()
	defer db.Close(context.Background())

	row := db.QueryRow(context.Background(), "SELECT value::text FROM systems WHERE name = 'ActiveLicenseId';")
	if err = row.Scan(&licenseId); err != nil {
		db.Close(context.Background())
		log.Fatalln(err)
	}
	return
}

func GetLicense(id string) (buffer []byte, err error) {
	db := initDBConn()
	defer db.Close(context.Background())

	var strBuf string
	row := db.QueryRow(context.Background(), "SELECT bytes::text FROM licenses WHERE id = $1;", id)
	if err = row.Scan(&strBuf); err != nil {
		db.Close(context.Background())
		log.Fatalln(err)
	} else {
		buffer = []byte(strBuf)
	}
	return
}
