package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
)

var ErrFieldsAreEmpty = "the fields are empty"

func NullIfEmpty(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func BuildSQLInsert(table string, fields []string) string {
	if len(fields) == 0 {
		return ErrFieldsAreEmpty
	}

	args := bytes.Buffer{}
	values := bytes.Buffer{}
	k := 0

	for _, v := range fields {
		k++
		args.WriteString(v)
		args.WriteString(", ")
		values.WriteString(fmt.Sprintf("$%d, ", k))
	}

	args.Truncate(args.Len() - 2)
	values.Truncate(values.Len() - 2)

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, args.String(), values.String())
}

func BuildSQLUpdatedByID(table string, fields []string) string {
	if len(fields) == 0 {
		return ErrFieldsAreEmpty
	}

	//move ID to the last
	fields = append(fields[1:], fields[0])
	args := bytes.Buffer{}
	k := 0

	/* The continue statement is used when you want to skip the remaining portion
	of the loop, and return to the top of the loop and continue a new iteration. */
	for _, v := range fields {
		if v == "created_at" {
			continue
		}
		args.WriteString(fmt.Sprintf("%s = $%d, ", v, k+1))
		k++
	}
	args.Truncate(args.Len() - 2)

	return fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, args.String(), k)
}

func BuildSQLDelete(table string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
}

func BuildSQLSelect(table string, fields []string) string {
	if len(fields) == 0 {
		return ErrFieldsAreEmpty
	}

	args := bytes.Buffer{}
	for _, v := range fields {
		args.WriteString(fmt.Sprintf("%s, ", v))
	}
	args.Truncate(args.Len() - 2)
	return fmt.Sprintf("SELECT %s FROM %s", args.String(), table)
}

func Int64ToNull(d int64) sql.NullInt64 {
	r := sql.NullInt64{Int64: d}
	if d > 0 {
		r.Valid = true
	}

	return r
}
