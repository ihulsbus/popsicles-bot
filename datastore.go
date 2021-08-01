package main

import (
	"context"
)

type Height struct {
	ID     string `json:"id"`
	Height int    `json:"height"`
}

func setHeight(user int, height int) error {
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

func setupDatastore() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS height (id INTEGER NOT NULL PRIMARY KEY, height INTEGER);
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
