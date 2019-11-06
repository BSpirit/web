package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func SelectQuery(table string, entries map[string][]string) (string, []interface{}) {
	var values []interface{}
	var whereClause []string
	query := "SELECT * FROM " + table
	if len(entries) != 0 {
		query += " WHERE "
		for key, value := range entries {
			if value[0] != "" {
				values = append(values, value[0])
				whereClause = append(whereClause, fmt.Sprintf("%s LIKE '%%' || ? || '%%'", key))
			}
		}
	}

	query = query + strings.Join(whereClause, " AND ")

	return query, values
}

func NewNullInt64(s string) sql.NullInt64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}
