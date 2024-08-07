package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// WriteGaoge writes gauge metric into db
func (d *DB) WriteGauge(ctx context.Context, name string, value float64) (float64, error) {
	var oldVal float64

	tx, err := d.DB.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	// Check if metric exists
	if err := tx.QueryRow(ctx, `SELECT gauge FROM metrics WHERE name=$1 AND type = 'gauge'`, name).Scan(&oldVal); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			if _, err := tx.Exec(ctx, `INSERT INTO metrics (name, type, gauge) VALUES ($1, 'gauge', $2)`, name, value); err != nil {
				d.Logger.Errorw("write gauge error", "error", err)
				return 0, err
			}
			return 0, nil
		}

		d.Logger.Errorw("write gauge error", "error", err)
		return 0, err
	}

	// Update metric if exists
	if _, err := tx.Exec(ctx, `UPDATE metrics SET gauge = $1 WHERE type = 'gauge' AND name = $2`, value, name); err != nil {
		d.Logger.Errorw("write gauge error", "error", err)
		return 0, err
	}

	return value, nil
}

// WriteCounter writes counter metric into db
func (d *DB) WriteCounter(ctx context.Context, name string, value int64) (int64, error) {
	tx, err := d.DB.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	// Check if metric exists
	var oldVal int64
	if err := tx.QueryRow(ctx, "SELECT counter FROM metrics WHERE name = $1 AND type = 'counter'", name).Scan(&oldVal); err != nil {
		// Insert new metric if not exists
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := tx.Exec(ctx, `INSERT INTO metrics (name, type, counter) VALUES ($1, 'counter', $2)`, name, value); err != nil {
				d.Logger.Errorw("write counter error", "error", err)
				return 0, err
			}

			return value, nil
		}
		return 0, err
	}

	// Update metric if exists
	value += oldVal
	if _, err := tx.Exec(ctx, `UPDATE metrics SET counter = $1 WHERE type = 'counter' AND name = $2`, value, name); err != nil {
		d.Logger.Errorw("write counter error", "error", err)
		return 0, err
	}

	return value, nil
}
