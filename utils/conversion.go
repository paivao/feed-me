package utils

import "database/sql"

func ConvertToSqlNull(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *s}
}

func BoolColapse(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}
