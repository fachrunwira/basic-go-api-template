package query

import (
	"fmt"
	"strings"
)

func (qb *queryBuilder) countRows() (*int, error) {
	query := fmt.Sprintf("SELECT count(*) as total FROM %s", qb.tableName)

	if len(qb.whereClause) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.whereClause, ", "))
	}

	if len(qb.groupClause) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", strings.Join(qb.groupClause, ", "))
	}

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.Query(qb.args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var total int
	for row.Next() {
		if err := row.Scan(&total); err != nil {
			return nil, err
		}
	}

	return &total, nil
}

func (qb *queryBuilder) Paginate() (map[string]interface{}, error) {
	if qb.pageSize == 0 {
		qb.pageSize = 15
	}

	if qb.offsetSize == 0 {
		qb.offsetSize = 1
	}

	total_row, err := qb.countRows()
	if err != nil {
		return nil, err
	}

	var paginate = make(map[string]interface{})

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

	results := []map[string]interface{}{}
	for rows.Next() {
		var (
			values     = make([]interface{}, len(columns))
			valuesPtrs = make([]interface{}, len(columns))
			rowMap     = make(map[string]interface{})
		)

		for i := range columns {
			valuesPtrs[i] = &values[i]
		}

		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		results = append(results, rowMap)
	}

	paginate["total_rows"] = total_row

	return nil, nil
}
