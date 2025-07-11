package query

import (
	"database/sql"
	"fmt"
	"strings"
)

func (qb *queryBuilder) initGetRows() (string, []interface{}) {
	query := "SELECT"

	if len(qb.fields) > 0 {
		query += fmt.Sprintf(" %s FROM %s", strings.Join(qb.fields, ", "), qb.tableName)
	} else {
		query += fmt.Sprintf(" * FROM %s", qb.tableName)
	}

	if qb.tableAlias != "" {
		query += fmt.Sprintf(" AS %s", qb.tableAlias)
	}

	if len(qb.whereClause) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.whereClause, " "))
	}

	if len(qb.groupClause) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", strings.Join(qb.groupClause, ", "))
	}

	if len(qb.orderClause) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", strings.Join(qb.orderClause, ", "))
	}

	if qb.pageSize > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.pageSize)
	}

	if qb.offsetSize > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offsetSize-1)
	}

	return query, qb.args
}

func (qb *queryBuilder) First() (map[string]interface{}, error) {
	qb.pageSize = 1
	query, args := qb.initGetRows()

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, sql.ErrNoRows
	}

	columns, err := row.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuesPtrs := make([]interface{}, len(columns))

	for k := range columns {
		valuesPtrs[k] = &values[k]
	}

	if err := row.Scan(valuesPtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]

		if b, ok := val.([]byte); ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return result, nil
}

func (qb *queryBuilder) Get() ([]map[string]interface{}, error) {
	query, args := qb.initGetRows()

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuesPtrs[i] = &values[i]
		}

		if err := rows.Scan(valuesPtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		result = append(result, row)
	}

	return result, nil
}
