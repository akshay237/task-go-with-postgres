package database

import (
	"fmt"
	"testing"
	"time"
)

func TestPostgresDBService(t *testing.T) {
	config := &ConfigPostgres{
		Host:     "127.0.0.1",
		User:     "akshay",
		Password: "password",
		Database: "akshay",
	}
	dbSvc := NewPostgresDBService(config)
	fmt.Println("DB Svc", dbSvc)
	conn, err := dbSvc.GetDBInstance()
	if err != nil {
		t.Log("Error", err)
	}
	query := `select task_id
	, task_title
	, task_description
	, task_priority
	, task_due_datetime
	, contact_name
	, is_deleted from tasks;`

	var task_id int
	var task_title, task_description, task_priority, contact_name string
	var is_deleted bool
	var timestamp time.Time

	rows := conn.QueryRow(query)
	err = rows.Scan(&task_id, &task_title, &task_description, &task_priority, &timestamp, &contact_name, &is_deleted)
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log("Got ")
}
