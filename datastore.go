package main

import (
	"context"
	"time"
)

type Height struct {
	ID     int `json:"id"`
	Height int `json:"height"`
}

type Shitlord struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestam"`
}

func setHeightInStore(user int, height int) error {
	tx, err := DS.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO height(id, height) values(?, ?)")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(user, height); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func getStoredHeight(uid int) (int, error) {
	var height int
	stmt, err := DS.Prepare("SELECT height FROM height WHERE id = ?")
	if err != nil {
		return 0, err
	}

	if err := stmt.QueryRow(uid).Scan(&height); err != nil {
		return 0, err
	}

	return height, nil
}

func setShitlord(shitlord Shitlord) error {
	tx, err := DS.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO shitlord(username, time) values(?, datetime('now'))")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(shitlord.ID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func getAllShitlords() (*[]Shitlord, error) {
	var shitlords []Shitlord

	stmt, err := DS.Prepare("SELECT * FROM shitlord")
	if err != nil {
		return &shitlords, err
	}

	if err := stmt.QueryRow().Scan(&shitlords); err != nil {
		return &shitlords, err
	}

	return &shitlords, nil
}

func getLastShitlord() (*Shitlord, error) {
	var shitlord Shitlord

	stmt, err := DS.Prepare("SELECT * FROM shitlords WHERE id = (SELECT MAX(id) FROM shitlords);")
	if err != nil {
		return &shitlord, err
	}

	if err := stmt.QueryRow().Scan(&shitlord); err != nil {
		return &shitlord, err
	}

	return &shitlord, nil
}

func setupDatastore() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS height (id INTEGER NOT NULL PRIMARY KEY, height INTEGER);
	CREATE TABLE IF NOT EXISTS shitlord (id INTEGER PRIMARY KEY AUTOINCREMENT, username INTEGER, time datetime);	
	`
	tx, err := DS.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(sqlStmt); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
