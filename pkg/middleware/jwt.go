package middleware

import "strings"

func getJWTToken(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
