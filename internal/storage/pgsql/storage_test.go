package pgsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"fruiting/job-parser/internal"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type storageSuite struct {
	suite.Suite

	ctx        context.Context
	db         *sql.DB
	mockDb     sqlmock.Sqlmock
	columnRows []string
	logs       *observer.ObservedLogs
	testErr    error

	position     internal.Name
	minSalary    internal.Salary
	maxSalary    internal.Salary
	medianSalary internal.Salary
	parser       internal.Parser
	fromYear     uint16
	toYear       uint16

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
	sqlDb := sqlx.NewDb(s.db, "sqlmock")
	s.columnRows = []string{
		"position_to_parse",
		"min_salary",
		"max_salary",
		"median_salary",
		"popular_skills",
		"parser",
		"mdate",
	}
	core, logs := observer.New(zap.InfoLevel)
	s.logs = logs
	s.testErr = errors.New("test err")

	s.position = "test position"
	s.minSalary = internal.Salary(1)
	s.maxSalary = internal.Salary(2)
	s.medianSalary = internal.Salary(1)
	s.parser = "test parser"
	s.fromYear = 2020
	s.toYear = 2023

	s.storage = NewStorage(sqlDb, zap.New(core))
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

func (s *storageSuite) TestGetQueryErr() {
	s.mockDb.
		ExpectQuery(regexp.QuoteMeta("select * from some_procedure($1,$2,$3)")).
		WithArgs(
			s.position,
			s.fromYear,
			s.toYear,
		).
		WillReturnError(s.testErr)

	resp, err := s.storage.Get(s.ctx, s.position, s.fromYear, s.toYear, s.parser)

	s.Equal(1, s.logs.FilterMessage("can't get jobs info").FilterField(zap.Error(s.testErr)).Len())
	s.Equal(internal.DatabaseErr, err)
	s.Nil(resp)
}

func (s *storageSuite) TestGetNoRows() {
	s.mockDb.
		ExpectQuery(regexp.QuoteMeta("select * from some_procedure($1,$2,$3)")).
		WithArgs(
			s.position,
			s.fromYear,
			s.toYear,
		).
		WillReturnRows(sqlmock.NewRows(s.columnRows))

	resp, err := s.storage.Get(s.ctx, s.position, s.fromYear, s.toYear, s.parser)

	s.Equal(0, s.logs.FilterMessage("can't get jobs info").FilterField(zap.Error(s.testErr)).Len())
	s.Nil(err)
	s.Nil(resp)
}

func (s *storageSuite) TestGetStructScanErr() {
	s.mockDb.
		ExpectQuery(regexp.QuoteMeta("select * from some_procedure($1,$2,$3)")).
		WithArgs(
			s.position,
			s.fromYear,
			s.toYear,
		).
		WillReturnRows(
			sqlmock.
				NewRows(s.columnRows).
				AddRow(
					"test-position",
					"100",
					"200",
					"150",
					"",
					"test-parser",
					"test",
				),
		)

	resp, err := s.storage.Get(s.ctx, s.position, s.fromYear, s.toYear, s.parser)

	s.Equal(1, s.logs.FilterMessage("can't scan raw to struct").Len())
	s.Equal(internal.DatabaseErr, err)
	s.Nil(resp)
}

func (s *storageSuite) TestGetOk() {
	popularSkillsJson, err := json.Marshal([]string{"test skill 1", "test skill 2", "test skill 3"})
	s.Nil(err)

	now := time.Unix(time.Now().Unix(), 0)
	s.mockDb.
		ExpectQuery(regexp.QuoteMeta("select * from some_procedure($1,$2,$3)")).
		WithArgs(
			s.position,
			s.fromYear,
			s.toYear,
		).
		WillReturnRows(
			sqlmock.
				NewRows(s.columnRows).
				AddRow(
					"test-position",
					"100",
					"200",
					"150",
					popularSkillsJson,
					"test-parser",
					now,
				),
		)

	resp, err := s.storage.Get(s.ctx, s.position, s.fromYear, s.toYear, s.parser)

	s.Equal(0, s.logs.FilterMessage("can't scan raw to struct").Len())
	s.Nil(err)
	s.Equal(&internal.JobsInfo{
		PositionToParse: "test-position",
		MinSalary:       100,
		MaxSalary:       200,
		MedianSalary:    150,
		PopularSkills:   []string{"test skill 1", "test skill 2", "test skill 3"},
		Parser:          "test-parser",
		Time:            &now,
	}, resp)
}

type dbJobsInfoResponseSuite struct {
	suite.Suite
}

func TestDbJobsInfoResponseSuite(t *testing.T) {
	suite.Run(t, &dbJobsInfoResponseSuite{})
}

func (s *dbJobsInfoResponseSuite) TestMapToDomain() {
	popularSkillsJson, err := json.Marshal([]string{"test skill 1", "test skill 2", "test skill 3"})
	s.Nil(err)
	now := time.Now()

	r := dbJobsInfoResponse{
		PositionToParse: "test-position",
		MinSalary:       100,
		MaxSalary:       200,
		MedianSalary:    150,
		PopularSkills:   popularSkillsJson,
		Parser:          "test-parser",
		Time:            &now,
	}

	expected := &internal.JobsInfo{
		PositionToParse: "test-position",
		MinSalary:       100,
		MaxSalary:       200,
		MedianSalary:    150,
		PopularSkills:   []string{"test skill 1", "test skill 2", "test skill 3"},
		Parser:          "test-parser",
		Time:            &now,
	}
	actual, err := r.mapToDomain()

	s.Nil(err)
	s.Equal(expected, actual)
}
