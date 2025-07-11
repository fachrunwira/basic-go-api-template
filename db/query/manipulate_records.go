package query

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/fachrunwira/basic-go-api-template/config"
)

func (qb *queryBuilder) Insert(attributes map[string]interface{}) error {
	var columnSets = make(map[string]struct{})

	for k := range attributes {
		columnSets[k] = struct{}{}
	}

	var columns []string
	for col := range columnSets {
		columns = append(columns, col)
	}
	sort.Strings(columns)

	var (
		placeholders []string
		args         []interface{}
		paramIndex   = 1
	)

	ph := make([]string, len(columns))
	for i, col := range columns {
		val, ok := attributes[col]
		if !ok {
			args = append(args, nil)
		} else {
			args = append(args, val)
		}
		ph[i] = insertType(paramIndex)
		paramIndex++
	}

	placeholders = append(placeholders, strings.Join(ph, ", "))

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUE (%s);", qb.tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return withTransaction(qb.db, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(args...)
		if err != nil {
			return err
		}

		return nil
	})
}

func (qb *queryBuilder) BatchInsert(attributes []map[string]interface{}) error {
	var columnSets = make(map[string]struct{})
	for _, attr := range attributes {
		for k := range attr {
			columnSets[k] = struct{}{}
		}
	}

	var columns []string
	for col := range columnSets {
		columns = append(columns, col)
	}
	sort.Strings(columns)

	var (
		placeholders []string
		args         []interface{}
		paramIndex   = 1
	)

	for _, row := range attributes {
		ph := make([]string, len(columns))
		for i, col := range columns {
			val, ok := row[col]
			if !ok {
				args = append(args, nil)
			} else {
				args = append(args, val)
			}

			ph[i] = insertType(paramIndex)
			paramIndex++
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(ph, ", ")))
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s;`, qb.tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	return withTransaction(qb.db, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(args...)
		if err != nil {
			return err
		}

		return nil
	})
}

func (qb *queryBuilder) Update(attributes map[string]interface{}) error {
	var columnSets = make(map[string]struct{})
	for k := range attributes {
		columnSets[k] = struct{}{}
	}

	var columns []string
	for k := range columnSets {
		columns = append(columns, k)
	}

	var (
		placeholders []string
		args         []interface{}
		paramIndex   = 1
	)

	ph := make([]string, len(columns))
	for i, col := range columns {
		val, ok := attributes[col]
		if !ok {
			args = append(args, nil)
		} else {
			args = append(args, val)
		}
		ph[i] = fmt.Sprintf("%s = %s", col, insertType(paramIndex))
	}

	placeholders = append(placeholders, strings.Join(ph, ", "))
	query := fmt.Sprintf("UPDATE %s SET %s", qb.tableName, strings.Join(placeholders, ", "))
	if len(qb.whereClause) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.whereClause, " "))
		args = append(args, qb.args...)
	}

	return withTransaction(qb.db, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(args...)
		if err != nil {
			return err
		}

		return nil
	})
}

func insertType(index int) string {
	db_config := config.LoadDBConfig()

	switch db_config.Driver {
	case "mysql":
		return "?"
	case "pgsql":
		return fmt.Sprintf("$%d", index)
	default:
		return ""
	}
}
