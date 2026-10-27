package db

import "testing"

func TestRebindPlaceholders(t *testing.T) {
	query := `SELECT '?' AS literal, "?" AS identifier /* ? */ WHERE a = ? AND b = '?' AND c = CURRENT_TIMESTAMP -- ?
AND d = ?`
	want := `SELECT '?' AS literal, "?" AS identifier /* ? */ WHERE a = $1 AND b = '?' AND c = CURRENT_TIMESTAMP::text -- ?
AND d = $2`
	if got := rebindPlaceholders(query); got != want {
		t.Fatalf("rebindPlaceholders() = %q, want %q", got, want)
	}
}

func TestAdaptSchemaQueryForPostgres(t *testing.T) {
	oldDriver := databaseDriver
	t.Cleanup(func() { databaseDriver = oldDriver })
	databaseDriver = "postgres"

	query := `CREATE TABLE sample (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`
	want := `CREATE TABLE sample (id BIGSERIAL PRIMARY KEY, created_at TEXT DEFAULT CURRENT_TIMESTAMP::text)`
	if got := adaptSchemaQuery(query); got != want {
		t.Fatalf("adaptSchemaQuery() = %q, want %q", got, want)
	}
}
