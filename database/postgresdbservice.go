package database

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

type PostgresDBService struct {
	config     *ConfigPostgres
	dbinstance *sql.DB
	lock       *sync.Mutex
}

func (o *PostgresDBService) createNewDBInstance() (*sql.DB, error) {
	// 1. Convert PostgreSQL Config to DB Config.
	constr := "user=akshay dbname=akshay password=password sslmode=disable"
	dbinst, err := sql.Open("postgres", constr)
	if err != nil {
		return nil, err
	}

	err = dbinst.Ping()
	if err != nil {
		dbinst.Close()
		return nil, err
	}

	return dbinst, nil
}

func (o *PostgresDBService) GetDBInstance() (*sql.DB, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.dbinstance != nil {
		return o.dbinstance, nil
	}

	dbinst, err := o.createNewDBInstance()
	if err != nil {
		return nil, err
	}
	o.dbinstance = dbinst
	return o.dbinstance, nil
}

func (o *PostgresDBService) ClearDBInstance() (bool, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.dbinstance == nil {
		return true, nil
	}

	err := o.dbinstance.Close()
	o.dbinstance = nil
	return true, err
}

func NewPostgresDBService(config *ConfigPostgres) *PostgresDBService {
	return &PostgresDBService{
		config:     config,
		dbinstance: nil,
		lock:       &sync.Mutex{},
	}
}
