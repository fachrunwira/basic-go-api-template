package query

import "fmt"

func (qb *queryBuilder) OrderBy(fields string, sortBy ...string) *queryBuilder {
	sort := "ASC"
	if len(sortBy) > 0 {
		sort = sortBy[0]
	}

	qb.orderClause = append(qb.orderClause, fmt.Sprintf("%s %s", fields, sort))

	return qb
}
