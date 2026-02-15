package db

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	table   string
	columns []string
	where   []string
	args    []any
	limit   int
	offset  int
	orderBy string
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:   table,
		columns: []string{"*"},
		where:   make([]string, 0),
		args:    make([]any, 0),
	}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		qb.columns = columns
	}
	return qb
}

func (qb *QueryBuilder) Where(condition string, args ...any) *QueryBuilder {
	if condition != "" {
		qb.where = append(qb.where, condition)
		qb.args = append(qb.args, args...)
	}
	return qb
}

func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) BuildSelect() (string, []any) {
	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(qb.columns, ", "))
	sb.WriteString(" FROM ")
	sb.WriteString(qb.table)

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.where, " AND "))
	}

	if qb.orderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(qb.orderBy)
	}

	args := make([]any, len(qb.args))
	copy(args, qb.args)

	if qb.limit > 0 {
		sb.WriteString(" LIMIT ?")
		args = append(args, qb.limit)
	}

	if qb.offset > 0 {
		sb.WriteString(" OFFSET ?")
		args = append(args, qb.offset)
	}

	return sb.String(), args
}

func (qb *QueryBuilder) BuildCount() (string, []any) {
	var sb strings.Builder

	sb.WriteString("SELECT COUNT(*) FROM ")
	sb.WriteString(qb.table)

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.where, " AND "))
	}

	return sb.String(), qb.args
}

func (qb *QueryBuilder) BuildInsert(data map[string]any, returning ...string) (string, []any) {
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(qb.table)

	cols := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	args := make([]any, 0, len(data))

	for k, v := range data {
		cols = append(cols, k)
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}

	sb.WriteString(" (")
	sb.WriteString(strings.Join(cols, ", "))
	sb.WriteString(") VALUES (")
	sb.WriteString(strings.Join(placeholders, ", "))
	sb.WriteString(")")

	if len(returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(returning, ", "))
	}

	return sb.String(), args
}

func (qb *QueryBuilder) BuildUpdate(data map[string]any, returning ...string) (string, []any) {
	var sb strings.Builder

	sb.WriteString("UPDATE ")
	sb.WriteString(qb.table)
	sb.WriteString(" SET ")

	sets := make([]string, 0, len(data))
	args := make([]any, 0, len(data))

	for k, v := range data {
		sets = append(sets, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}

	sb.WriteString(strings.Join(sets, ", "))

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.where, " AND "))
		args = append(args, qb.args...)
	}

	if len(returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(returning, ", "))
	}

	return sb.String(), args
}

func (qb *QueryBuilder) BuildDelete() (string, []any) {
	var sb strings.Builder

	sb.WriteString("DELETE FROM ")
	sb.WriteString(qb.table)

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.where, " AND "))
	}

	return sb.String(), qb.args
}
