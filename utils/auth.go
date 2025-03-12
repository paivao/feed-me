package utils

import (
	"context"

	"github.com/feed-me/database"
)

var defaultHash string

func init() {
	hash, err := passwordHash("")
	if err != nil {
		panic(err)
	}
	defaultHash = hash
}

func GetUser(ctx context.Context, queries *database.Queries, username, password string) *database.User {
	user, err := queries.GetUserByName(ctx, username)
	var hash string
	if err != nil {
		hash = defaultHash
	} else {
		hash = user.PasswordHash
	}
	if passwordVerify(password, hash) {
		return &user
	}
	return nil
}

func UserCan(queries *database.Queries, user *database.User, permission string) bool {
	return false
}
