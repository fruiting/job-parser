package pgsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"fruiting/job-parser/internal"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Storage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewStorage(db *sqlx.DB, logger *zap.Logger) *Storage {
	return &Storage{
		db:     db,
		logger: logger,
	}
}

func (s *Storage) Set(
	ctx context.Context,
	position internal.Name,
	minSalary internal.Salary,
	maxSalary internal.Salary,
	medianSalary internal.Salary,
	parser internal.Parser,
) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		"call some_procedure($1,$2,$3,$4,$5)", //todo
		position,
		minSalary,
		maxSalary,
		medianSalary,
		parser,
	)
	if err != nil {
		return fmt.Errorf("can't set jobs info: %w", err)
	}

	return nil
}

func (s *Storage) Get(
	ctx context.Context,
	positionName internal.Name,
	fromYear uint16,
	toYear uint16,
	parser internal.Parser,
) (*internal.JobsInfo, error) {
	ctxLogger := s.logger.With(
		zap.String("position_name", string(positionName)),
		zap.Uint16("from_year", fromYear),
		zap.Uint16("to_year", toYear),
		zap.String("parser", string(parser)),
	)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(
		ctx,
		"select * from some_procedure($1,$2,$3)",
		positionName,
		fromYear,
		toYear,
		parser,
	)
	if row.Err() != nil {
		s.logger.Error("can't get jobs info", zap.Error(row.Err()))
		return nil, internal.DatabaseErr
	}

	var raw dbJobsInfoResponse
	err := row.Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		ctxLogger.Error("can't scan raw to struct", zap.Error(err))
		return nil, internal.DatabaseErr
	}

	result, err := raw.mapToDomain()
	if err != nil {
		ctxLogger.Error("can't map jobs info to domain", zap.Error(err))
		return nil, internal.DatabaseErr
	}

	return result, nil
}

type dbJobsInfoResponse struct {
	PositionToParse string     `db:"position_to_parse"`
	MinSalary       uint32     `db:"min_salary"`
	MaxSalary       uint32     `db:"max_salary"`
	MedianSalary    uint32     `db:"median_salary"`
	PopularSkills   []byte     `db:"popular_skills"`
	Parser          string     `db:"parser"`
	Time            *time.Time `db:"mdate"`
}

func (r *dbJobsInfoResponse) mapToDomain() (*internal.JobsInfo, error) {
	var popularSkills []string
	err := json.Unmarshal(r.PopularSkills, &popularSkills)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal: %w", err)
	}

	return &internal.JobsInfo{
		PositionToParse: internal.Name(r.PositionToParse),
		MinSalary:       internal.Salary(r.MinSalary),
		MaxSalary:       internal.Salary(r.MaxSalary),
		MedianSalary:    internal.Salary(r.MedianSalary),
		PopularSkills:   popularSkills,
		Parser:          internal.Parser(r.Parser),
		Time:            r.Time,
	}, nil
}
