package sql

import (
	"database/sql"
	"testing"

	"github.com/shoenig/test/must"
)

func DumpTable(t testing.TB, db *sql.DB, tableName string) {
	t.Helper()

	rows, err := db.Query(
		// WARNING: this is basically sql injection, but this is only a test utility so it should
		// always be executed in a trusted context.
		`SELECT * FROM ` + tableName,
	)
	must.NoError(t, err)
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	must.NoError(t, err)

	t.Logf("=== DUMPING TABLE %s ===", tableName)

	for rows.Next() {
		targets := make([]any, len(colTypes))
		targetPtrs := make([]any, 0, len(colTypes))
		for i := range targets {
			targetPtrs = append(targetPtrs, &targets[i])
		}

		err := rows.Scan(targetPtrs...)
		must.NoError(t, err)
		t.Log(targets...)
	}

	must.NoError(t, rows.Err())
	must.NoError(t, rows.Close())
}
