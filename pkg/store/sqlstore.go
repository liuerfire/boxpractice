package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SQLStore struct {
	db *sqlx.DB
}

func NewSQLStore() (*SQLStore, error) {
	u, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		u = "mysql://root@tcp(127.0.0.1:3306)/boxpractice"
	}
	dsn := fmt.Sprintf("%s?parseTime=true", strings.TrimPrefix(u, "mysql://"))
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	return &SQLStore{
		db: db,
	}, nil
}

func (s *SQLStore) DB() *sqlx.DB {
	return s.db
}

func (s *SQLStore) Close() {
	s.db.Close()
}

func (s *SQLStore) Cleanup() error {
	var tables []string
	if err := s.DB().Select(&tables, "SHOW TABLES"); err != nil {
		return err
	}
	for _, table := range tables {
		if _, err := s.DB().Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)); err != nil {
			return err
		}
	}
	return nil
}

func IsErrDuplicateEntry(err error) bool {
	var mErr *mysql.MySQLError
	if errors.As(err, &mErr) {
		if mErr.Number == 1062 {
			return true
		}
	}
	return false
}

func IsErrNotFound(err error) bool {
	return err == sql.ErrNoRows
}
