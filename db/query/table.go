package query

func (qb *queryBuilder) Table(table string, alias ...string) *queryBuilder {
	qb.tableName = table
	if len(alias) > 0 {
		qb.tableAlias = alias[0]
	}
	return qb
}
