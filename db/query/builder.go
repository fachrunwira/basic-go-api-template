package query

import (
	"context"
	"database/sql"
	"log"

	"github.com/fachrunwira/basic-go-api-template/db"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"
)

type queryBuilder struct {
	db          *sql.DB
	tableName   string
	tableAlias  string
	fields      []string
	whereClause []string
	orderClause []string
	groupClause []string
	limit       int
	page        int
	args        []interface{}
}

var dbLogger *log.Logger = logger.SetLogger("./storage/log/db.log")

func Builder(ctx context.Context) *queryBuilder {
	return &queryBuilder{
		db:    db.FromContext(ctx),
		limit: 15,
		page:  1,
	}
}
