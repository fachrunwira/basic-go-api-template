package query

import "fmt"

func (qb *queryBuilder) Where(cond string, args ...interface{}) *queryBuilder {
	if len(qb.whereClause) > 0 {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("AND %s", cond))
	} else {
		qb.whereClause = append(qb.whereClause, cond)
	}

	qb.args = append(qb.args, args...)
	return qb
}

func (qb *queryBuilder) OrWhere(cond string, args ...interface{}) *queryBuilder {
	if len(qb.whereClause) > 0 {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("AND %s", cond))
	} else {
		qb.whereClause = append(qb.whereClause, cond)
	}

	qb.args = append(qb.args, args...)
	return qb
}

func (qb *queryBuilder) WhereRaw(raw string) *queryBuilder {
	if len(qb.whereClause) > 0 {
		qb.whereClause = append(qb.whereClause, raw)
	} else {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("AND %s", raw))
	}

	return qb
}

func (qb *queryBuilder) OrWhereRaw(raw string) *queryBuilder {
	if len(qb.whereClause) > 0 {
		qb.whereClause = append(qb.whereClause, raw)
	} else {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("OR %s", raw))
	}

	return qb
}
