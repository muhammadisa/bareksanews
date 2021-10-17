package cache

import (
	"context"
	"encoding/json"

	"github.com/muhammadisa/bareksanews/constant"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

func (c *cache) GetNewses(ctx context.Context) (res *pb.Newses, err error) {
	const funcName = `GetNewses`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var newses pb.Newses
	tagsMap := c.redis.HGetAll(ctx, constant.News).Val()
	for _, v := range tagsMap {
		var news pb.News
		err = json.Unmarshal([]byte(v), &news)
		if err != nil {
			return res, err
		}
		newses.Newses = append(newses.Newses, &news)
	}
	c.redis.Set(ctx, constant.ReloadNewses, false, 0)
	return &newses, nil
}

func (c *cache) GetFilter(ctx context.Context) string {
	const funcName = `GetFilter`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	return c.redis.Get(ctx, constant.FilterNewses).Val()
}

func (c *cache) SetFilter(ctx context.Context, filter string) error {
	const funcName = `SetFilter`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	return c.redis.Set(ctx, constant.FilterNewses, filter, 0).Err()
}

func (c *cache) UnsetNews(ctx context.Context, id string) error {
	const funcName = `UnsetNews`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	return c.redis.HDel(ctx, constant.News, id).Err()
}

func (c *cache) InvalidateNewses(ctx context.Context) error {
	const funcName = `InvalidateNewses`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	return c.redis.Set(ctx, constant.ReloadNewses, true, 0).Err()
}

func (c *cache) ReloadRequired(ctx context.Context, filter string) bool {
	const funcName = `ReloadRequired`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	fil := c.redis.Get(ctx, constant.FilterNewses).Val()
	state := c.redis.Get(ctx, constant.ReloadNewses).Val()
	return state != "0" || fil != filter
}

func (c *cache) ReloadNewses(ctx context.Context, newses *pb.Newses) (err error) {
	const funcName = `ReloadNewses`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	c.redis.Del(ctx, constant.News)
	data := make(map[string]interface{})
	for _, news := range newses.Newses {
		newsByte, err := json.Marshal(news)
		if err != nil {
			return err
		}
		data[news.Id] = string(newsByte)
	}
	c.redis.Set(ctx, constant.ReloadNewses, false, 0)
	return c.redis.HMSet(ctx, constant.News, data).Err()
}
