package db

import (
	"database/sql"
	"testing"
)

func TestListCardplatformStoredCDKCodesPageUsesPortableOrdering(t *testing.T) {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec(`
		CREATE TABLE cardplatform_cdk_codes (
			upstream_id INTEGER,
			code TEXT NOT NULL UNIQUE,
			code_prefix TEXT,
			plan TEXT,
			fee_amount_minor INTEGER,
			status TEXT,
			payment_country TEXT,
			created_at TEXT
		);
		CREATE TABLE cardplatform_cdk_notes (upstream_id INTEGER PRIMARY KEY, note TEXT);
		INSERT INTO cardplatform_cdk_codes
			(upstream_id, code, code_prefix, plan, fee_amount_minor, status, created_at)
		VALUES (2, 'ZC-TWO', 'ZC-TWO', 'plus', 15, 'unused', '2026-10-26 10:00:00'),
		       (1, 'ZC-ONE', 'ZC-ONE', 'plus', 15, 'unused', '2026-10-26 10:00:00');
	`); err != nil {
		t.Fatal(err)
	}

	oldDB := DB
	oldDriver := databaseDriver
	DB = conn
	databaseDriver = "sqlite"
	t.Cleanup(func() {
		DB = oldDB
		databaseDriver = oldDriver
	})

	list, total, err := ListCardplatformStoredCDKCodesPage("plus", "two", "", 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(list) != 1 || list[0].Code != "ZC-TWO" {
		t.Fatalf("got total=%d list=%+v", total, list)
	}
}
