package util

import "github.com/lib/pq"

// SanitizeSQLLit should be called on all values that will be passed as literals
// (arguments) to SQL queries
func SanitizeSQLLit(qArg string) string {
	return pq.QuoteLiteral(qArg)
}
