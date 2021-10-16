package repository

import (
	"context"

	"github.com/muhammadisa/bareksanews/repository/cache"
	_interface "github.com/muhammadisa/bareksanews/repository/interface"
	"github.com/muhammadisa/bareksanews/repository/sql"
	"github.com/muhammadisa/barektest-util/dbc"
	"go.opencensus.io/trace"
)

type Repository struct {
	ReadWriter      _interface.ReadWrite
	CacheReadWriter _interface.Cache
}

type RepoConf struct {
	SQL   dbc.Config
	Cache dbc.Config
}

func NewRepository(ctx context.Context, rc RepoConf, tracer trace.Tracer) (*Repository, error) {
	readWriter, err := sql.NewSQL(rc.SQL, tracer)
	if err != nil {
		return nil, err
	}
	cacheReadWriter, err := cache.NewCache(rc.Cache, tracer)
	if err != nil {
		return nil, err
	}
	return &Repository{
		ReadWriter:      readWriter,
		CacheReadWriter: cacheReadWriter,
	}, nil
}
