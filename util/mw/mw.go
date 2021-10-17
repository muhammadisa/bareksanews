package mw

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			req, _ := json.Marshal(request)
			defer func(begin time.Time) {
				level.Info(logger).Log(
					"err", err,
					"took", time.Since(begin),
					"request", string(req),
				)
			}(time.Now())
			resp, err = next(ctx, request)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}

func CircuitBreakerMiddleware(command string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			var logicErr error
			err = hystrix.Do(command, func() (err error) {
				resp, logicErr = next(ctx, request)
				return logicErr
			}, func(err error) error {
				return err
			})
			if logicErr != nil {
				return nil, logicErr
			}
			if err != nil {
				return nil, status.Error(
					codes.Unavailable,
					errors.New("service is busy or unavailable, please try again later").Error(),
				)
			}
			return resp, nil
		}
	}
}

