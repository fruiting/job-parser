package internal

//go:generate mockgen -source=storage.go -destination=./storage_mock.go -package=internal

import "context"

// Storage for long-term info
type Storage interface {
	Set(
		ctx context.Context,
		position Name,
		minSalary Salary,
		maxSalary Salary,
		medianSalary Salary,
		parser Parser,
	) error
	Get(
		ctx context.Context,
		positionName Name,
		fromYear uint16,
		toYear uint16,
		parser Parser,
	) (*JobsInfo, error)
}
