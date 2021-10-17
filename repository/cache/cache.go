package cache

import (
	"github.com/go-redis/redis/v8"
	_interface "github.com/muhammadisa/bareksanews/repository/interface"
	"github.com/muhammadisa/bareksanews/util/dbc"
	"go.opencensus.io/trace"
)

type cache struct {
	tracer trace.Tracer
	redis  redis.Cmdable
}

func NewCache(config dbc.Config, tracer trace.Tracer) (_interface.Cache, error) {
	redisDB, err := dbc.OpenRedis(config)
	if err != nil {
		return nil, err
	}
	return &cache{
		redis:  redisDB,
		tracer: tracer,
	}, nil
}
