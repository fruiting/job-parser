package pgsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"fruiting/job-parser/internal"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type storageSuite struct {
	suite.Suite

	ctx        context.Context
	db         *sql.DB
	mockDb     sqlmock.Sqlmock
	columnRows []string
	testErr    error

	position     internal.Name
	minSalary    internal.Salary
	maxSalary    internal.Salary
	medianSalary internal.Salary
	parser       internal.Parser

	storage *Storage
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &storageSuite{})
}

func (s *storageSuite) SetupTest() {
	s.ctx = context.Background()
	var err error
	s.db, s.mockDb, err = sqlmock.New()
	s.NoError(err)
	//sqlDb := sqlx.NewDb(s.db, "sqlmock")
	s.testErr = errors.New("test err")

	s.position = "test position"
	s.minSalary = internal.Salary(1)
	s.maxSalary = internal.Salary(2)
	s.medianSalary = internal.Salary(1)
	s.parser = "test parser"

	s.storage = NewStorage()
}

func (s *storageSuite) TestSetErr() {
	s.mockDb.
		ExpectExec(regexp.QuoteMeta("call some_procedure($1,$2,$3,$4,$5)")).
		WithArgs(s.position, s.minSalary, s.maxSalary, s.medianSalary, s.parser).
		WillReturnError(s.testErr)

	err := s.storage.Set(s.ctx, s.position, s.minSalary, s.maxSalary, s.medianSalary, s.parser)

	s.Equal(fmt.Errorf("can't set jobs info: %w", s.testErr), err)
}

func (s *storageSuite) TestSetOk() {
	s.mockDb.
		ExpectExec(regexp.QuoteMeta("call some_procedure($1,$2,$3,$4,$5)")).
		WithArgs(s.position, s.minSalary, s.maxSalary, s.medianSalary, s.parser).
		WillReturnResult(driver.ResultNoRows)

	err := s.storage.Set(s.ctx, s.position, s.minSalary, s.maxSalary, s.medianSalary, s.parser)

	s.Nil(err)
}
