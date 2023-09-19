package jsonb

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func RawMessage(v any) (postgres.Jsonb, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return postgres.Jsonb{}, err
	}
	return postgres.Jsonb{RawMessage: json.RawMessage(b)}, nil
}

func Filter(filters map[string][]string) string {
	clause := make([]string, 0, len(filters))
	for key, vals := range filters {
		clause = append(clause, fmt.Sprintf("data->>'%s' IN ('%v')", key, strings.Join(vals, "','")))
	}
	return strings.Join(clause, " AND ")
}

func Search(searches map[string]string, operator string) string {
	clause := make([]string, 0, len(searches))
	for key, val := range searches {
		if val != "" {
			clause = append(clause, fmt.Sprintf("UPPER(data->>'%s') LIKE '%%%v%%'", key, strings.ToUpper(val)))
		}
	}

	return strings.Join(clause, fmt.Sprintf(" %s ", strings.ToUpper(operator)))
}
