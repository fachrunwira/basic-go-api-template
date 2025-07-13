package query

func (qb *queryBuilder) Limit(size int) *queryBuilder {
	if size > 0 {
		qb.limit = size
	} else {
		qb.limit = 15
	}
	return qb
}

func (qb *queryBuilder) Page(page int) *queryBuilder {
	if page > 0 {
		qb.page = page
	} else {
		qb.page = 1
	}
	return qb
}
