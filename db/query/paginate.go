package query

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/fachrunwira/basic-go-api-template/lib/env"
	"github.com/labstack/echo/v4"
)

type navLinks struct {
	Active bool
	Label  string
	Url    string
}

func (qb *queryBuilder) countRows() (*int, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s", qb.tableName)

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

	var total int
	err = stmt.QueryRow(qb.args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &total, nil
}

func navigationLinks(c echo.Context, page, last_page int) []interface{} {
	var links []navLinks
	showBefore, showAfter := 5, 5
	start := max(1, page-showBefore)
	end := min(last_page, page+showAfter)

	if last_page <= (showAfter + showBefore + 3) {
		start = 1
		end = last_page
	}

	if page > 1 {
		links = append(links, navLinks{Active: false, Label: "&laquo Prev", Url: setNavigation(c, page-1)})
	} else {
		links = append(links, navLinks{Active: false, Label: "&laquo Prev", Url: ""})
	}

	if start > 1 {
		links = append(links, navLinks{Active: false, Label: "1", Url: setNavigation(c, 1)})
		if start > 2 {
			links = append(links, navLinks{Active: false, Label: "...", Url: ""})
		}
	}

	for i := start; i <= end; i++ {
		links = append(links, navLinks{Active: i == page, Label: strconv.Itoa(i), Url: setNavigation(c, i)})
	}

	if end < last_page {
		if end < last_page-1 {
			links = append(links, navLinks{Active: false, Label: "...", Url: ""})
		}
		links = append(links, navLinks{Active: false, Label: strconv.Itoa(last_page), Url: setNavigation(c, last_page)})
	}

	if last_page > page {
		links = append(links, navLinks{Active: false, Label: "&raquo", Url: setNavigation(c, page+1)})
	} else {
		links = append(links, navLinks{Active: false, Label: "&raquo", Url: ""})
	}

	newLinks := []interface{}{}
	for _, v := range links {
		list := map[string]interface{}{
			"active": v.Active,
			"label":  v.Label,
			"url":    v.Url,
		}
		newLinks = append(newLinks, list)
	}

	return newLinks
}

func setNavigation(c echo.Context, page int) string {
	host := env.Get("APP_URL", fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host))
	path := c.Request().URL.Path

	queryParams := setUri(c)
	queryParams["page"] = strconv.Itoa(page)

	var newQueryParams []string
	for k, v := range queryParams {
		newQueryParams = append(newQueryParams, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("%s%s?%s", host, path, strings.Join(newQueryParams, "&"))
}

func setUri(c echo.Context) map[string]string {
	uri := c.Request().RequestURI
	var queryParams = make(map[string]string)

	if idx := strings.Index(uri, "?"); idx != -1 {
		params := strings.Split(uri[idx+1:], "&")
		for _, param := range params {
			val := strings.SplitN(param, "=", 2)
			key := val[0]
			value := ""
			if len(val) > 1 {
				value = val[1]
			}

			queryParams[key] = value
		}
	}

	return queryParams
}

func pageLinks(c echo.Context, links string) string {
	host := env.Get("APP_URL", fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host))
	path := c.Request().URL.Path

	queryParams := setUri(c)

	page := 1
	if val, ok := queryParams["page"]; ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			if links == "next" {
				page = parsed + 1
			} else {
				page = parsed - 1
			}
		}
	} else {
		page = 2
	}

	var newQueryParams []string
	queryParams["page"] = strconv.Itoa(page)
	for k, v := range queryParams {
		newQueryParams = append(newQueryParams, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("%s%s?%s", host, path, strings.Join(newQueryParams, "&"))
}

func prevPage(page int, c echo.Context) *string {
	if page > 1 {
		url := pageLinks(c, "prev")
		return &url
	}

	return nil
}

func nextPage(page, last_page int, c echo.Context) *string {
	if last_page > page {
		url := pageLinks(c, "next")
		return &url
	}

	return nil
}

func (qb *queryBuilder) Paginate(c echo.Context) (map[string]interface{}, error) {
	total, err := qb.countRows()
	if err != nil {
		return nil, err
	}

	offset := (qb.page - 1) * qb.limit
	last_page := math.Ceil(float64(*total) / float64(qb.limit))
	query, args := qb.initGetRows()
	query += fmt.Sprintf(" LIMIT %d OFFSET %d;", qb.limit, offset)

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
		values := make([]interface{}, len(columns))
		valuesPtrs := make([]interface{}, len(columns))
		rowMap := make(map[string]interface{})

		for i := range columns {
			valuesPtrs[i] = &values[i]
		}

		if err = rows.Scan(valuesPtrs...); err != nil {
			return nil, err
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

	paginate := map[string]interface{}{
		"total":         total,
		"from":          offset + 1,
		"to":            qb.page * qb.limit,
		"data":          results,
		"last_page":     last_page,
		"current_page":  qb.page,
		"next_page_url": nextPage(qb.page, int(last_page), c),
		"prev_page_url": prevPage(qb.page, c),
		"links":         navigationLinks(c, qb.page, int(last_page)),
	}

	return paginate, nil
}
