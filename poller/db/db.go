package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"poller/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	username = "root"
	password = "pass1234"
	hostname = "mysql:3306"
	dbname   = "pollerdb"
)

func Init() {
	createDb()
	createTable()
}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func openDB() {
	var err error
	db, err = sql.Open("mysql", dsn())
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		panic(err)
	}
}

func CloseDB() {
	db.Close()
}

func createDb() {
	openDB()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)

	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected %d\n", no)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
}

func createTable() {
	// create table watch
	_, err1 := db.Exec(`CREATE TABLE IF NOT EXISTS pollerdb.watch(
		watch_id varchar(100) NOT NULL,
		user_id varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		zipcode varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		alerts json NOT NULL,
		watch_created datetime NOT NULL,
		watch_updated datetime NOT NULL,
		PRIMARY KEY (watch_id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
	// create table alert
	_, err2 := db.Exec(`CREATE TABLE IF NOT EXISTS pollerdb.alert(
		alert_id varchar(100) NOT NULL,
		watch_id varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		field_type ENUM('temp', 'feels_like', 'temp_min', 'temp_max', 'pressure','humidity') COLLATE utf8_unicode_ci NOT NULL,
		operator ENUM('gt', 'gte', 'eq', 'lt', 'lte') COLLATE utf8_unicode_ci NOT NULL,
		value int NOT NULL,
		alert_created datetime NOT NULL,
		alert_updated datetime NOT NULL,
		PRIMARY KEY (alert_id),
		FOREIGN KEY (watch_id) REFERENCES watch(watch_id) 
		)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)

	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
}

func InsertWatch(watch model.WATCH) bool {
	insert, err := db.Prepare(`INSERT INTO pollerdb.watch(watch_id, user_id, zipcode, alerts, watch_created, watch_updated) 
						VALUES (?, ?, ?, ?, ?, ?)`)

	log.Print("insert p")

	if err != nil {
		log.Print(err.Error())
		return false
	}

	log.Print("insert prepare")
	alertsJson, _ := json.Marshal(&watch.Alerts)

	res, err := insert.Exec(watch.ID, watch.UserId, watch.Zipcode, alertsJson, watch.WatchCreated, watch.WatchUpdated)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	log.Print(res.RowsAffected())

	for i := range watch.Alerts {
		uid, _ := uuid.NewRandom()
		watch.Alerts[i].ID = uid.String()
		watch.Alerts[i].WatchId = watch.ID
		watch.Alerts[i].AlertCreated = watch.WatchCreated
		watch.Alerts[i].AlertUpdated = watch.WatchCreated
	}

	for _, a := range watch.Alerts {
		if !insertAlert(a) {
			return false
		}
	}

	return true
}

func insertAlert(alert model.ALERT) bool {
	insert, err := db.Prepare(`INSERT INTO pollerdb.alert(alert_id, watch_id, field_type, operator, value, alert_created, alert_updated) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}
	_, err = insert.Exec(alert.ID, alert.WatchId, alert.FieldType, alert.Operator, alert.Value, alert.AlertCreated, alert.AlertUpdated)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}

func DeleteWatch(watch model.WATCH) {
	_, _ = db.Exec("DELETE FROM pollerdb.watch WHERE watch_id = ?", watch.ID)
	_, _ = db.Exec("DELETE FROM pollerdb.alert WHERE watch_id = ?", watch.ID)
}
