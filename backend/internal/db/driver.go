package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/stdlib"
)

type rebindDriver struct {
	inner driver.Driver
}

type rebindConn struct {
	inner driver.Conn
}

type rebindStmt struct {
	inner driver.Stmt
}

var registerPostgresDriver sync.Once

func openPostgresDriver() string {
	name := "pgx-rebind"
	registerPostgresDriver.Do(func() {
		sql.Register(name, &rebindDriver{inner: stdlib.GetDefaultDriver()})
	})
	return name
}

func (d *rebindDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.inner.Open(name)
	if err != nil {
		return nil, err
	}
	return &rebindConn{inner: conn}, nil
}

func (c *rebindConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.inner.Prepare(rebindPlaceholders(query))
	if err != nil {
		return nil, err
	}
	return &rebindStmt{inner: stmt}, nil
}

func (c *rebindConn) Close() error {
	return c.inner.Close()
}

func (c *rebindConn) Begin() (driver.Tx, error) {
	return c.inner.Begin()
}

func (c *rebindConn) Ping(ctx context.Context) error {
	pinger, ok := c.inner.(driver.Pinger)
	if !ok {
		return nil
	}
	return pinger.Ping(ctx)
}

func (c *rebindConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	beginner, ok := c.inner.(driver.ConnBeginTx)
	if !ok {
		return nil, driver.ErrSkip
	}
	return beginner.BeginTx(ctx, opts)
}

func (c *rebindConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.inner.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return execer.ExecContext(ctx, rebindPlaceholders(query), args)
}

func (c *rebindConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.inner.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return queryer.QueryContext(ctx, rebindPlaceholders(query), args)
}

func (c *rebindConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	preparer, ok := c.inner.(driver.ConnPrepareContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	stmt, err := preparer.PrepareContext(ctx, rebindPlaceholders(query))
	if err != nil {
		return nil, err
	}
	return &rebindStmt{inner: stmt}, nil
}

func (c *rebindConn) CheckNamedValue(value *driver.NamedValue) error {
	checker, ok := c.inner.(driver.NamedValueChecker)
	if !ok {
		return nil
	}
	return checker.CheckNamedValue(value)
}

func (s *rebindStmt) Close() error {
	return s.inner.Close()
}

func (s *rebindStmt) NumInput() int {
	return s.inner.NumInput()
}

func (s *rebindStmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.inner.Exec(args)
}

func (s *rebindStmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.inner.Query(args)
}

func (s *rebindStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := s.inner.(driver.StmtExecContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return execer.ExecContext(ctx, args)
}

func (s *rebindStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := s.inner.(driver.StmtQueryContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return queryer.QueryContext(ctx, args)
}

func (s *rebindStmt) CheckNamedValue(value *driver.NamedValue) error {
	checker, ok := s.inner.(driver.NamedValueChecker)
	if !ok {
		return nil
	}
	return checker.CheckNamedValue(value)
}

func rebindPlaceholders(query string) string {
	var builder strings.Builder
	builder.Grow(len(query) + 8)
	placeholder := 0
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for index := 0; index < len(query); index++ {
		current := query[index]
		var next byte
		if index+1 < len(query) {
			next = query[index+1]
		}

		if inLineComment {
			builder.WriteByte(current)
			if current == '\n' {
				inLineComment = false
			}
			continue
		}
		if inBlockComment {
			builder.WriteByte(current)
			if current == '*' && next == '/' {
				builder.WriteByte(next)
				index++
				inBlockComment = false
			}
			continue
		}
		if !inDoubleQuote && !inBacktick && current == '-' && next == '-' {
			builder.WriteByte(current)
			builder.WriteByte(next)
			index++
			inLineComment = true
			continue
		}
		if !inSingleQuote && !inDoubleQuote && !inBacktick && current == '/' && next == '*' {
			builder.WriteByte(current)
			builder.WriteByte(next)
			index++
			inBlockComment = true
			continue
		}
		if !inDoubleQuote && !inBacktick && current == '\'' {
			builder.WriteByte(current)
			if inSingleQuote && next == '\'' {
				builder.WriteByte(next)
				index++
			} else {
				inSingleQuote = !inSingleQuote
			}
			continue
		}
		if !inSingleQuote && !inBacktick && current == '"' {
			inDoubleQuote = !inDoubleQuote
			builder.WriteByte(current)
			continue
		}
		if !inSingleQuote && !inDoubleQuote && current == '`' {
			inBacktick = !inBacktick
			builder.WriteByte(current)
			continue
		}
		if !inSingleQuote && !inDoubleQuote && !inBacktick && strings.HasPrefix(query[index:], "CURRENT_TIMESTAMP") {
			beforeIdentifier := index == 0 || !isSQLIdentifierByte(query[index-1])
			afterIndex := index + len("CURRENT_TIMESTAMP")
			afterIdentifier := afterIndex == len(query) || !isSQLIdentifierByte(query[afterIndex])
			if beforeIdentifier && afterIdentifier {
				if strings.HasPrefix(query[afterIndex:], "::") {
					builder.WriteString("CURRENT_TIMESTAMP")
				} else {
					builder.WriteString("CURRENT_TIMESTAMP::text")
				}
				index = afterIndex - 1
				continue
			}
		}
		if current == '?' && !inSingleQuote && !inDoubleQuote && !inBacktick {
			placeholder++
			builder.WriteString(fmt.Sprintf("$%d", placeholder))
			continue
		}
		builder.WriteByte(current)
	}
	return builder.String()
}

func isSQLIdentifierByte(value byte) bool {
	return value == '_' || value >= '0' && value <= '9' || value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z'
}
