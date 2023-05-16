package internal

import (
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mock/redis_mock.go -package=mock github.com/redis/go-redis/v9 Cmdable
