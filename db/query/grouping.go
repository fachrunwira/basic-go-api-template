package query

func (qb *queryBuilder) GroupBy(fields ...string) *queryBuilder {
	qb.groupClause = append(qb.groupClause, fields...)
	return qb
}
