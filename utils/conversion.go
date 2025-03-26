package utils

import (
	"database/sql"
	"time"
)

func ConvertToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *s}
}

func ConvertToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Valid: true, Time: *t}
}

func BoolColapse(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}
