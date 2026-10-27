package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFileDoesNotOverrideEnvironment(t *testing.T) {
	t.Setenv("CONFIG_TEST_EXISTING", "from-environment")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("CONFIG_TEST_EXISTING=from-file\nCONFIG_TEST_FILE=\"from-file\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := loadDotEnvFile(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("CONFIG_TEST_EXISTING"); got != "from-environment" {
		t.Fatalf("existing environment value changed to %q", got)
	}
	if got := os.Getenv("CONFIG_TEST_FILE"); got != "from-file" {
		t.Fatalf("file value = %q, want %q", got, "from-file")
	}
}

func TestDatabaseDSNEscapesCredentials(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user@example.com",
		Password: "p@ss/word",
		DBName:   "cdk_recharge",
		SSLMode:  "require",
	}
	if got, want := cfg.DSN(), "postgres://user%40example.com:p%40ss%2Fword@localhost:5432/cdk_recharge?sslmode=require"; got != want {
		t.Fatalf("DSN() = %q, want %q", got, want)
	}
}
