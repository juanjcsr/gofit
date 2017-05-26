package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=gofit")
	if err != nil {
		return nil, fmt.Errorf("could not create database")
	}
	return db, nil
}
