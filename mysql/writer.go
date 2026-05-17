package mysql

import (
	"adv_analitics_solution/internal/stats"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	insertQuery = `INSERT INTO %s (ts, country, os, browser, campaign_id, requests, impressions) VALUES (?, ?, ?, ?, ?, ?, ?)`
)

type writer struct {
	db        *sql.DB
	tableName string
}

func (w *writer) Insert(rows stats.Rows) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := w.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(insertQuery, w.tableName))
	if err != nil {
		return err
	}

	for k, v := range rows {
		ts := time.Unix(k.Timestamp)
		if _, err := stmt.Exec(ts, k.Country, k.Os, k.Browser, k.CampaignId, v.Requests, v.Impressions); err != nil {
			return err
		}
	}

	return tx.Commit()

}
