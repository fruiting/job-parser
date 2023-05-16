package internal

import "context"

//go:generate mockgen -source=storage.go -destination=./storage_mock.go -package=internal

// Storage for long-term info
type Storage interface {
	Set(ctx context.Context, jobsInfo *JobsInfo) error
}
