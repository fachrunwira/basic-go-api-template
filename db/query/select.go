package query

func (qb *queryBuilder) Select(fields ...string) *queryBuilder {
	qb.fields = append(qb.fields, fields...)
	return qb
}
