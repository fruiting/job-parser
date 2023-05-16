package pgsql

import (
	"context"

	"fruiting/job-parser/internal"
)

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Set(
	ctx context.Context,
	position internal.Name,
	minSalary internal.Salary,
	maxSalary internal.Salary,
	medianSalary internal.Salary,
	parser internal.Parser,
) error {
	return nil
}
