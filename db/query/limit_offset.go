package query

func (qb *queryBuilder) Limit(size ...int) *queryBuilder {
	pageSize := 10
	if len(size) > 0 {
		pageSize = size[0]
	}

	qb.pageSize = pageSize
	return qb
}

func (qb *queryBuilder) Offset(page int) *queryBuilder {
	qb.offsetSize = (page - 1)
	return qb
}
