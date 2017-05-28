package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=0.0.0.0 user=postgres password=password  dbname=gofit sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not create database: %v", err)
	}
	return db, nil
}
