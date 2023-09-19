package jsonb

import (
	"fmt"
	"strings"
)

type conditions struct {
	list []string
	err  error
}

func (c conditions) clause() (string, error) {
	if len(c.list) == 0 {
		return "", nil
	}
	if c.err != nil {
		return "", fmt.Errorf("conditions: %w", c.err)
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(c.list, " AND ")), nil
}

// ensure always a new underlying array
func (c *conditions) add(cond string) {
	l := len(c.list)
	conds := make([]string, l, l+1)

	copy(conds, c.list)
	conds = append(conds, cond)
	c.list = conds
}

func (c *conditions) addErr(err error) {
	if c.err == nil {
		c.err = err
		return
	}
	c.err = fmt.Errorf("%v; %w", c.err, err)
}

// Match adds the condition that ALL key/value pairs in fields must be a subset of the model.
//
// Examples for fields:
//
//	// models with field "one" equal to "value"
//	map[string]any{"one": "value"}
//	// models with field "slice" equal to a slice containing "a"
//	map[string]any{"slice": []string{"a"}}
//	// models with field "slice" equal to a slice containing both "a" and "b"
//	map[string]any{"slice": []string{"a", "b"}}
//	// models with field "one" equal to "value" and "slice" equal to a slice containing "c"
//	map[string]any{"one": "value", "slice": []string{"c"}}
//	// models with field "slice" equal to a non-nil slice
//	map[string]any{"slice": []string{}}
//	// models with field "slice" equal to nil
//	map[string]any{"slice": nil}
//	// models with field "map" equal to a map or struct with the previous rules
//	map[string]any{"map": map[string]any{...}}
func (d DB[M]) Match(fields map[string]any) *DB[M] {
	match, err := parse("", fields)
	if err != nil {
		d.conds.addErr(err)
	}

	d.conds.add(fmt.Sprintf("doc @> '%s'::jsonb", match.doc))
	return &d
}

// Filter adds the condition that for each key in fields, the model must have a corresponding field equal to one of the values.
// Empty or nil slices are ignored.
func (d DB[M]) Filter(fields map[string][]string) *DB[M] {
	conds := make([]string, 0, len(fields))
	for field, vals := range fields {
		if len(vals) > 0 {
			conds = append(conds, fmt.Sprintf("doc->>'%s' IN ('%v')", field, strings.Join(vals, "','")))
		}
	}
	if len(conds) > 0 {
		d.conds.add(strings.Join(conds, " AND "))
	}
	return &d
}

// Search adds the condition that for at least one key in fields, the model must have a corresponding field containing the value.
// Comparisons are case-insensitive. Blank values for a field are ignored.
//
// NOTE: The `_` character is a wildcard, escape with square brackets `[_]` for literal value.
func (d DB[M]) Search(fields map[string]string) *DB[M] {
	conds := make([]string, 0, len(fields))
	for field, val := range fields {
		if val != "" {
			conds = append(conds, fmt.Sprintf("UPPER(doc->>'%s') LIKE '%%%v%%'", field, strings.ToUpper(val)))
		}
	}
	if len(conds) > 0 {
		d.conds.add(fmt.Sprintf("(%s)", strings.Join(conds, " OR ")))
	}
	return &d
}
