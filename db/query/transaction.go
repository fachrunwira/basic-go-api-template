package query

import "database/sql"

func withTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		dbLogger.Printf("Failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			dbLogger.Printf("Recovered from panic, rolled back transaction: %v", err)
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		dbLogger.Printf("Transaction function returned, rolling back transaction: %v", err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		dbLogger.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
