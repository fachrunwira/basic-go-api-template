package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/fachrunwira/basic-go-api-template/config"
)

type queryBuilder struct {
	db          *sql.DB
	TableName   string
	TableAlias  string
	Fields      []string
	WhereClause []string
	OrderClause []string
	GroupClause []string
	PageSize    int
	OffsetSize  int
	args        []interface{}
}

func (qb *queryBuilder) Table(table string, alias ...string) *queryBuilder {
	qb.TableName = table
	if len(alias) > 0 {
		qb.TableAlias = alias[0]
	}
	return qb
}

func (qb *queryBuilder) Select(fields ...string) *queryBuilder {
	qb.Fields = append(qb.Fields, fields...)
	return qb
}

func (qb *queryBuilder) Where(cond string, args ...interface{}) *queryBuilder {
	if len(qb.WhereClause) > 0 {
		qb.WhereClause = append(qb.WhereClause, fmt.Sprintf("AND %s", cond))
	} else {
		qb.WhereClause = append(qb.WhereClause, cond)
	}

	qb.args = append(qb.args, args...)
	return qb
}

func (qb *queryBuilder) OrWhere(cond string, args ...interface{}) *queryBuilder {
	if len(qb.WhereClause) > 0 {
		qb.WhereClause = append(qb.WhereClause, fmt.Sprintf("AND %s", cond))
	} else {
		qb.WhereClause = append(qb.WhereClause, cond)
	}

	qb.args = append(qb.args, args...)
	return qb
}

func (qb *queryBuilder) WhereRaw(raw string) *queryBuilder {
	if len(qb.WhereClause) > 0 {
		qb.WhereClause = append(qb.WhereClause, raw)
	} else {
		qb.WhereClause = append(qb.WhereClause, fmt.Sprintf("AND %s", raw))
	}

	return qb
}

func (qb *queryBuilder) OrWhereRaw(raw string) *queryBuilder {
	if len(qb.WhereClause) > 0 {
		qb.WhereClause = append(qb.WhereClause, raw)
	} else {
		qb.WhereClause = append(qb.WhereClause, fmt.Sprintf("OR %s", raw))
	}

	return qb
}

func (qb *queryBuilder) GroupBy(fields ...string) *queryBuilder {
	qb.GroupClause = append(qb.GroupClause, fields...)
	return qb
}

func (qb *queryBuilder) OrderBy(fields string, sortBy ...string) *queryBuilder {
	sort := "ASC"
	if len(sortBy) > 0 {
		sort = sortBy[0]
	}

	qb.OrderClause = append(qb.OrderClause, fmt.Sprintf("%s %s", fields, sort))

	return qb
}

func (qb *queryBuilder) Limit(size ...int) *queryBuilder {
	pageSize := 10
	if len(size) > 0 {
		pageSize = size[0]
	}

	qb.PageSize = pageSize
	return qb
}

func (qb *queryBuilder) Offset(page int) *queryBuilder {
	qb.OffsetSize = (page - 1)
	return qb
}

func Builder() *queryBuilder {
	return &queryBuilder{}
}

func (qb *queryBuilder) initGetRows() (string, []interface{}) {
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.Fields, ", "), qb.TableName)
	if qb.TableAlias != "" {
		query += fmt.Sprintf(" AS %s", qb.TableAlias)
	}

	if len(qb.WhereClause) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.WhereClause, " "))
	}

	if len(qb.GroupClause) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", strings.Join(qb.GroupClause, ", "))
	}

	if len(qb.OrderClause) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", strings.Join(qb.OrderClause, ", "))
	}

	if qb.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.PageSize)
	}

	if qb.OffsetSize > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.OffsetSize)
	}

	return query, qb.args
}

func (qb *queryBuilder) Get(dest ...any) error {
	query, args := qb.initGetRows()
	row := qb.db.QueryRow(query, args...)
	return row.Scan(dest...)
}

func (qb *queryBuilder) First(dest ...any) error {
	qb.PageSize = 1
	query, args := qb.initGetRows()
	row := qb.db.QueryRow(query, args...)
	return row.Scan(dest...)
}

func withTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		DBLogger.Printf("Failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			DBLogger.Printf("Recovered from panic, rolled back transaction: %v", err)
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		DBLogger.Printf("Transaction function returned, rolling back transaction: %v", err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		DBLogger.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (qb *queryBuilder) Insert(attributes map[string]string) error {
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

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUE (%s);", qb.TableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

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

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s;`, qb.TableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
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
	query := fmt.Sprintf("UPDATE %s SET %s", qb.TableName, strings.Join(placeholders, ", "))
	if len(qb.WhereClause) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.WhereClause, " "))
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
